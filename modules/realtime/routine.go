package realtime

import (
	"log"
	"reflect"
	"sync"

	"github.com/spaceuptech/space-cloud/model"
	"github.com/spaceuptech/space-cloud/utils"
)

func (m *Module) helperSendFeed(data *model.FeedData) {
	clientsTemp, ok := m.groups.Load(createGroupKey(data.DBType, data.Group))
	if !ok {
		log.Println("Realtime hanlder could not find key:", createGroupKey(data.DBType, data.Group))
		return
	}

	clients := clientsTemp.(*clientsStub)
	clients.clients.Range(func(key interface{}, value interface{}) bool {
		queries := value.(*sync.Map)
		queries.Range(func(id interface{}, value interface{}) bool {
			query := value.(*queryStub)

			dataPoint := &model.FeedData{
				QueryID: id.(string), DocID: data.DocID, Group: data.Group, Payload: data.Payload,
				TimeStamp: data.TimeStamp, Type: data.Type, DBType: data.DBType,
			}

			switch data.Type {
			case utils.RealtimeDelete:
				query.sendFeed(dataPoint)

			case utils.RealtimeInsert, utils.RealtimeUpdate:
				if validate(query.whereObj, data.Payload) {
					query.sendFeed(dataPoint)
				}

			default:
				log.Println("Realtime Module Error: Invalid event type received -", data.Type)
			}
			return true
		})
		return true
	})
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
