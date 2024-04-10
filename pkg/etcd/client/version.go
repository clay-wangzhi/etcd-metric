package client

import (
	"errors"
	"sync"

	etcdv1alpha1 "etcd-operator/api/etcd/v1alpha1"

	"k8s.io/klog/v2"
)

type Factory func(cluster *VersionContext) (VersionClient, error)

var (
	mutex     sync.Mutex
	providers = make(map[etcdv1alpha1.EtcdStorageBackend]Factory)
)

// RegisterEtcdClientFactory registers the specified etcd client
func RegisterEtcdClientFactory(name etcdv1alpha1.EtcdStorageBackend, factory Factory) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, found := providers[name]; found {
		klog.V(2).Infof("etcdcluster provider %s was registered twice", name)
	}

	klog.V(2).Infof("register etcdCluster provider %s", name)
	providers[name] = factory
}

// GetEtcdClientProvider gets the specified etcd client
func GetEtcdClientProvider(
	name etcdv1alpha1.EtcdStorageBackend,
	ctx *VersionContext,
) (VersionClient, error) {
	mutex.Lock()
	defer mutex.Unlock()

	// compatible with existing clusters
	if name == "" {
		name = etcdv1alpha1.EtcdStorageV3
	}
	f, found := providers[name]

	klog.V(1).Infof("get provider name %s,status:%t", name, found)
	if !found {
		return nil, errors.New("fatal error,etcd cluster provider not found")
	}
	return f(ctx)
}
