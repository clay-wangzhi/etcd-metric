package util

import (
	"fmt"
	"strconv"
	"strings"

	etcdv1alpha1 "etcd-operator/api/etcd/v1alpha1"
	// "tkestack.io/kstone/pkg/inspection/metrics"
)

const (
	FeatureStatusEnabled  = "enabled"
	FeatureStatusDisabled = "disabled"
)

type ConsistencyType string

const (
	ConsistencyKeyTotal             ConsistencyType = "keyTotal"
	ConsistencyRevision             ConsistencyType = "revision"
	ConsistencyIndex                ConsistencyType = "index"
	ConsistencyRaftRaftAppliedIndex ConsistencyType = "raftAppliedIndex"
	ConsistencyRaftIndex            ConsistencyType = "raftIndex"
)

const (
	OneDaySeconds = 24 * 60 * 60
)

func IsFeatureGateEnabled(annotations map[string]string, name etcdv1alpha1.KStoneFeature) bool {
	if gates, found := annotations[etcdv1alpha1.KStoneFeatureAnno]; found && gates != "" {
		featurelist := strings.Split(gates, ",")
		for _, item := range featurelist {
			features := strings.Split(item, "=")
			if len(features) != 2 {
				continue
			}

			enabled, _ := strconv.ParseBool(features[1])
			if etcdv1alpha1.KStoneFeature(features[0]) == name && enabled {
				return true
			}
		}
	}
	return false
}

func IncrFailedInspectionCounter(clusterName string, featureName etcdv1alpha1.KStoneFeature) {
	labels := map[string]string{
		"clusterName":    clusterName,
		"inspectionType": string(featureName),
	}
	fmt.Print(labels)
	// metrics.EtcdInspectionFailedNum.With(labels).Inc()
}
