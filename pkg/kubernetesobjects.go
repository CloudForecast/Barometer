package pkg

import (
	"fmt"
	"github.com/CloudForecast/Barometer/pkg/barometerApi"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"strings"
)

type ErrorReporter func(error)

func fetchKubernetesResources(resourceTypes []string, errorReporter ErrorReporter) (map[string][]resource.Info, error) {
	configFlags := genericclioptions.NewConfigFlags(true)
	// This is a bit hacky, but will work
	if k := viper.GetString("kubeconfg"); k != "" {
		configFlags.KubeConfig = &k
	}

	builder := resource.NewBuilder(configFlags)
	builder = builder.Unstructured().AllNamespaces(true).SelectAllParam(true)
	builder = builder.Flatten().RequireObject(true).ContinueOnError()

	// check to see what resource types aren't present in the cluster, exclude from fetching, and report the error.
	// if we attempt to fetch even one unavailable resource, the entire set of API calls will fail.
	var filteredResourceTypes []string
	var excludedResourceTypes []string
	for _, resourceType := range resourceTypes {
		if DoesClusterContainResourceName(resourceType) {
			filteredResourceTypes = append(filteredResourceTypes, resourceType)
		} else {
			excludedResourceTypes = append(excludedResourceTypes, resourceType)
		}
	}
	if len(excludedResourceTypes) > 0 {
		err := fmt.Errorf("the following requested kubernetes resource types are not available in the cluster: %s", strings.Join(excludedResourceTypes, ","))
		defer errorReporter(err)
		log.Error().Err(err).Msg("")
	}

	result := builder.ResourceTypes(filteredResourceTypes...).Do()
	outputResults := make(map[string][]resource.Info)
	err := result.Visit(func(info *resource.Info, err error) error {
		log.Debug().Msgf("handling resource.Info: %v", *info)
		if err != nil {
			log.Error().Err(err).Msgf("received error on resource item")
			defer errorReporter(err)
			// Docs: https://pkg.go.dev/k8s.io/cli-runtime/pkg/resource@v0.20.6#VisitorFunc
			// To aggregate or handle errors but continue on, return nil rather than an error
			return nil
		}

		kind := info.Object.GetObjectKind().GroupVersionKind().Kind
		if _, ok := outputResults[kind]; ok {
			outputResults[kind] = append(outputResults[kind], *info)
		} else {
			outputResults[kind] = []resource.Info{*info}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return outputResults, nil
}

func fetchAndRunKubeInstructions(b barometerApi.ApiClient) (*barometerApi.BarometerK8sApiResultsEventData, error) {
	queryInstructions, err := b.GetKubeInstructions()
	if err != nil {
		return nil, err
	}

	kindsToFetch := queryInstructions.KindsToFetch
	if len(kindsToFetch) == 0 {
		log.Warn().Msg("no k8s kinds returned from barometer API to fetch")
		return nil, nil
	}

	var stringKinds []string
	for _, i := range kindsToFetch {
		stringKinds = append(stringKinds, string(i))
	}

	dryRun := viper.GetBool("dryrun")
	var errorReporter func(err error)
	if dryRun {
		errorReporter = func(err error) {}
	} else {
		errorReporter = func(err error) {
			_ = b.SendExceptionEvent(err)
		}
	}

	result, err := fetchKubernetesResources(stringKinds, errorReporter)
	if err != nil {
		return nil, err
	}

	outMap := make(map[string][]interface{})
	for key, elem := range result {
		for _, i := range elem {
			var obj map[string]interface{}
			err := mapstructure.Decode(i.Object, &obj)
			if err != nil {
				panic(err)
			}
			outMap[key] = append(outMap[key], obj["Object"])
		}
	}

	log.Trace().Msgf("outMap: %v", outMap)

	eventData := barometerApi.BarometerK8sApiResultsEventData{
		K8sClusterInformation: outMap,
	}
	return &eventData, nil
}

func FetchAndSubmitKubernetesObjects(b barometerApi.ApiClient) error {
	results, err := fetchAndRunKubeInstructions(b)
	if err != nil {
		return err
	}
	if results == nil {
		return nil
	}

	log.Debug().Msg("sending k8s result data to Barometer...")
	err = b.SendK8sAPIResultsEvent(*results)
	if err != nil {
		return err
	}
	return nil
}
