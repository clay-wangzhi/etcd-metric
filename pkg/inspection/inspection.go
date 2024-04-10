package inspection

import (
	"context"
	"fmt"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	etcdv1alpha1 "etcd-operator/api/etcd/v1alpha1"
	"etcd-operator/pkg/controllers/util"
	"etcd-operator/pkg/etcd"
	"etcd-operator/pkg/featureprovider"
	featureutil "etcd-operator/pkg/featureprovider/util"
	clientset "etcd-operator/pkg/generated/clientset/versioned"
	platformscheme "etcd-operator/pkg/generated/clientset/versioned/scheme"
)

const (
	DefaultInspectionInterval = 300 * time.Second
	DefaultInspectionPath     = ""
)

type Server struct {
	cli                *clientset.Clientset
	kubeCli            kubernetes.Interface
	client             map[string]*clientv3.Client
	wchan              map[string]clientv3.WatchChan
	watcher            map[string]clientv3.Watcher
	eventCh            map[string]chan *clientv3.Event
	mux                sync.Mutex
	clientConfigGetter etcd.ClientConfigGetter
}

// NewInspectionServer generates the server of inspection
func NewInspectionServer(ctx *featureprovider.FeatureContext) (*Server, error) {
	cli, err := clientset.NewForConfig(ctx.ClientBuilder.ConfigOrDie())
	if err != nil {
		klog.Errorf("failed to init etcdinspection client, err is %v", err)
		return nil, err
	}
	return &Server{
		kubeCli:            ctx.ClientBuilder.ClientOrDie(),
		cli:                cli,
		client:             make(map[string]*clientv3.Client),
		wchan:              make(map[string]clientv3.WatchChan),
		watcher:            make(map[string]clientv3.Watcher),
		eventCh:            make(map[string]chan *clientv3.Event),
		clientConfigGetter: ctx.ClientConfigGetter,
	}, nil
}

// GetEtcdCluster gets etcdcluster
func (c *Server) GetEtcdCluster(namespace, name string) (*etcdv1alpha1.EtcdCluster, error) {
	return c.cli.EtcdV1alpha1().EtcdClusters(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

// GetEtcdInspection gets etcdinspection
func (c *Server) GetEtcdInspection(namespace, name string) (*etcdv1alpha1.EtcdInspection, error) {
	inspectionTask, err := c.cli.EtcdV1alpha1().EtcdInspections(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			klog.Errorf("failed to get etcdinspection, err: %v, namespace is %s, name is %s", err, namespace, name)
		}
		return nil, err
	}
	return inspectionTask, nil
}

// DeleteEtcdInspection deletes etcdinspection
func (c *Server) DeleteEtcdInspection(namespace, name string) error {
	err := c.cli.EtcdV1alpha1().EtcdInspections(namespace).
		Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		klog.Errorf(
			"failed to delete etcdinspection, namespace is %s, name is %s, err is %v",
			namespace,
			name,
			err,
		)
		return err
	}
	return nil
}

// CreateEtcdInspection creates etcdinspection
func (c *Server) CreateEtcdInspection(inspection *etcdv1alpha1.EtcdInspection) (*etcdv1alpha1.EtcdInspection, error) {
	newinspectionTask, err := c.cli.EtcdV1alpha1().EtcdInspections(inspection.Namespace).
		Create(context.TODO(), inspection, metav1.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		klog.Errorf(
			"failed to create etcdinspection, namespace is %s, name is %s, err is %v",
			inspection.Namespace,
			inspection.Name,
			err,
		)
		return newinspectionTask, err
	}
	return newinspectionTask, nil
}

func (c *Server) initInspectionTask(
	cluster *etcdv1alpha1.EtcdCluster,
	inspectionFeatureName etcdv1alpha1.KStoneFeature,
) (*etcdv1alpha1.EtcdInspection, error) {
	name := cluster.Name + "-" + string(inspectionFeatureName)
	inspectionTask := &etcdv1alpha1.EtcdInspection{}
	inspectionTask.ObjectMeta = metav1.ObjectMeta{
		Name:      name,
		Namespace: cluster.Namespace,
		Labels:    cluster.Labels,
	}
	inspectionTask.Spec = etcdv1alpha1.EtcdInspectionSpec{
		InspectionType: string(inspectionFeatureName),
		ClusterName:    cluster.Name,
	}
	inspectionTask.Status = etcdv1alpha1.EtcdInspectionStatus{
		LastUpdatedTime: metav1.Time{
			Time: time.Now(),
		},
	}

	err := controllerutil.SetOwnerReference(cluster, inspectionTask, platformscheme.Scheme)
	if err != nil {
		klog.Errorf("set inspection task's owner failed, err is %v", err)
		return inspectionTask, err
	}
	return inspectionTask, nil
}

