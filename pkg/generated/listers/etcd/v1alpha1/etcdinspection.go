/*
Copyright 2023 Clay.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "etcd-operator/api/etcd/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// EtcdInspectionLister helps list EtcdInspections.
// All objects returned here must be treated as read-only.
type EtcdInspectionLister interface {
	// List lists all EtcdInspections in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.EtcdInspection, err error)
	// EtcdInspections returns an object that can list and get EtcdInspections.
	EtcdInspections(namespace string) EtcdInspectionNamespaceLister
	EtcdInspectionListerExpansion
}

// etcdInspectionLister implements the EtcdInspectionLister interface.
type etcdInspectionLister struct {
	indexer cache.Indexer
}

// NewEtcdInspectionLister returns a new EtcdInspectionLister.
func NewEtcdInspectionLister(indexer cache.Indexer) EtcdInspectionLister {
	return &etcdInspectionLister{indexer: indexer}
}

// List lists all EtcdInspections in the indexer.
func (s *etcdInspectionLister) List(selector labels.Selector) (ret []*v1alpha1.EtcdInspection, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.EtcdInspection))
	})
	return ret, err
}

// EtcdInspections returns an object that can list and get EtcdInspections.
func (s *etcdInspectionLister) EtcdInspections(namespace string) EtcdInspectionNamespaceLister {
	return etcdInspectionNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// EtcdInspectionNamespaceLister helps list and get EtcdInspections.
// All objects returned here must be treated as read-only.
type EtcdInspectionNamespaceLister interface {
	// List lists all EtcdInspections in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.EtcdInspection, err error)
	// Get retrieves the EtcdInspection from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.EtcdInspection, error)
	EtcdInspectionNamespaceListerExpansion
}

// etcdInspectionNamespaceLister implements the EtcdInspectionNamespaceLister
// interface.
type etcdInspectionNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all EtcdInspections in the indexer for a given namespace.
func (s etcdInspectionNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.EtcdInspection, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.EtcdInspection))
	})
	return ret, err
}

// Get retrieves the EtcdInspection from the indexer for a given namespace and name.
func (s etcdInspectionNamespaceLister) Get(name string) (*v1alpha1.EtcdInspection, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("etcdinspection"), name)
	}
	return obj.(*v1alpha1.EtcdInspection), nil
}
