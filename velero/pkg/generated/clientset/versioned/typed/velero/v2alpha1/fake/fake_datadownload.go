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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v2alpha1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v2alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeDataDownloads implements DataDownloadInterface
type FakeDataDownloads struct {
	Fake *FakeVeleroV2alpha1
	ns   string
}

var datadownloadsResource = schema.GroupVersionResource{Group: "velero.io", Version: "v2alpha1", Resource: "datadownloads"}

var datadownloadsKind = schema.GroupVersionKind{Group: "velero.io", Version: "v2alpha1", Kind: "DataDownload"}

// Get takes name of the dataDownload, and returns the corresponding dataDownload object, and an error if there is any.
func (c *FakeDataDownloads) Get(ctx context.Context, name string, options v1.GetOptions) (result *v2alpha1.DataDownload, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(datadownloadsResource, c.ns, name), &v2alpha1.DataDownload{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2alpha1.DataDownload), err
}

// List takes label and field selectors, and returns the list of DataDownloads that match those selectors.
func (c *FakeDataDownloads) List(ctx context.Context, opts v1.ListOptions) (result *v2alpha1.DataDownloadList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(datadownloadsResource, datadownloadsKind, c.ns, opts), &v2alpha1.DataDownloadList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v2alpha1.DataDownloadList{ListMeta: obj.(*v2alpha1.DataDownloadList).ListMeta}
	for _, item := range obj.(*v2alpha1.DataDownloadList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested dataDownloads.
func (c *FakeDataDownloads) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(datadownloadsResource, c.ns, opts))

}

// Create takes the representation of a dataDownload and creates it.  Returns the server's representation of the dataDownload, and an error, if there is any.
func (c *FakeDataDownloads) Create(ctx context.Context, dataDownload *v2alpha1.DataDownload, opts v1.CreateOptions) (result *v2alpha1.DataDownload, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(datadownloadsResource, c.ns, dataDownload), &v2alpha1.DataDownload{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2alpha1.DataDownload), err
}

// Update takes the representation of a dataDownload and updates it. Returns the server's representation of the dataDownload, and an error, if there is any.
func (c *FakeDataDownloads) Update(ctx context.Context, dataDownload *v2alpha1.DataDownload, opts v1.UpdateOptions) (result *v2alpha1.DataDownload, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(datadownloadsResource, c.ns, dataDownload), &v2alpha1.DataDownload{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2alpha1.DataDownload), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeDataDownloads) UpdateStatus(ctx context.Context, dataDownload *v2alpha1.DataDownload, opts v1.UpdateOptions) (*v2alpha1.DataDownload, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(datadownloadsResource, "status", c.ns, dataDownload), &v2alpha1.DataDownload{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2alpha1.DataDownload), err
}

// Delete takes name of the dataDownload and deletes it. Returns an error if one occurs.
func (c *FakeDataDownloads) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(datadownloadsResource, c.ns, name), &v2alpha1.DataDownload{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeDataDownloads) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(datadownloadsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v2alpha1.DataDownloadList{})
	return err
}

// Patch applies the patch and returns the patched dataDownload.
func (c *FakeDataDownloads) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v2alpha1.DataDownload, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(datadownloadsResource, c.ns, name, pt, data, subresources...), &v2alpha1.DataDownload{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2alpha1.DataDownload), err
}