func (c *Server) GetEtcdClusterInfo(namespace, name string) (*etcdv1alpha1.EtcdCluster, *etcd.ClientConfig, error) {
	cluster, err := c.GetEtcdCluster(namespace, name)
	if err != nil {
		klog.Errorf("failed to get cluster info, namespace is %s, name is %s, err is %v", namespace, name, err)
		return nil, nil, err
	}

	annotations := cluster.ObjectMeta.Annotations
	secretName := ""
	if annotations != nil {
		if _, found := annotations[util.ClusterTLSSecretName]; found {
			secretName = annotations[util.ClusterTLSSecretName]
		}
	}
	path := fmt.Sprintf("%s/%s", cluster.Namespace, cluster.Name)
	clientConfig, err := c.clientConfigGetter.New(path, secretName)
	if err != nil {
		klog.Errorf("failed to get cluster, namespace is %s, name is %s, err is %v", namespace, name, err)
		return nil, nil, err
	}
	return cluster, clientConfig, nil
}

func (c *Server) Equal(cluster *etcdv1alpha1.EtcdCluster, inspectionFeatureName etcdv1alpha1.KStoneFeature) bool {
	if !featureutil.IsFeatureGateEnabled(cluster.ObjectMeta.Annotations, inspectionFeatureName) {
		if cluster.Status.FeatureGatesStatus[inspectionFeatureName] != featureutil.FeatureStatusDisabled {
			return c.checkEqualIfDisabled(cluster, inspectionFeatureName)
		}
		return true
	}
	return c.checkEqualIfEnabled(cluster, inspectionFeatureName)
}

func (c *Server) Sync(cluster *etcdv1alpha1.EtcdCluster, inspectionFeatureName etcdv1alpha1.KStoneFeature) error {
	if !featureutil.IsFeatureGateEnabled(cluster.ObjectMeta.Annotations, inspectionFeatureName) {
		return c.cleanInspectionTask(cluster, inspectionFeatureName)
	}
	return c.syncInspectionTask(cluster, inspectionFeatureName)
}

// CheckEqualIfDisabled Checks whether the inspection task has been deleted if the inspection feature is disabled.
func (c *Server) checkEqualIfDisabled(cluster *etcdv1alpha1.EtcdCluster, inspectionFeatureName etcdv1alpha1.KStoneFeature) bool {
	name := cluster.Name + "-" + string(inspectionFeatureName)
	if _, err := c.GetEtcdInspection(cluster.Namespace, name); apierrors.IsNotFound(err) {
		return true
	}
	return false
}

// CheckEqualIfEnabled check whether the desired inspection task are consistent with the actual task,
// if the inspection feature is enabled.
func (c *Server) checkEqualIfEnabled(cluster *etcdv1alpha1.EtcdCluster, inspectionFeatureName etcdv1alpha1.KStoneFeature) bool {
	name := cluster.Name + "-" + string(inspectionFeatureName)
	if _, err := c.GetEtcdInspection(cluster.Namespace, name); err == nil {
		return true
	}
	return false
}

// CleanInspectionTask cleans inspection task
func (c *Server) cleanInspectionTask(cluster *etcdv1alpha1.EtcdCluster, inspectionFeatureName etcdv1alpha1.KStoneFeature) error {
	name := cluster.Name + "-" + string(inspectionFeatureName)
	return c.DeleteEtcdInspection(cluster.Namespace, name)
}

// SyncInspectionTask syncs inspection task
func (c *Server) syncInspectionTask(cluster *etcdv1alpha1.EtcdCluster, inspectionFeatureName etcdv1alpha1.KStoneFeature) error {
	task, err := c.initInspectionTask(cluster, inspectionFeatureName)
	if err != nil {
		return err
	}
	_, err = c.CreateEtcdInspection(task)
	if err != nil {
		return err
	}
	return nil
}
