package middlewares

import (
	"context"
	"net/http"

	"github.com/spacecloud-io/space-cloud/model"
	"github.com/spf13/viper"
)

// GetRequestParams returns the request params stored in request context
func GetRequestParams(r *http.Request) *model.RequestParams {
	value := r.Context().Value("req-params")
	if value == nil {
		return nil
	}

	return value.(*model.RequestParams)
}

// StoreRequestParams stores the provided request params in the request context
func StoreRequestParams(r *http.Request, reqParams *model.RequestParams) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), "req-params", reqParams))
}

// IsRequestAuthenticated checks if the request is authenticated
func IsRequestAuthenticated(reqParams *model.RequestParams, isAdmin bool) bool {
	return reqParams.Claims != nil || (isAdmin && viper.GetBool("dev"))
}
