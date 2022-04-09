package file

import "github.com/spacecloud-io/space-cloud/model"

func addResource(moduleType []*model.ResourceObject, resourceObj *model.ResourceObject) []*model.ResourceObject {
	for index, resource := range moduleType {
		if resource.Meta.Name == resourceObj.Meta.Name && matchParent(resource.Meta.Parents, resourceObj.Meta.Parents) {
			moduleType[index] = resourceObj
			return moduleType
		}
	}
	moduleType = append(moduleType, resourceObj)
	return moduleType
}

func deleteResources(moduleType []*model.ResourceObject, meta *model.ResourceMeta) []*model.ResourceObject {
	tempModuleType := make([]*model.ResourceObject, 0)
	for _, resourceObj := range moduleType {
		if !matchParent(resourceObj.Meta.Parents, meta.Parents) {
			tempModuleType = append(tempModuleType, resourceObj)
		}
	}
	return tempModuleType
}

func matchParent(resourceParents, providedParents map[string]string) bool {
	for k, v := range providedParents {
		v1, ok := resourceParents[k]
		if !ok || v1 != v {
			return false
		}
	}
	return true
}
