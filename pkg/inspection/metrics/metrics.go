package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	EtcdNodeDiffTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "inspection",
		Name:      "etcd_node_diff_total",
		Help:      "total etcd node diff key",
	}, []string{"clusterName"})

	EtcdEndpointHealthy = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "inspection",
		Name:      "etcd_endpoint_healthy",
		Help:      "The healthy of etcd member",
	}, []string{"clusterName", "endpoint"})

	EtcdRequestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Subsystem: "inspection",
		Name:      "etcd_request_total",
		Help:      "The total number of etcd requests",
	}, []string{"clusterName", "grpcMethod", "etcdPrefix", "resourceName"})

	EtcdKeyTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "inspection",
		Name:      "etcd_key_total",
		Help:      "The total number of etcd key",
	}, []string{"clusterName", "etcdPrefix", "resourceName"})

	EtcdEndpointAlarm = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "inspection",
		Name:      "etcd_endpoint_alarm",
		Help:      "The alarm of etcd member",
	}, []string{"clusterName", "endpoint", "alarmType"})

	EtcdNodeRevisionDiff = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "inspection",
		Name:      "etcd_node_revision_diff_total",
		Help:      "The revision difference between all member",
	}, []string{"clusterName"})

	EtcdNodeIndexDiff = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "inspection",
		Name:      "etcd_node_index_diff_total",
		Help:      "The index difference between all member",
	}, []string{"clusterName"})

	EtcdNodeRaftAppliedIndexDiff = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "inspection",
		Name:      "etcd_node_raft_applied_index_diff_total",
		Help:      "The raft applied index difference between all member",
	}, []string{"clusterName"})

	EtcdNodeRaftIndexDiff = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "inspection",
		Name:      "etcd_node_raft_index_diff_total",
		Help:      "The raft index difference between all member",
	}, []string{"clusterName"})

	EtcdBackupFiles = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "inspection",
		Name:      "etcd_backup_files",
		Help:      "The Number of backup files in the last day",
	}, []string{"clusterName"})

	EtcdFailedBackupFiles = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: "inspection",
		Name:      "etcd_failed_backup_files",
		Help:      "The Number of failed backup files in the last day",
	}, []string{"clusterName"})

	EtcdInspectionFailedNum = prometheus.NewCounterVec(prometheus.CounterOpts{
		Subsystem: "inspection",
		Name:      "failed_num",
		Help:      "The total Number of failed inspection",
	}, []string{"clusterName", "inspectionType"})
)

func init() {
	prometheus.MustRegister(EtcdNodeDiffTotal)
	prometheus.MustRegister(EtcdEndpointHealthy)
	prometheus.MustRegister(EtcdRequestTotal)
	prometheus.MustRegister(EtcdKeyTotal)
	prometheus.MustRegister(EtcdEndpointAlarm)
	prometheus.MustRegister(EtcdNodeRevisionDiff)
	prometheus.MustRegister(EtcdNodeIndexDiff)
	prometheus.MustRegister(EtcdNodeRaftAppliedIndexDiff)
	prometheus.MustRegister(EtcdNodeRaftIndexDiff)
	prometheus.MustRegister(EtcdBackupFiles)
	prometheus.MustRegister(EtcdFailedBackupFiles)
	prometheus.MustRegister(EtcdInspectionFailedNum)
}
