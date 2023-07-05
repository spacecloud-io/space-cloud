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

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/spacecloud-io/space-cloud/pkg/apis/core/v1alpha1"
	scheme "github.com/spacecloud-io/space-cloud/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// JwtRSASecretsGetter has a method to return a JwtRSASecretInterface.
// A group's client should implement this interface.
type JwtRSASecretsGetter interface {
	JwtRSASecrets(namespace string) JwtRSASecretInterface
}

// JwtRSASecretInterface has methods to work with JwtRSASecret resources.
type JwtRSASecretInterface interface {
	Create(ctx context.Context, jwtRSASecret *v1alpha1.JwtRSASecret, opts v1.CreateOptions) (*v1alpha1.JwtRSASecret, error)
	Update(ctx context.Context, jwtRSASecret *v1alpha1.JwtRSASecret, opts v1.UpdateOptions) (*v1alpha1.JwtRSASecret, error)
	UpdateStatus(ctx context.Context, jwtRSASecret *v1alpha1.JwtRSASecret, opts v1.UpdateOptions) (*v1alpha1.JwtRSASecret, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.JwtRSASecret, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.JwtRSASecretList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.JwtRSASecret, err error)
	JwtRSASecretExpansion
}

// jwtRSASecrets implements JwtRSASecretInterface
type jwtRSASecrets struct {
	client rest.Interface
	ns     string
}

// newJwtRSASecrets returns a JwtRSASecrets
func newJwtRSASecrets(c *CoreV1alpha1Client, namespace string) *jwtRSASecrets {
	return &jwtRSASecrets{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the jwtRSASecret, and returns the corresponding jwtRSASecret object, and an error if there is any.
func (c *jwtRSASecrets) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.JwtRSASecret, err error) {
	result = &v1alpha1.JwtRSASecret{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("jwtrsasecrets").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of JwtRSASecrets that match those selectors.
func (c *jwtRSASecrets) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.JwtRSASecretList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.JwtRSASecretList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("jwtrsasecrets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested jwtRSASecrets.
func (c *jwtRSASecrets) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("jwtrsasecrets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a jwtRSASecret and creates it.  Returns the server's representation of the jwtRSASecret, and an error, if there is any.
func (c *jwtRSASecrets) Create(ctx context.Context, jwtRSASecret *v1alpha1.JwtRSASecret, opts v1.CreateOptions) (result *v1alpha1.JwtRSASecret, err error) {
	result = &v1alpha1.JwtRSASecret{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("jwtrsasecrets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(jwtRSASecret).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a jwtRSASecret and updates it. Returns the server's representation of the jwtRSASecret, and an error, if there is any.
func (c *jwtRSASecrets) Update(ctx context.Context, jwtRSASecret *v1alpha1.JwtRSASecret, opts v1.UpdateOptions) (result *v1alpha1.JwtRSASecret, err error) {
	result = &v1alpha1.JwtRSASecret{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("jwtrsasecrets").
		Name(jwtRSASecret.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(jwtRSASecret).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *jwtRSASecrets) UpdateStatus(ctx context.Context, jwtRSASecret *v1alpha1.JwtRSASecret, opts v1.UpdateOptions) (result *v1alpha1.JwtRSASecret, err error) {
	result = &v1alpha1.JwtRSASecret{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("jwtrsasecrets").
		Name(jwtRSASecret.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(jwtRSASecret).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the jwtRSASecret and deletes it. Returns an error if one occurs.
func (c *jwtRSASecrets) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("jwtrsasecrets").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *jwtRSASecrets) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("jwtrsasecrets").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched jwtRSASecret.
func (c *jwtRSASecrets) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.JwtRSASecret, err error) {
	result = &v1alpha1.JwtRSASecret{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("jwtrsasecrets").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}