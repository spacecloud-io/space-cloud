package k8s

import (
	"context"

	"github.com/spacecloud-io/space-cloud/managers/source"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8s) loadConfiguration() error {

	// CompiledGraphQLSource
	compiledGraphqlSourceList, err := k.dc.Resource(compiledgraphqlsourcesResource).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, obj := range compiledGraphqlSourceList.Items {
		kind := obj.GetKind()
		key := source.GetModuleName(obj.GetAPIVersion(), obj.GetKind())

		k.configuration[kind] = append(k.configuration[kind], &obj)
		k.configurationN[key] = append(k.configurationN[key], &obj)
	}

	// GraphQLSource
	graphqlSourceList, err := k.dc.Resource(graphqlsourcesResource).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, obj := range graphqlSourceList.Items {
		kind := obj.GetKind()
		key := source.GetModuleName(obj.GetAPIVersion(), obj.GetKind())

		k.configuration[kind] = append(k.configuration[kind], &obj)
		k.configurationN[key] = append(k.configurationN[key], &obj)
	}

	// JWTHSASecrets
	jwtHSASecretsList, err := k.dc.Resource(jwthsasecretsResource).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, obj := range jwtHSASecretsList.Items {
		kind := obj.GetKind()
		key := source.GetModuleName(obj.GetAPIVersion(), obj.GetKind())

		k.configuration[kind] = append(k.configuration[kind], &obj)
		k.configurationN[key] = append(k.configurationN[key], &obj)
	}

	// OPAPolicies
	opaPoliciesList, err := k.dc.Resource(opapoliciesResource).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, obj := range opaPoliciesList.Items {
		kind := obj.GetKind()
		key := source.GetModuleName(obj.GetAPIVersion(), obj.GetKind())

		k.configuration[kind] = append(k.configuration[kind], &obj)
		k.configurationN[key] = append(k.configurationN[key], &obj)
	}

	return nil
}
