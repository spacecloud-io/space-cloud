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
	v1alpha1 "github.com/spacecloud-io/space-cloud/pkg/client/clientset/versioned/typed/core/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeCoreV1alpha1 struct {
	*testing.Fake
}

func (c *FakeCoreV1alpha1) CompiledGraphqlSources(namespace string) v1alpha1.CompiledGraphqlSourceInterface {
	return &FakeCompiledGraphqlSources{c, namespace}
}

func (c *FakeCoreV1alpha1) GraphqlSources(namespace string) v1alpha1.GraphqlSourceInterface {
	return &FakeGraphqlSources{c, namespace}
}

func (c *FakeCoreV1alpha1) JwtHSASecrets(namespace string) v1alpha1.JwtHSASecretInterface {
	return &FakeJwtHSASecrets{c, namespace}
}

func (c *FakeCoreV1alpha1) OPAPolicies(namespace string) v1alpha1.OPAPolicyInterface {
	return &FakeOPAPolicies{c, namespace}
}

func (c *FakeCoreV1alpha1) OpenAPISources(namespace string) v1alpha1.OpenAPISourceInterface {
	return &FakeOpenAPISources{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeCoreV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
