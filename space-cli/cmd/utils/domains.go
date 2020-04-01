package utils

import "fmt"

// GetServiceDomain is used for getting the main service domain
func GetServiceDomain(projectID, serviceID string) string {
	return fmt.Sprintf("%s.%s.svc.cluster.local", serviceID, projectID)
}

// GetInternalServiceDomain is used for getting internal service domain
func GetInternalServiceDomain(projectID, serviceID, version string) string {
	return fmt.Sprintf("%s.%s-%s.svc.cluster.local", serviceID, projectID, version)
}
