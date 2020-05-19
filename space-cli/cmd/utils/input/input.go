package input

import "github.com/AlecAivazis/survey/v2"

type inputInterface interface {
	AskOne(p survey.Prompt, respone interface{}, opts ...survey.AskOpt) error
}

// Input struct for parameter
type input struct{}

// AskOne calls survey.AskOne
func (i *input) AskOne(p survey.Prompt, respone interface{}, opts ...survey.AskOpt) error {
	if err := survey.AskOne(p, respone, opts...); err != nil {
		return err
	}
	return nil
}

// Survey package for survey.Askone
var Survey inputInterface

func init() {
	Survey = &input{}
}
