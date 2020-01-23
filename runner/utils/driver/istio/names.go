package istio

import (
	"fmt"

	"github.com/spaceuptech/space-cloud/runner/model"
)

func getServiceUniqueName(project, service, version string) string {
	return fmt.Sprintf("%s-%s-%s", project, service, version)
}

func getServiceAccountName(service *model.Service) string {
	return fmt.Sprintf("%s-%s", service.ProjectID, service.ID)
}

func getDeploymentName(service *model.Service) string {
	return fmt.Sprintf("%s-%s", service.ID, service.Version)
}

func getServiceName(serviceID string) string {
	return serviceID
}

func getVirtualServiceName(serviceID string) string {
	return serviceID
}

func getDestRuleName(serviceID string) string {
	return serviceID
}

func getAuthorizationPolicyName(service *model.Service) string {
	return fmt.Sprintf("auth-%s-%s", service.ProjectID, service.ID)
}

func getSidecarName(serviceID string) string {
	return serviceID
}
