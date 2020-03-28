package istio

import (
	"fmt"

	"github.com/spaceuptech/space-cloud/runner/model"
)

func getServiceUniqueID(projectID, serviceID, version string) string {
	return fmt.Sprintf("%s:%s:%s", projectID, serviceID, version)
}
func getServiceUniqueName(project, service, version string) string {
	return fmt.Sprintf("%s-%s-%s", project, service, version)
}

func getServiceAccountName(serviceID string) string {
	return serviceID
}

func getDeploymentName(serviceID, version string) string {
	return fmt.Sprintf("%s-%s", serviceID, version)
}

func getServiceName(serviceID string) string {
	return serviceID
}

func getServiceDomainName(projectID, serviceID string) string {
	return fmt.Sprintf("%s.%s.svc.cluster.local", serviceID, projectID)
}
func getInternalServiceName(serviceID, version string) string {
	return fmt.Sprintf("%s-%s-internal", serviceID, version)
}

func getInternalServiceDomain(projectID, serviceID, version string) string {
	return fmt.Sprintf("%s.%s.svc.cluster.local", getInternalServiceName(serviceID, version), projectID)
}

func getVirtualServiceName(serviceID string) string {
	return serviceID
}

func getGeneralDestRuleName(serviceID string) string {
	return serviceID
}
func getInternalDestRuleName(serviceID, version string) string {
	return fmt.Sprintf("%s-%s", serviceID, version)
}

func getAuthorizationPolicyName(projectID, serviceID, version string) string {
	return fmt.Sprintf("auth-%s-%s-%s", projectID, serviceID, version)
}

func getSidecarName(serviceID, version string) string {
	return fmt.Sprintf("%s-%s", serviceID, version)
}

func getGeneratedByAnnotationName() string {
	return fmt.Sprintf("space-cloud-runner-%s", model.Version)
}
