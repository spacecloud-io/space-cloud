package utils

import (
	"fmt"
	"time"
)

// CheckParse checks if the string can be parsed or not
func CheckParse(s string) (time.Time, error) {
	var value time.Time
	var err error
	value, err = time.Parse(time.RFC3339, s)
	if err != nil {
		value, err = time.Parse("2006-01-02", s)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid date format (%s) provided", s)
		}
	}
	return value, nil
}
