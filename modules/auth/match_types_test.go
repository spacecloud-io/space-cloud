package auth

import (
	"github.com/spaceuptech/space-cloud/config"
	"testing"
)

func TestMatchString(t *testing.T) {
	var testCases = []struct {
		name          string
		err           error
		isErrExpected bool
		rule          *config.Rule
		args          map[string]interface{}
	}{
		{name: "Match String !=", isErrExpected: false, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F1: "interfacestring", F2: "interfaceString"}},
		{name: "Match String != without passing args", isErrExpected: false, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F1: "interfacestring", F2: "interfaceString"}},
		{name: "Error Match String != passing only one value", isErrExpected: true, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F1: "interfacestring"}},
		{name: "Error Match String != passing float value", isErrExpected: true, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F1: "interfacestring", F2: 3.0}},
		{name: "Error Match String == value as an interface", isErrExpected: true, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "==", Type: "string", F1: "interfacestring", F2: []interface{}{"interfacestring"}}},
		{name: "Match String ==", isErrExpected: false, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "==", Type: "string", F1: "interfaceString1", F2: "interfaceString1", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{name: "Error Match String eval is not provided !=", isErrExpected: true, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{name: "Error Match String !=", isErrExpected: true, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{name: "Error Match String !=", isErrExpected: true, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F1: "interfaceString1", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			gotErr := matchString(testCase.rule, testCase.args)
			if (gotErr != nil) != testCase.isErrExpected {
				t.Errorf("got %v wanted %v", gotErr, testCase.isErrExpected)
			}

		})
	}
}

func TestMatchNumber(t *testing.T) {
	var authMatchNumber = []struct {
		testName      string
		isErrExpected bool
		rule          *config.Rule
		args          map[string]interface{}
	}{
		{testName: "Match Number ==", isErrExpected: false, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "==", Type: "number", F1: 1.0, F2: 1, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match Number <=", isErrExpected: false, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "<=", Type: "number", F1: 2.54, F2: 3.67, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match Number >=", isErrExpected: false, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: ">=", Type: "number", F1: 3.67, F2: 2.54, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match Number <", isErrExpected: false, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "<", Type: "number", F1: 1, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match Number >", isErrExpected: false, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: ">", Type: "number", F1: 3.0, F2: 2.99, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match Number !=", isErrExpected: false, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "number", F1: 10.1, F2: 10, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match Number eval is not provided", isErrExpected: true, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "", Type: "number", F1: 1.0, F2: 100, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match number !=(single field F2)", isErrExpected: true, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match number !=(single field F1)", isErrExpected: true, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F1: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
	}
	for _, test := range authMatchNumber {
		t.Run(test.testName, func(t *testing.T) {
			err := matchNumber(test.rule, test.args)
			if (err != nil) != test.isErrExpected {
				t.Error("| Got This |", err, "| Wanted Error |", test.isErrExpected)
			}
		})
	}
}

func TestMatchBool(t *testing.T) {
	var authMatchBool = []struct {
		testName      string
		isErrExpected bool
		rule          *config.Rule
		args          map[string]interface{}
	}{
		{testName: "Match bool ==", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "==", Type: "bool", F1: true, F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match bool ==/invalid match error", isErrExpected: true, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "==", Type: "bool", F1: true, F2: false, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match bool !=", isErrExpected: false, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "bool", F1: false, F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match bool !=", isErrExpected: true, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "bool", F1: false, F2: false, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match bool eval is not provided", isErrExpected: true, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "", Type: "bool", F1: true, F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match bool !=(only F2 provided)", isErrExpected: true, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match bool !=(only F1 provided)", isErrExpected: true, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F1: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
	}

	for _, test := range authMatchBool {
		t.Run(test.testName, func(t *testing.T) {
			err := matchBool(test.rule, test.args)
			if (err != nil) != test.isErrExpected {
				t.Error("| Got This |", err, "| Wanted Error |", test.isErrExpected)
			}
		})
	}
}
