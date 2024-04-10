package featureprovider

import (
	"errors"
	"sync"

	"k8s.io/klog/v2"
)

var (
	mutex                sync.Mutex
	EtcdFeatureProviders = make(map[string]FeatureFactory)
)

type FeatureFactory func(cfg *FeatureContext) (Feature, error)

// RegisterFeatureFactory registers the specified feature provider
func RegisterFeatureFactory(name string, factory FeatureFactory) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, found := EtcdFeatureProviders[name]; found {
		klog.V(2).Infof("feature provider:%s was registered twice", name)
	}

	klog.V(2).Infof("feature provider:%s", name)
	EtcdFeatureProviders[name] = factory
}

// GetFeatureProvider gets the specified feature provider
func GetFeatureProvider(name string, ctx *FeatureContext) (Feature, error) {
	mutex.Lock()
	defer mutex.Unlock()
	f, found := EtcdFeatureProviders[name]

	klog.V(1).Infof("get provider name %s,status:%t", name, found)
	if !found {
		return nil, errors.New("fatal error,feature provider not found")
	}
	return f(ctx)
}

// ListFeatureProvider lists all feature provider
func ListFeatureProvider() []string {
	var features []string
	mutex.Lock()
	defer mutex.Unlock()
	for feature := range EtcdFeatureProviders {
		features = append(features, feature)
	}
	return features
}
