package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/spaceuptech/helpers"
)

// CheckParse checks if the string can be parsed or not
func CheckParse(s string) (time.Time, error) {
	var value time.Time
	var err error
	value, err = time.Parse(time.RFC3339Nano, s)
	if err != nil {
		value, err = time.Parse("2006-01-02", s)
		if err != nil {
			return time.Time{}, helpers.Logger.LogError(helpers.GetRequestID(context.TODO()), fmt.Sprintf("Invalid date format (%s) provided", s), nil, nil)
		}
	}
	return value, nil
}
