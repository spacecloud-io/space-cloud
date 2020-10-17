package istio

import (
	"fmt"
)

func getServiceUniqueID(projectID, serviceID, version string) string {
	return fmt.Sprintf("%s:%s:%s", projectID, serviceID, version)
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

func getKedaScaledObjectName(serviceID, version string) string {
	return getDeploymentName(serviceID, version)
}

func getKedaTriggerAuthName(serviceID, version, triggerName string) string {
	return fmt.Sprintf("%s-%s", getDeploymentName(serviceID, version), triggerName)
}
