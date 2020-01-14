package istio

import (
	"fmt"

	"github.com/spaceuptech/space-cloud/runner/model"
)

func getNamespaceName(project, env string) string {
	return fmt.Sprintf("%s-%s", project, env)
}

func getServiceUniqueName(project, service, environment, version string) string {
	return fmt.Sprintf("%s-%s-%s-%s", project, service, environment, version)
}

func getServiceAccountName(service *model.Service) string {
	return fmt.Sprintf("%s-%s", service.ProjectID, service.ID)
}

func getDeploymentName(service *model.Service) string {
	return fmt.Sprintf("%s-%s", service.ID, service.Version)
}

func getAuthorizationPolicyName(service *model.Service) string {
	return fmt.Sprintf("auth-%s-%s", service.ProjectID, service.ID)
}

func getGatewayName(service *model.Service) string {
	return fmt.Sprintf("gateway-%s", service.ID)
}
