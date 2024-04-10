package clusterprovider

import (
	"errors"
	"sync"

	"k8s.io/klog/v2"

	etcdv1alpha1 "etcd-operator/api/etcd/v1alpha1"
)

type EtcdFactory func(cluster *ClusterContext) (Cluster, error)

var (
	mutex     sync.Mutex
	providers = make(map[etcdv1alpha1.EtcdClusterType]EtcdFactory)
)

// RegisterEtcdClusterFactory registers the specified cluster provider
func RegisterEtcdClusterFactory(name etcdv1alpha1.EtcdClusterType, factory EtcdFactory) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, found := providers[name]; found {
		klog.V(2).Infof("etcdcluster provider %s was registered twice", name)
	}

	klog.V(2).Infof("register etcdCluster provider %s", name)
	providers[name] = factory
}

// GetEtcdClusterProvider gets the specified cluster provider
func GetEtcdClusterProvider(
	name etcdv1alpha1.EtcdClusterType,
	ctx *ClusterContext,
) (Cluster, error) {
	mutex.Lock()
	defer mutex.Unlock()
	f, found := providers[name]

	klog.V(1).Infof("get provider name %s,status:%t", name, found)
	if !found {
		return nil, errors.New("fatal error,etcd cluster provider not found")
	}
	return f(ctx)
}
