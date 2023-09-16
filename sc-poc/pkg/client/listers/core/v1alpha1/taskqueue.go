/*
Copyright The Space Cloud Authors.

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
	v1alpha1 "github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// TaskQueueLister helps list TaskQueues.
// All objects returned here must be treated as read-only.
type TaskQueueLister interface {
	// List lists all TaskQueues in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.TaskQueue, err error)
	// TaskQueues returns an object that can list and get TaskQueues.
	TaskQueues(namespace string) TaskQueueNamespaceLister
	TaskQueueListerExpansion
}

// taskQueueLister implements the TaskQueueLister interface.
type taskQueueLister struct {
	indexer cache.Indexer
}

// NewTaskQueueLister returns a new TaskQueueLister.
func NewTaskQueueLister(indexer cache.Indexer) TaskQueueLister {
	return &taskQueueLister{indexer: indexer}
}

// List lists all TaskQueues in the indexer.
func (s *taskQueueLister) List(selector labels.Selector) (ret []*v1alpha1.TaskQueue, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.TaskQueue))
	})
	return ret, err
}

// TaskQueues returns an object that can list and get TaskQueues.
func (s *taskQueueLister) TaskQueues(namespace string) TaskQueueNamespaceLister {
	return taskQueueNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// TaskQueueNamespaceLister helps list and get TaskQueues.
// All objects returned here must be treated as read-only.
type TaskQueueNamespaceLister interface {
	// List lists all TaskQueues in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.TaskQueue, err error)
	// Get retrieves the TaskQueue from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.TaskQueue, error)
	TaskQueueNamespaceListerExpansion
}

// taskQueueNamespaceLister implements the TaskQueueNamespaceLister
// interface.
type taskQueueNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all TaskQueues in the indexer for a given namespace.
func (s taskQueueNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.TaskQueue, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.TaskQueue))
	})
	return ret, err
}

// Get retrieves the TaskQueue from the indexer for a given namespace and name.
func (s taskQueueNamespaceLister) Get(name string) (*v1alpha1.TaskQueue, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("taskqueue"), name)
	}
	return obj.(*v1alpha1.TaskQueue), nil
}
