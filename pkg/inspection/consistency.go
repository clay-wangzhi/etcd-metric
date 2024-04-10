package inspection

import (
	"context"
	"sort"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
	"k8s.io/klog/v2"

	etcdv1alpha1 "etcd-operator/api/etcd/v1alpha1"
	"etcd-operator/pkg/etcd"
	featureutil "etcd-operator/pkg/featureprovider/util"
	"etcd-operator/pkg/inspection/metrics"
)

type ConsistencyInfo struct {
	Path     string `json:"path,omitempty"`
	Interval int    `json:"interval,omitempty"`
}

// getEtcdConsistentMetadata gets the etcd consistent metadata.
func (c *Server) getEtcdConsistentMetadata(
	cluster *etcdv1alpha1.EtcdCluster,
	keyPrefix string,
	cli *etcd.ClientConfig,
) (map[featureutil.ConsistencyType][]uint64, error) {

	var mu sync.Mutex
	endpointMetadata := make(map[featureutil.ConsistencyType][]uint64)

	ctx, cancel := context.WithTimeout(context.Background(), etcd.DefaultDialTimeout)
	g, ctx := errgroup.WithContext(ctx)
	defer cancel()

	for _, member := range cluster.Status.Members {
		member := member
		g.Go(func() error {
			backendStorage, extensionClientURL := etcd.EtcdV3Backend, member.ExtensionClientUrl
			if strings.HasPrefix(member.Version, "2") {
				backendStorage = etcd.EtcdV2Backend
			}
			backend, err := etcd.NewEtcdStatBackend(backendStorage)
			if err != nil {
				klog.Errorf("failed to get etcd stat backend,backend %s,err is %v", extensionClientURL, err)
				return err
			}

			cli.Endpoints = []string{extensionClientURL}

			err = backend.Init(cli)
			if err != nil {
				klog.Errorf(
					"failed to get new etcd clientv3,cluster name is %s,endpoint is %s,err is %v",
					cluster.Name,
					extensionClientURL,
					err,
				)
				return err
			}
			defer backend.Close()
			metadata, err := backend.GetIndex(ctx, extensionClientURL)
			if err != nil {
				return err
			}
			totalKey, err := backend.GetTotalKeyNum(ctx, keyPrefix)
			if err != nil {
				return err
			}
			metadata[featureutil.ConsistencyKeyTotal] = totalKey

			mu.Lock()
			for t, v := range metadata {
				endpointMetadata[t] = append(endpointMetadata[t], v)
			}
			mu.Unlock()
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return endpointMetadata, err
	}
	return endpointMetadata, nil
}

// CollectClusterConsistentData collects the etcd metadata info, calculate the difference, and
// transfer them to prometheus metrics
func (c *Server) CollectClusterConsistentData(inspection *etcdv1alpha1.EtcdInspection) error {
	namespace, name := inspection.Namespace, inspection.Spec.ClusterName
	cluster, clientConfig, err := c.GetEtcdClusterInfo(namespace, name)
	defer func() {
		if err != nil {
			featureutil.IncrFailedInspectionCounter(name, etcdv1alpha1.KStoneFeatureConsistency)
		}
	}()
	if err != nil {
		klog.Errorf("failed to load tls config, namespace is %s, name is %s, err is %v", namespace, name, err)
		return err
	}
	endpointMetadataDiff := make(map[featureutil.ConsistencyType]uint64)
	endpointMetadata, err := c.getEtcdConsistentMetadata(cluster, DefaultInspectionPath, clientConfig)
	if err != nil {
		klog.Errorf("failed to getEtcdConsistentMetadata, etcd cluster %s, err is %v", cluster.Name, err)
	} else {
		for t, values := range endpointMetadata {
			sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })
			endpointMetadataDiff[t] = values[len(values)-1] - values[0]
		}
	}
	labels := map[string]string{
		"clusterName": cluster.Name,
	}
	for t, v := range endpointMetadataDiff {
		switch t {
		case featureutil.ConsistencyKeyTotal:
			metrics.EtcdNodeDiffTotal.With(labels).Set(float64(v))
		case featureutil.ConsistencyRevision:
			metrics.EtcdNodeRevisionDiff.With(labels).Set(float64(v))
		case featureutil.ConsistencyIndex:
			metrics.EtcdNodeIndexDiff.With(labels).Set(float64(v))
		case featureutil.ConsistencyRaftRaftAppliedIndex:
			metrics.EtcdNodeRaftAppliedIndexDiff.With(labels).Set(float64(v))
		case featureutil.ConsistencyRaftIndex:
			metrics.EtcdNodeRaftIndexDiff.With(labels).Set(float64(v))
		}
	}
	return nil
}
