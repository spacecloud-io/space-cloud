package utils

import (
	"net/http"
	"strconv"

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

//ExtractRequestParams extract request info from http request & stores it in reqParam variable
func ExtractRequestParams(r *http.Request, reqParams *model.RequestParams, body interface{}) {
	if reqParams == nil {
		reqParams = &model.RequestParams{}
	}
	reqParams.RequestID = r.Header.Get(helpers.HeaderRequestID)
	reqParams.Method = r.Method
	reqParams.Path = r.URL.Path
	reqParams.Headers = r.Header
	reqParams.Payload = body
}
