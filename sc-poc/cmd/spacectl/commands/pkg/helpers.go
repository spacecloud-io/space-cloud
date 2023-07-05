package pkg

func findElement(arr []string, target string) int {
	for i, str := range arr {
		if str == target {
			return i
		}
	}
	return -1
}

func deleteElement(arr []string, index int) []string {
	// Check if the index is out of range
	if index < 0 || index >= len(arr) {
		return arr
	}

	// Create a new slice with the element at the specified index removed
	return append(arr[:index], arr[index+1:]...)
}
