package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

// ErrRuleNotFound is thrown when an error is not present in the auth object
var ErrRuleNotFound = errors.New("auth: No rule has been provided")

// ErrIncorrectRuleFieldType is thrown when the field type of a rule is of incorrect type
var ErrIncorrectRuleFieldType = errors.New("auth: Incorrect rule field type")

// ErrIncorrectMatch is thrown when the field type of a rule is of incorrect type
var ErrIncorrectMatch = errors.New("auth: The two fields do not match")

// FormatError check whether error is provided in config.Rule
func formatError(ctx context.Context, rule *config.Rule, err error) error {
	if err == nil {
		return nil
	}

	name := rule.Name
	if name == "" {
		name = "with no name"
	}

	_ = helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("Rule (%s) of type (%s) failed", name, rule.Rule), err, map[string]interface{}{})

	if rule.Error == "" {
		return err
	}
	return errors.New(rule.Error)
}
