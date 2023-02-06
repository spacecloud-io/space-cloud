package utils

// MergeMaps merges the overlay map on top of the base map
func MergeMaps(base, overlay map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(base))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range overlay {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = MergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}

// ShallowCopy returns the copied map
func ShallowCopy[V any](m map[string]V) map[string]V {
	newMap := make(map[string]V, len(m))
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}
