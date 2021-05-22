package pkg

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeConfig    *rest.Config
	KubeResources map[string]schema.GroupVersionResource
)

func Setup() *rest.Config {
	if kubeConfig != nil {
		return kubeConfig
	}

	log.Debug().Msg("setting up kubeconfig for first time...")
	k := viper.GetString("kubeconfig")
	// if no kubeconfig is passed, k is "", which will result in inClusterConfig.
	config, err := clientcmd.BuildConfigFromFlags("", k)
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

// DiscoverKubeResources idempotently fetches the list of available resources in the Kubernetes cluster
// and returns them as map[string]GVR
func DiscoverKubeResources() (map[string]schema.GroupVersionResource, error) {
	if KubeResources != nil {
		return KubeResources, nil
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
				Group:    resource.Group,
				Version:  resource.Version,
				Resource: resource.Name,
			}
		}
	}

	KubeResources = resources
	log.Debug().Msg("k8s resource discovery complete")
	return resources, nil
}

func DoesClusterContainResourceName(resourceName string) bool {
	resources, err := DiscoverKubeResources()
	if err != nil {
		panic(err)
	}
	_, ok := resources[resourceName]
	return ok
}
