package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

// AcceptableIDType converts a provied id to string
func AcceptableIDType(id interface{}) (string, bool) {
	switch v := id.(type) {
	case string:
		return v, true
	case int:
		return strconv.Itoa(v), true
	case int32:
		return strconv.FormatInt(int64(v), 10), true
	case int64:
		return strconv.FormatInt(v, 10), true
	case float64:
		// json.Unmarshal converts all numbers to float64
		if float64(int64(v)) == v { // v is actually an int
			return strconv.FormatInt(int64(v), 10), true
		}
		return "", false
	default:
		return "", false
	}
}

// GetIDVariable gets the id variable for the provided db type
func GetIDVariable(dbAlias string) string {
	idVar := "id"
	if model.DBType(dbAlias) == model.Mongo {
		idVar = "_id"
	}

	return idVar
}

// ArrayContains checks if the array contains the value provided
func ArrayContains(array []interface{}, value interface{}) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

// ExtractRequestParams extract request info from http request & stores it in reqParam variable
func ExtractRequestParams(r *http.Request, reqParams model.RequestParams, body interface{}) model.RequestParams {
	reqParams.RequestID = r.Header.Get(helpers.HeaderRequestID)
	reqParams.Method = r.Method
	reqParams.Path = r.URL.Path
	reqParams.Headers = r.Header
	reqParams.Payload = body
	return reqParams
}

// ExtractJoinInfoForInstantInvalidate extracts join info
func ExtractJoinInfoForInstantInvalidate(join []*model.JoinOption, joinKeysMapping map[string]map[string]string) {
	for _, j := range join {
		GenerateJoinKeysForInstantInvalidate(j.Table, j.On, joinKeysMapping)
		if j.Join != nil {
			ExtractJoinInfoForInstantInvalidate(j.Join, joinKeysMapping)
		}
	}
}

// GenerateJoinKeys generates join keys
func GenerateJoinKeys(joinTable string, joinOn map[string]interface{}, databaseRow map[string]interface{}, joinKeysMapping map[string]map[string]string) {
	isValidJoin, columnName := IsValidJoin(joinOn, joinTable)
	dbRow, ok := databaseRow[joinTable+"__"+columnName]
	if isValidJoin && ok {
		outerKey := fmt.Sprintf("%s::%s::%s", joinTable, "join", columnName)
		rowValue := fmt.Sprintf("%v", dbRow)
		_, ok := joinKeysMapping[outerKey]
		if !ok {
			joinKeysMapping[outerKey] = map[string]string{rowValue: ""}
		} else {
			joinKeysMapping[outerKey][rowValue] = ""
		}
	} else {
		outerKey := fmt.Sprintf("%s::%s::%s", joinTable, "always", "none")
		joinKeysMapping[outerKey] = map[string]string{"none": ""}
	}
}

// GenerateJoinKeysForInstantInvalidate generates join keys
func GenerateJoinKeysForInstantInvalidate(joinTable string, joinOn map[string]interface{}, joinKeysMapping map[string]map[string]string) {
	isValidJoin, columnName := IsValidJoin(joinOn, joinTable)
	if isValidJoin {
		outerKey := fmt.Sprintf("%s::%s::%s", joinTable, "join", columnName)
		joinKeysMapping[outerKey] = map[string]string{}
		outerKey = fmt.Sprintf("%s::%s::%s", joinTable, "always", "none")
		joinKeysMapping[outerKey] = map[string]string{"none": ""}
	} else {
		outerKey := fmt.Sprintf("%s::%s::%s", joinTable, "always", "none")
		joinKeysMapping[outerKey] = map[string]string{"none": ""}
	}
}

// IsValidJoin checks if join is valid
func IsValidJoin(on map[string]interface{}, jointTableName string) (bool, string) {
	if len(on) > 1 {
		return false, "none"
	}

	// Its not a valid join when the or condition is used
	_, ok := on["$or"]
	if ok {
		return false, "none"
	}

	columnName := "none"
	for leftKey, value := range on {
		// See if left side is our value
		if arr := strings.Split(leftKey, "."); len(arr) > 0 && arr[0] == jointTableName {
			columnName = arr[1]
		}

		isValid, col := checkJoinType(value, jointTableName)
		if !isValid {
			return false, "none"
		}
		if col != "none" {
			columnName = col
		}
	}
	return columnName != "none", columnName
}

func checkJoinType(value interface{}, jointTableName string) (bool, string) {
	switch t := value.(type) {
	case string:
		columnName := "none"
		if arr := strings.Split(t, "."); len(arr) > 0 && arr[0] == jointTableName {
			columnName = arr[1]
		}
		return true, columnName
	case map[string]interface{}:
		value2, ok := t["$eq"]
		if !ok {
			return false, "none"
		}
		return checkJoinType(value2, jointTableName)
	default:
		return false, "none"
	}
}
