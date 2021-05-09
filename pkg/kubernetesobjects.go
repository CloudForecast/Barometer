package pkg

import (
	"context"
	"fmt"
	"github.com/CloudForecast/Barometer/pkg/barometerApi"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

var (
	kubeConfig *rest.Config
	kubeResources map[string]schema.GroupVersionResource
)

func Setup() *rest.Config {
	if kubeConfig != nil {
		return kubeConfig
	}

	log.Debug().Msg("setting up kubeconfig for first time...")
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	kubeConfig = config
	return kubeConfig
}


func getKubeClient() (dynamic.Interface, error) {
	config := Setup()

	dc, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return dc, nil
}

func discoverKubeResources() (map[string]schema.GroupVersionResource, error) {
	if kubeResources != nil {
		return kubeResources, nil
	}

	log.Debug().Msg("doing k8s resource discovery...")
	config := Setup()
	discoveryclient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}
	resourceLists, err := discoveryclient.ServerPreferredResources()
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch resource list")
	}

	resources := make(map[string]schema.GroupVersionResource)

	for _, resourceList := range resourceLists {
		for _, resource := range resourceList.APIResources {
			resources[resource.Name] = schema.GroupVersionResource{
				Group: resource.Group,
				Version: resource.Version,
				Resource: resource.Name,
			}
		}
	}

	kubeResources = resources
	log.Debug().Msg("k8s resource discovery complete")
	return resources, nil
}

func fetchKubeResource(client dynamic.Interface, kubeResource string) (*unstructured.UnstructuredList, error) {
	resources, err := discoverKubeResources()
	if err != nil {
		return nil, err
	}

	if val, ok := resources[kubeResource]; ok {
		list, err := client.Resource(val).List(context.TODO(), v1.ListOptions{})
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("error fetching %s", kubeResource))
		}
		return list, nil
	} else {
		log.Warn().Msgf("tried to fetch resource '%s' but not found by discovery", kubeResource)
		return nil, nil
	}
}

func fetchKubernetesResources(resourceTypes []string) (map[string][]resource.Info, error) {
	configFlags := genericclioptions.NewConfigFlags(true)

	builder := resource.NewBuilder(configFlags)
	builder = builder.Unstructured().AllNamespaces(true).SelectAllParam(true)
	builder = builder.Flatten().RequireObject(true).ContinueOnError()
	result := builder.ResourceTypes(resourceTypes...).Do()
	outputResults := make(map[string][]resource.Info)
	err := result.Visit(func (info *resource.Info, err error) error {
		if err != nil {
			// TODO: send to Cloudforecast
			return err
		}

		kind := info.Object.GetObjectKind().GroupVersionKind().String()
		if _, ok := outputResults[kind]; ok {
			outputResults[kind] = append(outputResults[kind], *info)
		} else {
			outputResults[kind] = []resource.Info{*info}
		}
		return nil
	})

	log.Debug().Msgf("outputResults: %v", outputResults)

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
	result, err := fetchKubernetesResources(stringKinds)
	if err != nil {
		return nil, err
	}

	// convert from []resource.Info to []interface{}
	outMap := make(map[string][]interface{})
	for key, elem := range result {
		for _, i := range elem {
			outMap[key] = append(outMap[key], interface{}(i))
		}
	}

	eventData := barometerApi.BarometerK8sApiResultsEventData{
		K8sClusterInformation: outMap,
	}
	return &eventData, nil
}

//func oldfetchAndRunKubeInstructions(b barometerApi.ApiClient) (*barometerApi.BarometerK8sApiResultsEventData, error) {
//	outputResults := make(map[string]interface{})
//
//	kubeClient, err := getKubeClient()
//	if err != nil {
//		return nil, err
//	}
//
//	queryInstructions, err := b.GetKubeInstructions()
//	if err != nil {
//		return nil, err
//	}
//
//	kindsToFetch := queryInstructions.KindsToFetch
//	if len(kindsToFetch) == 0 {
//		log.Warn().Msg("no k8s kinds returned from barometer API to fetch")
//		return nil, nil
//	}
//	for _, kindToFetch := range kindsToFetch {
//		log.Debug().Msgf("fetching kube resource '%s'...", kindToFetch)
//		list, err := fetchKubeResource(kubeClient, string(kindToFetch))
//		if err != nil {
//			err = errors.Wrap(err, fmt.Sprintf("error fetching kube resource '%s'", kindToFetch))
//			//return nil, err
//			log.Error().Err(err).Msg("")
//			continue
//		}
//		if list == nil {
//			continue
//		}
//
//		outputResults[string(kindToFetch)] = list
//		log.Debug().Msgf("fetched %d of kube resource %s", len(list.Items), kindsToFetch)
//	}
//
//	return &barometerApi.BarometerK8sApiResultsEventData{
//		K8sClusterInformation: outputResults,
//	}, nil
//}

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