package helpers

import (
	"context"
	"net/http"

	"github.com/segmentio/ksuid"
)

type contextKey string

const internalRequestID = "internal"
const noValueRequestID = "noValue"

const contextKeyRequestID = contextKey("requestId")

func CreateContext(r *http.Request) context.Context {
	return context.WithValue(r.Context(), contextKeyRequestID, r.Header.Get(HeaderRequestID))
}

func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return internalRequestID + "-" + ksuid.New().String()
	}
	value := ctx.Value(contextKeyRequestID)
	if value == nil {
		return noValueRequestID + "-" + ksuid.New().String()
	}
	return value.(string)
}
