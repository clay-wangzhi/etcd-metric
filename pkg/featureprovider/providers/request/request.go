package request

import (
	"sync"

	etcdv1alpha1 "etcd-operator/api/etcd/v1alpha1"
	"etcd-operator/pkg/featureprovider"
	"etcd-operator/pkg/inspection"
)

const (
	ProviderName = string(etcdv1alpha1.KStoneFeatureRequest)
)

var (
	once     sync.Once
	instance *FeatureRequest
)

type FeatureRequest struct {
	name       string
	inspection *inspection.Server
	ctx        *featureprovider.FeatureContext
}

func init() {
	featureprovider.RegisterFeatureFactory(
		ProviderName,
		func(ctx *featureprovider.FeatureContext) (featureprovider.Feature, error) {
			return initFeatureRequestInstance(ctx)
		},
	)
}

func initFeatureRequestInstance(ctx *featureprovider.FeatureContext) (featureprovider.Feature, error) {
	var err error
	once.Do(func() {
		instance = &FeatureRequest{
			name: ProviderName,
			ctx:  ctx,
		}
		instance.inspection, err = inspection.NewInspectionServer(ctx)
	})
	return instance, err
}

func (c *FeatureRequest) Equal(cluster *etcdv1alpha1.EtcdCluster) bool {
	return c.inspection.Equal(cluster, etcdv1alpha1.KStoneFeatureRequest)
}

func (c *FeatureRequest) Sync(cluster *etcdv1alpha1.EtcdCluster) error {
	return c.inspection.Sync(cluster, etcdv1alpha1.KStoneFeatureRequest)
}

func (c *FeatureRequest) Do(inspection *etcdv1alpha1.EtcdInspection) error {
	return c.inspection.CollectEtcdClusterRequest(inspection)
}
