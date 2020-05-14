package auth

import (
	"errors"
	"fmt"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// ErrRuleNotFound is thrown when an error is not present in the auth object
var ErrRuleNotFound = errors.New("auth: No rule has been provided")

// ErrIncorrectRuleFieldType is thrown when the field type of a rule is of incorrect type
var ErrIncorrectRuleFieldType = errors.New("auth: Incorrect rule field type")

// ErrIncorrectMatch is thrown when the field type of a rule is of incorrect type
var ErrIncorrectMatch = errors.New("auth: The two fields do not match")

// formatError check whether error is provided in config.Rule
func formatError(rule *config.Rule, err error) error {
	if err == nil {
		return nil
	}
	if rule.Name == "" {
		rule.Name = "no name"
	}
	_ = utils.LogError(fmt.Sprintf("Rule (%s) of type (%s) failed", rule.Name, rule.Rule), err)

	if rule.Error == "" {
		return err
	}
	return errors.New(rule.Error)
}
