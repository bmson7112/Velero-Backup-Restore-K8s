/*
Copyright the Velero contributors.

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

package v2alpha1

import (
	"context"
	time "time"

	velerov2alpha1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v2alpha1"
	versioned "github.com/vmware-tanzu/velero/pkg/generated/clientset/versioned"
	internalinterfaces "github.com/vmware-tanzu/velero/pkg/generated/informers/externalversions/internalinterfaces"
	v2alpha1 "github.com/vmware-tanzu/velero/pkg/generated/listers/velero/v2alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// DataDownloadInformer provides access to a shared informer and lister for
// DataDownloads.
type DataDownloadInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v2alpha1.DataDownloadLister
}

type dataDownloadInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewDataDownloadInformer constructs a new informer for DataDownload type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewDataDownloadInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredDataDownloadInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredDataDownloadInformer constructs a new informer for DataDownload type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredDataDownloadInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.VeleroV2alpha1().DataDownloads(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.VeleroV2alpha1().DataDownloads(namespace).Watch(context.TODO(), options)
			},
		},
		&velerov2alpha1.DataDownload{},
		resyncPeriod,
		indexers,
	)
}

func (f *dataDownloadInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredDataDownloadInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *dataDownloadInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&velerov2alpha1.DataDownload{}, f.defaultInformer)
}

func (f *dataDownloadInformer) Lister() v2alpha1.DataDownloadLister {
	return v2alpha1.NewDataDownloadLister(f.Informer().GetIndexer())
}
