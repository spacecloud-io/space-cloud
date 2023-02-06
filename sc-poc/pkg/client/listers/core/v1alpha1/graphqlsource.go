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

// GraphqlSourceLister helps list GraphqlSources.
// All objects returned here must be treated as read-only.
type GraphqlSourceLister interface {
	// List lists all GraphqlSources in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.GraphqlSource, err error)
	// GraphqlSources returns an object that can list and get GraphqlSources.
	GraphqlSources(namespace string) GraphqlSourceNamespaceLister
	GraphqlSourceListerExpansion
}

// graphqlSourceLister implements the GraphqlSourceLister interface.
type graphqlSourceLister struct {
	indexer cache.Indexer
}

// NewGraphqlSourceLister returns a new GraphqlSourceLister.
func NewGraphqlSourceLister(indexer cache.Indexer) GraphqlSourceLister {
	return &graphqlSourceLister{indexer: indexer}
}

// List lists all GraphqlSources in the indexer.
func (s *graphqlSourceLister) List(selector labels.Selector) (ret []*v1alpha1.GraphqlSource, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.GraphqlSource))
	})
	return ret, err
}

// GraphqlSources returns an object that can list and get GraphqlSources.
func (s *graphqlSourceLister) GraphqlSources(namespace string) GraphqlSourceNamespaceLister {
	return graphqlSourceNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// GraphqlSourceNamespaceLister helps list and get GraphqlSources.
// All objects returned here must be treated as read-only.
type GraphqlSourceNamespaceLister interface {
	// List lists all GraphqlSources in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.GraphqlSource, err error)
	// Get retrieves the GraphqlSource from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.GraphqlSource, error)
	GraphqlSourceNamespaceListerExpansion
}

// graphqlSourceNamespaceLister implements the GraphqlSourceNamespaceLister
// interface.
type graphqlSourceNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all GraphqlSources in the indexer for a given namespace.
func (s graphqlSourceNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.GraphqlSource, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.GraphqlSource))
	})
	return ret, err
}

// Get retrieves the GraphqlSource from the indexer for a given namespace and name.
func (s graphqlSourceNamespaceLister) Get(name string) (*v1alpha1.GraphqlSource, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("graphqlsource"), name)
	}
	return obj.(*v1alpha1.GraphqlSource), nil
}
