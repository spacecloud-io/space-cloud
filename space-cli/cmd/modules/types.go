package modules

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/stretchr/testify/mock"
)

type mockInputInterface struct {
	mock.Mock
}

func (m *mockInputInterface) AskOne(p survey.Prompt, respone interface{}, opts ...survey.AskOpt) error {
	c := m.Called(p, respone, opts)
	return c.Error(0)
}
