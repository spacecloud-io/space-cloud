package typescript

import (
	"strings"
)

func getTypeName(name string, skipFirst bool) string {
	arr := strings.Split(name, "-")
	for i, item := range arr {
		if i == 0 && skipFirst {
			arr[i] = item
			continue
		}

		arr[i] = strings.Title(item)
	}

	return strings.Join(arr, "")
}
