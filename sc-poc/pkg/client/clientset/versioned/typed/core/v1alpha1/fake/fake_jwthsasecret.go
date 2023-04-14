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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeJwtHSASecrets implements JwtHSASecretInterface
type FakeJwtHSASecrets struct {
	Fake *FakeCoreV1alpha1
	ns   string
}

var jwthsasecretsResource = schema.GroupVersionResource{Group: "core.space-cloud.io", Version: "v1alpha1", Resource: "jwthsasecrets"}

var jwthsasecretsKind = schema.GroupVersionKind{Group: "core.space-cloud.io", Version: "v1alpha1", Kind: "JwtHSASecret"}

// Get takes name of the jwtHSASecret, and returns the corresponding jwtHSASecret object, and an error if there is any.
func (c *FakeJwtHSASecrets) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.JwtHSASecret, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(jwthsasecretsResource, c.ns, name), &v1alpha1.JwtHSASecret{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.JwtHSASecret), err
}

// List takes label and field selectors, and returns the list of JwtHSASecrets that match those selectors.
func (c *FakeJwtHSASecrets) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.JwtHSASecretList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(jwthsasecretsResource, jwthsasecretsKind, c.ns, opts), &v1alpha1.JwtHSASecretList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.JwtHSASecretList{ListMeta: obj.(*v1alpha1.JwtHSASecretList).ListMeta}
	for _, item := range obj.(*v1alpha1.JwtHSASecretList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested jwtHSASecrets.
func (c *FakeJwtHSASecrets) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(jwthsasecretsResource, c.ns, opts))

}

// Create takes the representation of a jwtHSASecret and creates it.  Returns the server's representation of the jwtHSASecret, and an error, if there is any.
func (c *FakeJwtHSASecrets) Create(ctx context.Context, jwtHSASecret *v1alpha1.JwtHSASecret, opts v1.CreateOptions) (result *v1alpha1.JwtHSASecret, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(jwthsasecretsResource, c.ns, jwtHSASecret), &v1alpha1.JwtHSASecret{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.JwtHSASecret), err
}

// Update takes the representation of a jwtHSASecret and updates it. Returns the server's representation of the jwtHSASecret, and an error, if there is any.
func (c *FakeJwtHSASecrets) Update(ctx context.Context, jwtHSASecret *v1alpha1.JwtHSASecret, opts v1.UpdateOptions) (result *v1alpha1.JwtHSASecret, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(jwthsasecretsResource, c.ns, jwtHSASecret), &v1alpha1.JwtHSASecret{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.JwtHSASecret), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeJwtHSASecrets) UpdateStatus(ctx context.Context, jwtHSASecret *v1alpha1.JwtHSASecret, opts v1.UpdateOptions) (*v1alpha1.JwtHSASecret, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(jwthsasecretsResource, "status", c.ns, jwtHSASecret), &v1alpha1.JwtHSASecret{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.JwtHSASecret), err
}

// Delete takes name of the jwtHSASecret and deletes it. Returns an error if one occurs.
func (c *FakeJwtHSASecrets) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(jwthsasecretsResource, c.ns, name, opts), &v1alpha1.JwtHSASecret{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeJwtHSASecrets) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(jwthsasecretsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.JwtHSASecretList{})
	return err
}

// Patch applies the patch and returns the patched jwtHSASecret.
func (c *FakeJwtHSASecrets) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.JwtHSASecret, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(jwthsasecretsResource, c.ns, name, pt, data, subresources...), &v1alpha1.JwtHSASecret{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.JwtHSASecret), err
}