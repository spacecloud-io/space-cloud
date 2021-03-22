package istio

import (
	"fmt"
	"strings"
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

func splitInternalServiceName(n string) (serviceID, version string) {
	arr := strings.Split(n, "-")
	if len(arr) < 3 || arr[len(arr)-1] != "internal" {
		return "", ""
	}
	return strings.Join(arr[:len(arr)-2], "-"), arr[len(arr)-2]
}

func getInternalServiceDomain(projectID, serviceID, version string) string {
	return fmt.Sprintf("%s.%s.svc.cluster.local", getInternalServiceName(serviceID, version), projectID)
}

func checkIfInternalServiceDomain(projectID, serviceID, internalServiceDomain string) bool {
	if strings.HasSuffix(internalServiceDomain, "svc.cluster.local") {
		arr := strings.Split(internalServiceDomain, ".")
		serID, _ := splitInternalServiceName(arr[0])
		if len(arr) == 5 && arr[1] == projectID && serID == serviceID {
			return true
		}
	}
	return false
}

func splitInternalServiceDomain(s string) (projectID, serviceID, version string) {
	arr := strings.Split(s, ".")
	if len(arr) < 5 {
		return "", "", ""
	}
	projectID = arr[1]
	serviceID, version = splitInternalServiceName(arr[0])
	if serviceID == "" || version == "" {
		projectID = ""
	}
	return projectID, serviceID, version
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
