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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	corev1alpha1 "github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
	versioned "github.com/spacecloud-io/space-cloud/pkg/client/clientset/versioned"
	internalinterfaces "github.com/spacecloud-io/space-cloud/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/spacecloud-io/space-cloud/pkg/client/listers/core/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// PubsubChannelInformer provides access to a shared informer and lister for
// PubsubChannels.
type PubsubChannelInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.PubsubChannelLister
}

type pubsubChannelInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewPubsubChannelInformer constructs a new informer for PubsubChannel type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewPubsubChannelInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredPubsubChannelInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredPubsubChannelInformer constructs a new informer for PubsubChannel type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredPubsubChannelInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CoreV1alpha1().PubsubChannels(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CoreV1alpha1().PubsubChannels(namespace).Watch(context.TODO(), options)
			},
		},
		&corev1alpha1.PubsubChannel{},
		resyncPeriod,
		indexers,
	)
}

func (f *pubsubChannelInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredPubsubChannelInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *pubsubChannelInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&corev1alpha1.PubsubChannel{}, f.defaultInformer)
}

func (f *pubsubChannelInformer) Lister() v1alpha1.PubsubChannelLister {
	return v1alpha1.NewPubsubChannelLister(f.Informer().GetIndexer())
}
