package healthy

import (
	"sync"

	etcdv1alpha1 "etcd-operator/api/etcd/v1alpha1"
	"etcd-operator/pkg/featureprovider"
	"etcd-operator/pkg/inspection"
)

var (
	once     sync.Once
	instance *FeatureHealthy
)

type FeatureHealthy struct {
	name       string
	inspection *inspection.Server
	ctx        *featureprovider.FeatureContext
}

const (
	ProviderName = string(etcdv1alpha1.KStoneFeatureHealthy)
)

func init() {
	featureprovider.RegisterFeatureFactory(
		ProviderName,
		func(ctx *featureprovider.FeatureContext) (featureprovider.Feature, error) {
			return initFeatureHealthyInstance(ctx)
		},
	)
}

func initFeatureHealthyInstance(ctx *featureprovider.FeatureContext) (featureprovider.Feature, error) {
	var err error
	once.Do(func() {
		instance = &FeatureHealthy{
			name: ProviderName,
			ctx:  ctx,
		}
		instance.inspection, err = inspection.NewInspectionServer(ctx)
	})
	return instance, err
}

func (c *FeatureHealthy) Equal(cluster *etcdv1alpha1.EtcdCluster) bool {
	return c.inspection.Equal(cluster, etcdv1alpha1.KStoneFeatureHealthy)
}

func (c *FeatureHealthy) Sync(cluster *etcdv1alpha1.EtcdCluster) error {
	return c.inspection.Sync(cluster, etcdv1alpha1.KStoneFeatureHealthy)
}

func (c *FeatureHealthy) Do(inspection *etcdv1alpha1.EtcdInspection) error {
	return c.inspection.CollectMemberHealthy(inspection)
}
