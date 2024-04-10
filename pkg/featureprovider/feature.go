package featureprovider

import (
	"etcd-operator/api/etcd/v1alpha1"
	"etcd-operator/pkg/controllers/util"
	"etcd-operator/pkg/etcd"
)

// Feature is an abstract, pluggable interface for cluster features.
type Feature interface {
	// Equal checks whether the feature needs to be updated
	Equal(cluster *v1alpha1.EtcdCluster) bool

	// Sync synchronizes the latest feature configuration
	Sync(cluster *v1alpha1.EtcdCluster) error

	// Do executes inspection tasks.
	Do(task *v1alpha1.EtcdInspection) error
}
type FeatureContext struct {
	ClientBuilder      util.ClientBuilder
	ClientConfigGetter etcd.ClientConfigGetter
}
