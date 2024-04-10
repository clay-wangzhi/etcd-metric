package alarm

import (
	"sync"

	etcdv1alpha1 "etcd-operator/api/etcd/v1alpha1"
	"etcd-operator/pkg/featureprovider"
	"etcd-operator/pkg/inspection"
)

const (
	ProviderName = string(etcdv1alpha1.KStoneFeatureAlarm)
)

var (
	once     sync.Once
	instance *FeatureAlarm
)

type FeatureAlarm struct {
	name       string
	inspection *inspection.Server
	ctx        *featureprovider.FeatureContext
}

func init() {
	featureprovider.RegisterFeatureFactory(
		ProviderName,
		func(ctx *featureprovider.FeatureContext) (featureprovider.Feature, error) {
			return initFeatureAlarmInstance(ctx)
		},
	)
}

func initFeatureAlarmInstance(ctx *featureprovider.FeatureContext) (featureprovider.Feature, error) {
	var err error
	once.Do(func() {
		instance = &FeatureAlarm{
			name: ProviderName,
			ctx:  ctx,
		}
		instance.inspection, err = inspection.NewInspectionServer(ctx)
	})
	return instance, err
}

func (c *FeatureAlarm) Equal(cluster *etcdv1alpha1.EtcdCluster) bool {
	return c.inspection.Equal(cluster, etcdv1alpha1.KStoneFeatureAlarm)
}

func (c *FeatureAlarm) Sync(cluster *etcdv1alpha1.EtcdCluster) error {
	return c.inspection.Sync(cluster, etcdv1alpha1.KStoneFeatureAlarm)
}

func (c *FeatureAlarm) Do(inspection *etcdv1alpha1.EtcdInspection) error {
	return c.inspection.CollectAlarmList(inspection)
}
