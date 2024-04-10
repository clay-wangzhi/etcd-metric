package providers

import (
	// register consistency inspection feature/
	_ "etcd-operator/pkg/featureprovider/providers/consistency"
	// register healthy inspection feature
	_ "etcd-operator/pkg/featureprovider/providers/healthy"
	// register request inspection feature
	_ "etcd-operator/pkg/featureprovider/providers/request"
	// register alarm inspection feature
	_ "etcd-operator/pkg/featureprovider/providers/alarm"
)
