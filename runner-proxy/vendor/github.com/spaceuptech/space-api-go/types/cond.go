package types

// M is a type representing a map
type M map[string]interface{}

// Cond is a function to write a condition
func Cond(f1, eval string, f2 interface{}) M {
	return M{"type": "cond", "f1": f1, "eval": eval, "f2": f2}
}

// And is a function to "and" multiple conditions together
func And(conds ...M) M {
	return M{"type": "and", "conds": conds}
}

// Or is a function to "or" multiple conditions together
func Or(conds ...M) M {
	return M{"type": "or", "conds": conds}
}

// GenerateFind generates a Mongo db find clause from the provided condition
func GenerateFind(condition M) M {
	m := M{}
	switch condition["type"].(string) {
	case "and":
		conds := condition["conds"].([]M)
		for _, c := range conds {
			t := GenerateFind(c)
			if typ, ok := c["type"]; ok && typ == "cond" {
				if f1, ok := m[c["f1"].(string)]; ok {
					for k, v := range t[c["f1"].(string)].(map[string]interface{}) {
						f1.(map[string]interface{})[k] = v
					}
				} else {
					for k, v := range t {
						m[k] = v
					}
				}
			} else {
				for k, v := range t {
					m[k] = v
				}
			}
		}

	case "or":
		conds := condition["conds"].([]M)
		t := []M{}
		for _, c := range conds {
			t = append(t, GenerateFind(c))
		}
		m["$or"] = t

	case "cond":
		f1 := condition["f1"].(string)
		eval := condition["eval"].(string)
		f2 := condition["f2"]

		switch eval {
		case "==":
			m[f1] = map[string]interface{}{"$eq": f2}
		case "!=":
			m[f1] = map[string]interface{}{"$ne": f2}
		case ">":
			m[f1] = map[string]interface{}{"$gt": f2}
		case "<":
			m[f1] = map[string]interface{}{"$lt": f2}
		case ">=":
			m[f1] = map[string]interface{}{"$gte": f2}
		case "<=":
			m[f1] = map[string]interface{}{"$lte": f2}
		case "in":
			m[f1] = map[string]interface{}{"$in": f2}
		case "notIn":
			m[f1] = map[string]interface{}{"$nin": f2}
		case "regex":
			m[f1] = map[string]interface{}{"$regex": f2}
		case "contains":
			m[f1] = map[string]interface{}{"$contains": f2}
		}
	}

	return m
}
