package utils

import (
	"encoding/json"

	"github.com/AlecAivazis/survey/v2"
	"github.com/stretchr/testify/mock"
)

// MockInputInterface is the mock interface for survey package
type MockInputInterface struct {
	mock.Mock
}

// AskOne is the mock method for MockInputInterface
func (m *MockInputInterface) AskOne(p survey.Prompt, response interface{}) error {
	c := m.Called(p, response)
	a, _ := json.Marshal(c.Get(1))
	_ = json.Unmarshal(a, response)
	return c.Error(0)
}
