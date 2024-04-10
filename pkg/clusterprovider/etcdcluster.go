package clusterprovider

import (
	"k8s.io/client-go/dynamic"

	etcdv1alpha1 "etcd-operator/api/etcd/v1alpha1"
	"etcd-operator/pkg/controllers/util"
	"etcd-operator/pkg/etcd"
)

// Cluster is an abstract, pluggable interface for etcd clusters.
type Cluster interface {

	// BeforeCreate does some things before creating the cluster
	BeforeCreate(cluster *etcdv1alpha1.EtcdCluster) error
	// Create creates the cluster
	Create(cluster *etcdv1alpha1.EtcdCluster) error
	// AfterCreate does some things after creating the cluster
	AfterCreate(cluster *etcdv1alpha1.EtcdCluster) error

	// BeforeUpdate does some things before updating the cluster
	BeforeUpdate(cluster *etcdv1alpha1.EtcdCluster) error
	// Update updates the cluster
	Update(cluster *etcdv1alpha1.EtcdCluster) error
	// AfterUpdate does some things after updating the cluster
	AfterUpdate(cluster *etcdv1alpha1.EtcdCluster) error

	// BeforeDelete does some things before deleting the cluster
	BeforeDelete(cluster *etcdv1alpha1.EtcdCluster) error
	// Delete deletes the cluster
	Delete(cluster *etcdv1alpha1.EtcdCluster) error
	// AfterDelete does some things after deleting the cluster
	AfterDelete(cluster *etcdv1alpha1.EtcdCluster) error

	// Equal checks whether the cluster needs to be updated
	Equal(cluster *etcdv1alpha1.EtcdCluster) (bool, error)

	// Status gets the cluster status
	Status(config *etcd.ClientConfig, cluster *etcdv1alpha1.EtcdCluster) (etcdv1alpha1.EtcdClusterStatus, error)
}

type ClusterContext struct {
	Clientbuilder util.ClientBuilder
	Client        dynamic.Interface
}
