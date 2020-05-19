package utils

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/stretchr/testify/mock"
)

// MockInputInterface is the mock interface for survey package
type MockInputInterface struct {
	mock.Mock
}

// AskOne is the mock method for MockInputInterface
func (m *MockInputInterface) AskOne(p survey.Prompt, respone interface{}, opts ...survey.AskOpt) error {
	c := m.Called(p, respone, opts)
	return c.Error(0)
}
