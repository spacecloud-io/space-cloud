package schema

func (s *Schema) CheckIfEventingIsPossible(dbAlias, col string, find map[string]interface{}) (findForUpdate map[string]interface{}, present bool) {
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
				t := &trackCols{find: map[string]interface{}{}}
				tracker["i:"+fieldSchema.IndexInfo.Group] = t
			}

			// Add count for want
			t.want++

			// Check if field is present in the find clause
			if value, ok := isFieldPresentInFindAndIsValidForEventing(fieldName, find); ok {
				t.find[fieldName] = value
			}
		}

		// Process for primary key
		if fieldSchema.IsPrimary {
			t, p := tracker["p"]
			if !p {
				t := &trackCols{}
				tracker["p"] = t
			}

			// Add count for want
			t.want++

			// Check if field is present in the find clause
			if value, ok := isFieldPresentInFindAndIsValidForEventing(fieldName, find); ok {
				find[fieldName] = value
			}
		}
	}

	for _, t := range tracker {
		if t.want == len(t.find) {
			return t.find, true
		}
	}

	return nil, false
}

func isFieldPresentInFindAndIsValidForEventing(fieldName string, find map[string]interface{}) (interface{}, bool) {
	if findValue, p := find[fieldName]; p {
		findValueObj, ok := findValue.(map[string]interface{})
		if !ok {
			return findValue, true
		}

		if findValue, p := findValueObj["$eq"]; p {
			return findValue, true
		}
	}
	return nil, false
}
