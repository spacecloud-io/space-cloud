package auth

import (
	"errors"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/config"
)

func TestMatchString(t *testing.T) {
	var authMatchString = []struct {
		testName string
		err      error
		rule     *config.Rule
		args     map[string]interface{}
	}{
		{testName: "Match String !=", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match String ==", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "==", Type: "string", F1: "interfaceString1", F2: "interfaceString1", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match String eval is not provided !=", err: ErrIncorrectMatch, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match String !=", err: ErrIncorrectRuleFieldType, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match String !=", err: ErrIncorrectRuleFieldType, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F1: "interfaceString1", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
	}

	successTestCases := 1
	for i, test := range authMatchString {
		t.Run(test.testName, func(t *testing.T) {
			err := matchString(test.rule, test.args)
			if i <= successTestCases {
				if !reflect.DeepEqual(err, test.err) {
					t.Error("Success Test ", "| Got This |", err, "| Wanted This |", test.err)
				}
			} else {

				if !reflect.DeepEqual(err, test.err) {
					t.Error("Error Test", "| Got This |", err, "| Wanted This |", test.err)
				}
			}

		})
	}
}

func TestMatchNumber(t *testing.T) {
	var authMatchNumber = []struct {
		testName string
		err      error
		rule     *config.Rule
		args     map[string]interface{}
	}{
		{testName: "Match Number ==", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "==", Type: "number", F1: 1.0, F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match Number <=", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "<=", Type: "number", F1: 2.0, F2: 3.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match Number >=", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: ">=", Type: "number", F1: 3.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match Number <", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "<", Type: "number", F1: 1.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match Number >", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: ">", Type: "number", F1: 3.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match Number !=", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "number", F1: 1.0, F2: 10.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Erro Match Number eval is not provided ==", err: ErrIncorrectRuleFieldType, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "", Type: "number", F1: 1.0, F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match number !=", err: errors.New("Store: Cloud not load value"), args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match number !=", err: errors.New("Store: Cloud not load value"), args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F1: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
	}
	successTestCases := 5
	for i, test := range authMatchNumber {
		t.Run(test.testName, func(t *testing.T) {
			err := matchNumber(test.rule, test.args)
			if i <= successTestCases {
				if !reflect.DeepEqual(err, test.err) {
					t.Error("Success Test ", "| Got This |", err, "| Wanted This |", test.err)
				}
			} else {

				if !reflect.DeepEqual(err, test.err) {
					t.Error("Error Test ", "| Got This |", err, "| Wanted This |", test.err)
				}
			}

		})
	}
}

func TestMatchBool(t *testing.T) {
	var authMatchBool = []struct {
		testName string
		err      error
		rule     *config.Rule
		args     map[string]interface{}
	}{
		{testName: "Match bool ==", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "==", Type: "bool", F1: true, F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match bool !=", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "bool", F1: false, F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match bool == eval is not provided", err: ErrIncorrectRuleFieldType, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "", Type: "bool", F1: true, F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match bool !=", err: errors.New("Store: Cloud not load value"), args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match bool !=", err: errors.New("Store: Cloud not load value"), args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F1: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
	}

	successTestCases := 1
	for i, test := range authMatchBool {
		t.Run(test.testName, func(t *testing.T) {
			err := matchBool(test.rule, test.args)
			if i <= successTestCases {
				if !reflect.DeepEqual(err, test.err) {
					t.Error("Success Test ", "| Got This |", err, "| Wanted This |", test.err)
				}
			} else {

				if !reflect.DeepEqual(err, test.err) {
					t.Error("Error Test", "| Got This |", err, "| Wanted This |", test.err)
				}
			}

		})
	}
}
