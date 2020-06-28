package schema

// CheckIfEventingIsPossible checks if eventing is possible
func (s *Schema) CheckIfEventingIsPossible(dbAlias, col string, obj map[string]interface{}, isFind bool) (findForUpdate map[string]interface{}, present bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	// Struct to track counts
	type trackCols struct {
		want int
		find map[string]interface{}
	}

	dbSchema, p := s.SchemaDoc[dbAlias]
	if !p {
		return nil, false
	}

	colSchema, p := dbSchema[col]
	if !p {
		return nil, false
	}

	tracker := map[string]*trackCols{}
	for fieldName, fieldSchema := range colSchema {
		// Process for unique index
		if fieldSchema.IsIndex && fieldSchema.IsUnique {
			t, p := tracker["i:"+fieldSchema.IndexInfo.Group]
			if !p {
				t = &trackCols{find: map[string]interface{}{}}
				tracker["i:"+fieldSchema.IndexInfo.Group] = t
			}

			// Add count for want
			t.want++

			// Check if field is present in the find clause
			if value, ok := isFieldPresentInFindAndIsValidForEventing(fieldName, obj, isFind); ok {
				t.find[fieldName] = value
			}
		}

		// Process for primary key
		if fieldSchema.IsPrimary {
			t, p := tracker["p"]
			if !p {
				t = &trackCols{find: map[string]interface{}{}}
				tracker["p"] = t
			}

			// Add count for want
			t.want++

			// Check if field is present in the find clause
			if value, ok := isFieldPresentInFindAndIsValidForEventing(fieldName, obj, isFind); ok {
				t.find[fieldName] = value
			}
		}
	}

	// First check if the primary key was provided
	if t, p := tracker["p"]; p && t.want == len(t.find) {
		return t.find, true
	}

	for _, t := range tracker {
		if t.want == len(t.find) {
			return t.find, true
		}
	}

	return nil, false
}

func isFieldPresentInFindAndIsValidForEventing(fieldName string, obj map[string]interface{}, isFind bool) (interface{}, bool) {
	if findValue, p := obj[fieldName]; p {
		findValueObj, ok := findValue.(map[string]interface{})
		if !ok {
			return findValue, true
		}

		if !isFind {
			return findValue, true
		}

		if findValue, p := findValueObj["$eq"]; p {
			return findValue, true
		}
	}
	return nil, false
}
