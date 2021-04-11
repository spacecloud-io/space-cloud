package utils

import (
	"fmt"
	"strings"
)

// GetHelmChartDownloadURL adjusts the url prefixes according to the version
func GetHelmChartDownloadURL(url, version string) string {
	arr := strings.Split(url, "/")
	chartName := fmt.Sprintf("%s-%s.tgz", arr[len(arr)-1], version)
	arr = append(arr, chartName)
	return strings.Join(arr, "/")
}
