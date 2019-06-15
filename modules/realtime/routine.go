package realtime

import (
	"reflect"
	"sync"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (m *Module) initWorkers(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go m.worker()
	}
}

func (m *Module) worker() {
	if !m.enabled {
		return
	}

	for data := range m.feed {
		clientsTemp, ok := m.groups.Load(data.Group)
		if !ok {
			continue
		}
		clients := clientsTemp.(*sync.Map)
		clients.Range(func(key interface{}, value interface{}) bool {
			queries := value.(*sync.Map)
			queries.Range(func(id interface{}, value interface{}) bool {
				query := value.(*queryStub)

				// Send feed data if type is delete or the where clause matches
				if data.Type == utils.RealtimeDelete || validate(query.whereObj, data.Payload) {
					dataPoint := &model.FeedData{
						QueryID: id.(string), DocID: data.DocID, Group: data.Group, Payload: data.Payload,
						TimeStamp: data.TimeStamp, Type: data.Type, DBType: data.DBType,
					}
					query.sendFeed(dataPoint)
				}
				return true
			})
			return true
		})
	}
}

func validate(where map[string]interface{}, obj interface{}) bool {
	if res, ok := obj.(map[string]interface{}); ok {
		for k, temp := range where {
			if k == "$or" {
				array, ok := temp.([]interface{})
				if !ok {
					return false
				}
				for _, val := range array {
					value := val.(map[string]interface{})
					if validate(value, res) {
						return true
					}
				}
				return false
			}

			val, p := res[k]
			if !p {
				return false
			}
			// And clause
			cond, ok := temp.(map[string]interface{})
			if !ok {
				if temp != val {
					return false
				}
				continue
			}

			// match condition
			for k2, v2 := range cond {
				if reflect.TypeOf(val) != reflect.TypeOf(v2) {
					return false
				}
				switch k2 {
				case "$eq":
					if val != v2 {
						return false
					}
				case "$neq":
					if val == v2 {
						return false
					}
				case "$gt":
					switch val2 := val.(type) {
					case string:
						v2String := v2.(string)
						if val2 <= v2String {
							return false
						}
					case int64:
						v2Int := v2.(int64)
						if val2 <= v2Int {
							return false
						}
					case float64:
						v2Float := v2.(float64)
						if val2 <= v2Float {
							return false
						}
					default:
						return false
					}
				case "$gte":
					switch val2 := val.(type) {
					case string:
						v2String := v2.(string)
						if val2 < v2String {
							return false
						}
					case int64:
						v2Int := v2.(int64)
						if val2 < v2Int {
							return false
						}
					case float64:
						v2Float := v2.(float64)
						if val2 < v2Float {
							return false
						}
					default:
						return false
					}

				case "$lt":
					switch val2 := val.(type) {
					case string:
						v2String := v2.(string)
						if val2 >= v2String {
							return false
						}
					case int64:
						v2Int := v2.(int64)
						if val2 >= v2Int {
							return false
						}
					case float64:
						v2Float := v2.(float64)
						if val2 >= v2Float {
							return false
						}
					default:
						return false
					}

				case "$lte":
					switch val2 := val.(type) {
					case string:
						v2String := v2.(string)
						if val2 > v2String {
							return false
						}
					case int64:
						v2Int := v2.(int64)
						if val2 > v2Int {
							return false
						}
					case float64:
						v2Float := v2.(float64)
						if val2 > v2Float {
							return false
						}
					default:
						return false
					}
				}
			}
		}
		return true
	}
	if resArray, ok := obj.([]interface{}); ok {
		for _, res := range resArray {
			tempObj, ok := res.(map[string]interface{})
			if !ok {
				return false
			}
			if !validate(where, tempObj) {
				return false
			}
		}
		return true
	}
	return false
}
