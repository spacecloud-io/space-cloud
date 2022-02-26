package connectors

import (
	"github.com/spacecloud-io/space-cloud/config"
)

func areSchemasSimilar(a, b config.DatabaseSchemas) bool {
	// Simply return false if length of arrays are not similar
	if len(a) != len(b) {
		return false
	}

	// Check if all keys of first map are present in the second map
	for k, s1 := range a {
		s2, p := b[k]
		if !p {
			return false
		}

		if s1.Table != s2.Table {
			return false
		}
	}

	return true
}

func sanitizePrepareQueries(queries config.DatabasePreparedQueries) config.DatabasePreparedQueries {
	temp := make(config.DatabasePreparedQueries, len(queries))

	for _, v := range queries {
		temp[v.ID] = v
	}

	return temp
}
