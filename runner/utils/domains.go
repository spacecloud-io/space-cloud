package utils

import "fmt"

func GetServiceDomain(projectID, serviceID string) string {
	return fmt.Sprintf("%s.%s.svc.cluster.local", serviceID, projectID)
}

func GetInternalServiceDomain(projectID, serviceID, version string) string {
	return fmt.Sprintf("%s.%s.%s.svc.cluster.local", serviceID, projectID, version)
}
