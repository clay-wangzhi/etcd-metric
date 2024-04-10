package consistency

import (
	"sync"

	etcdv1alpha1 "etcd-operator/api/etcd/v1alpha1"
	"etcd-operator/pkg/featureprovider"
	"etcd-operator/pkg/inspection"
)

const (
	ProviderName = string(etcdv1alpha1.KStoneFeatureConsistency)
)

var (
	once     sync.Once
	instance *FeatureConsistency
)

type FeatureConsistency struct {
	name       string
	inspection *inspection.Server
	ctx        *featureprovider.FeatureContext
}

func init() {
	featureprovider.RegisterFeatureFactory(
		ProviderName,
		func(ctx *featureprovider.FeatureContext) (featureprovider.Feature, error) {
			return initFeatureConsistencyInstance(ctx)
		},
	)
}

func initFeatureConsistencyInstance(ctx *featureprovider.FeatureContext) (featureprovider.Feature, error) {
	var err error
	once.Do(func() {
		instance = &FeatureConsistency{
			name: ProviderName,
			ctx:  ctx,
		}
		instance.inspection, err = inspection.NewInspectionServer(ctx)
	})
	return instance, err
}

func (c *FeatureConsistency) Equal(cluster *etcdv1alpha1.EtcdCluster) bool {
	return c.inspection.Equal(cluster, etcdv1alpha1.KStoneFeatureConsistency)
}

func (c *FeatureConsistency) Sync(cluster *etcdv1alpha1.EtcdCluster) error {
	return c.inspection.Sync(cluster, etcdv1alpha1.KStoneFeatureConsistency)
}

func (c *FeatureConsistency) Do(inspection *etcdv1alpha1.EtcdInspection) error {
	return c.inspection.CollectClusterConsistentData(inspection)
}
