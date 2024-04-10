package clusterprovider

import (
	"fmt"
	"strings"

	"k8s.io/klog/v2"

	etcdv1alpha1 "etcd-operator/api/etcd/v1alpha1"
	"etcd-operator/pkg/etcd"
	versionClient "etcd-operator/pkg/etcd/client"
	_ "etcd-operator/pkg/etcd/client/versions" // import etcd client including v2 and v3
)

type EtcdAlarm struct {
	MemberID  uint64
	AlarmType string
}

// GetStorageMemberEndpoints get member of cluster status
func GetStorageMemberEndpoints(cluster *etcdv1alpha1.EtcdCluster) []string {
	members := cluster.Status.Members
	endpoints := make([]string, 0)
	if len(members) == 0 {
		return endpoints
	}

	for _, m := range members {
		endpoints = append(endpoints, m.ExtensionClientUrl)
	}
	return endpoints
}

// populateExtensionCientURLMap generate extensionClientURLs map
func populateExtensionCientURLMap(extensionClientURLs string) (map[string]string, error) {
	urlMap := make(map[string]string)
	if extensionClientURLs == "" {
		return nil, nil
	}
	items := strings.Split(extensionClientURLs, ",")
	for i := 0; i < len(items); i++ {
		eps := strings.Split(items[i], "->")
		if len(eps) == 2 {
			urlMap[eps[0]] = eps[1]
		} else {
			return urlMap, fmt.Errorf("invalid extensionClientURLs %s", items[i])
		}
	}
	return urlMap, nil
}

// GetRuntimeEtcdMembers get members of etcd
func GetRuntimeEtcdMembers(
	storageBackend string,
	endpoints []string,
	extensionClientURLs string,
	config *etcd.ClientConfig) ([]etcdv1alpha1.MemberStatus, error) {
	etcdMembers := make([]etcdv1alpha1.MemberStatus, 0)

	config.Endpoints = endpoints
	versioned, err := versionClient.GetEtcdClientProvider(etcdv1alpha1.EtcdStorageBackend(storageBackend),
		&versionClient.VersionContext{Config: config})
	if err != nil {
		klog.Errorf("failed get etcd version, err is %v", err)
		return etcdMembers, err
	}

	defer versioned.Close()

	memberRsp, err := versioned.MemberList()
	if err != nil {
		klog.Errorf("failed to get member list, endpoints is %s,err is %v", endpoints, err)
		return etcdMembers, err
	}

	extensionClientURLMap, err := populateExtensionCientURLMap(extensionClientURLs)
	if err != nil {
		klog.Errorf("failed to populate extension clientURL,err is %v", err)
		return etcdMembers, err
	}

	for _, m := range memberRsp {
		// parse url
		if m.ClientURLs == nil {
			continue
		}
		items := strings.Split(m.ClientURLs[0], ":")
		endPoint := strings.TrimPrefix(items[1], "//")

		extensionClientURL := m.ClientURLs[0]
		if extensionClientURLs != "" {
			var ep string
			if strings.HasPrefix(m.ClientURLs[0], "https://") {
				ep = strings.TrimPrefix(m.ClientURLs[0], "https://")
				if _, ok := extensionClientURLMap[ep]; ok {
					extensionClientURL = "https://" + extensionClientURLMap[ep]
				}
			} else {
				ep = strings.TrimPrefix(m.ClientURLs[0], "http://")
				if _, ok := extensionClientURLMap[ep]; ok {
					extensionClientURL = "http://" + extensionClientURLMap[ep]
				}
			}
		}

		// default info
		memberVersion, memberStatus, memberRole := "", etcdv1alpha1.MemberPhaseUnStarted, etcdv1alpha1.EtcdMemberUnKnown
		var errors []string
		statusRsp, err := versioned.Status(extensionClientURL)
		if err == nil && statusRsp != nil {
			memberStatus = etcdv1alpha1.MemberPhaseRunning
			memberVersion = statusRsp.Version
			if statusRsp.IsLearner {
				memberRole = etcdv1alpha1.EtcdMemberLearner
			} else if statusRsp.Leader == m.ID {
				memberRole = etcdv1alpha1.EtcdMemberLeader
			} else {
				memberRole = etcdv1alpha1.EtcdMemberFollower
			}
		} else {
			klog.Errorf("failed to get member %s status,err is %v", extensionClientURL, err)
			errors = append(errors, err.Error())
		}

		etcdMembers = append(etcdMembers, etcdv1alpha1.MemberStatus{
			Name:               m.Name,
			MemberId:           m.ID,
			ClientUrl:          m.ClientURLs[0],
			ExtensionClientUrl: extensionClientURL,
			Role:               memberRole,
			Status:             memberStatus,
			Endpoint:           endPoint,
			Port:               items[2],
			Version:            memberVersion,
			Errors:             errors,
		})
	}

	return etcdMembers, nil
}

// GetEtcdClusterMemberStatus check healthy of cluster and member
func GetEtcdClusterMemberStatus(
	members []etcdv1alpha1.MemberStatus,
	config *etcd.ClientConfig) ([]etcdv1alpha1.MemberStatus, etcdv1alpha1.EtcdClusterPhase) {
	clusterStatus := etcdv1alpha1.EtcdClusterRunning
	newMembers := make([]etcdv1alpha1.MemberStatus, 0)
	for _, m := range members {
		healthy, err := etcd.MemberHealthy(m.ExtensionClientUrl, config)
		if err != nil {
			m.Status = etcdv1alpha1.MemberPhaseUnKnown
		} else {
			if healthy {
				m.Status = etcdv1alpha1.MemberPhaseRunning
			} else {
				m.Status = etcdv1alpha1.MemberPhaseUnHealthy
			}
		}

		if m.Status != etcdv1alpha1.MemberPhaseRunning && clusterStatus == etcdv1alpha1.EtcdClusterRunning {
			clusterStatus = etcdv1alpha1.EtcdClusterUnhealthy
		}

		newMembers = append(newMembers, m)
	}

	return newMembers, clusterStatus
}

// GetEtcdAlarms get alarm list of etcd
func GetEtcdAlarms(
	endpoints []string,
	config *etcd.ClientConfig) ([]EtcdAlarm, error) {
	etcdAlarms := make([]EtcdAlarm, 0)

	config.Endpoints = endpoints

	client, err := etcd.NewClientv3(config)
	if err != nil {
		klog.Errorf("failed to get new etcd clientv3, err is %v ", err)
		return etcdAlarms, err
	}
	defer client.Close()

	alarmRsp, err := etcd.AlarmList(client)
	if err != nil {
		klog.Errorf("failed to get alarm list, err is %v", err)
		return etcdAlarms, err
	}

	for _, a := range alarmRsp.Alarms {
		etcdAlarms = append(etcdAlarms, EtcdAlarm{
			MemberID:  a.MemberID,
			AlarmType: a.Alarm.String(),
		})
	}
	return etcdAlarms, nil
}
