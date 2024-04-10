package imported

import (
	"sync"

	etcdv1alpha1 "etcd-operator/api/etcd/v1alpha1"
	"etcd-operator/pkg/clusterprovider"
	"etcd-operator/pkg/controllers/util"
	"etcd-operator/pkg/etcd"
)

const (
	AnnoImportedURI = "importedAddr"
)

var (
	once     sync.Once
	instance *EtcdClusterImported
)

// EtcdClusterImported is the etcd cluster imported from kstone-dashboard
type EtcdClusterImported struct {
	name etcdv1alpha1.EtcdClusterType
	ctx  *clusterprovider.ClusterContext
}

// init registers an imported etcd cluster provider
func init() {
	clusterprovider.RegisterEtcdClusterFactory(
		etcdv1alpha1.EtcdClusterImported,
		func(ctx *clusterprovider.ClusterContext) (clusterprovider.Cluster, error) {
			return initEtcdClusterImportedInstance(ctx)
		},
	)
}

func initEtcdClusterImportedInstance(ctx *clusterprovider.ClusterContext) (clusterprovider.Cluster, error) {
	once.Do(func() {
		instance = &EtcdClusterImported{
			name: etcdv1alpha1.EtcdClusterImported,
			ctx: &clusterprovider.ClusterContext{
				Clientbuilder: ctx.Clientbuilder,
				Client:        ctx.Clientbuilder.DynamicClientOrDie(),
			},
		}
	})
	return instance, nil
}

func (c *EtcdClusterImported) BeforeCreate(cluster *etcdv1alpha1.EtcdCluster) error {
	return nil
}

func (c *EtcdClusterImported) Create(cluster *etcdv1alpha1.EtcdCluster) error {
	return nil
}

func (c *EtcdClusterImported) AfterCreate(cluster *etcdv1alpha1.EtcdCluster) error {
	return nil
}

func (c *EtcdClusterImported) BeforeUpdate(cluster *etcdv1alpha1.EtcdCluster) error {
	return nil
}

func (c *EtcdClusterImported) Update(cluster *etcdv1alpha1.EtcdCluster) error {
	return nil
}

func (c *EtcdClusterImported) AfterUpdate(cluster *etcdv1alpha1.EtcdCluster) error {
	return nil
}

func (c *EtcdClusterImported) BeforeDelete(cluster *etcdv1alpha1.EtcdCluster) error {
	return nil
}

func (c *EtcdClusterImported) Delete(cluster *etcdv1alpha1.EtcdCluster) error {
	return nil
}

func (c *EtcdClusterImported) AfterDelete(cluster *etcdv1alpha1.EtcdCluster) error {
	return nil
}

func (c *EtcdClusterImported) Equal(cluster *etcdv1alpha1.EtcdCluster) (bool, error) {
	return true, nil
}

// Status gets the imported etcd cluster status
func (c *EtcdClusterImported) Status(config *etcd.ClientConfig, cluster *etcdv1alpha1.EtcdCluster) (etcdv1alpha1.EtcdClusterStatus, error) {
	status := cluster.Status

	annotations := cluster.ObjectMeta.Annotations
	if annotations == nil {
		annotations = make(map[string]string)
	}

	endpoints := clusterprovider.GetStorageMemberEndpoints(cluster)

	if len(endpoints) == 0 {
		if addr, found := annotations[AnnoImportedURI]; found {
			endpoints = append(endpoints, addr)
			status.ServiceName = addr
		} else {
			status.Phase = etcdv1alpha1.EtcdClusterUnknown
			return status, nil
		}
	}

	members, err := clusterprovider.GetRuntimeEtcdMembers(
		cluster.Spec.StorageBackend,
		endpoints,
		cluster.Annotations[util.ClusterExtensionClientURL],
		config,
	)
	if err != nil && len(members) == 0 {
		status.Phase = etcdv1alpha1.EtcdClusterUnknown
		return status, err
	}

	status.Members, status.Phase = clusterprovider.GetEtcdClusterMemberStatus(members, config)
	return status, err
}
