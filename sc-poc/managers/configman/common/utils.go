package common

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

func ConvertToList(sources []*unstructured.Unstructured) *unstructured.UnstructuredList {
	list := &unstructured.UnstructuredList{}
	if len(sources) == 0 {
		return list
	}

	kind := sources[0].GetKind() + "List"
	apiVersion := sources[0].GetAPIVersion()

	list.SetAPIVersion(apiVersion)
	list.SetKind(kind)
	list.Items = []unstructured.Unstructured{}

	for _, src := range sources {
		list.Items = append(list.Items, *src)
	}

	return list
}
