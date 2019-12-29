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
		{name: "Match String !=", isErrExpected: false, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F1: "interfacestring", F2: "interfaceString"}},
		{name: "Match String notin-Fail", isErrExpected: true, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "notin", Type: "string", F1: "interface", F2: []interface{}{"suyash", "interface"}}},
		{name: "Match String notin-Success", isErrExpected: false, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "notin", Type: "string", F1: "interface", F2: []interface{}{"interface1", "suyash"}}},
		{name: "Match String != loaded from state", isErrExpected: false, args: map[string]interface{}{"v1": "val1", "v2": "val2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F1: "args.v1", F2: "args.v2"}},
		{name: "Match String in loaded from state-success", isErrExpected: false, args: map[string]interface{}{"v1": "var", "v2": []interface{}{"suyash", "var"}}, rule: &config.Rule{Rule: "Rule", Eval: "in", Type: "string", F1: "args.v1", F2: "args.v2"}},
		{name: "Match String in loaded from state of state-Success", isErrExpected: true, args: map[string]interface{}{"v1": "val1", "v2": map[string]interface{}{"v3": []interface{}{"suyash", "val"}}}, rule: &config.Rule{Rule: "Rule", Eval: "in", Type: "string", F1: "args.v1", F2: "args.v2.v3"}},
		{name: "Match String notin loaded from state-Fail", isErrExpected: true, args: map[string]interface{}{"v1": "var", "v2": []interface{}{"suyash", "var"}}, rule: &config.Rule{Rule: "Rule", Eval: "notin", Type: "string", F1: "args.v1", F2: "args.v2"}},
		{name: "Match String notin loaded from state-Success", isErrExpected: false, args: map[string]interface{}{"v1": "var", "v2": []interface{}{"suyash", "val"}}, rule: &config.Rule{Rule: "Rule", Eval: "notin", Type: "string", F1: "args.v1", F2: "args.v2"}},
		{name: "Match String != without passing args", isErrExpected: false, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F1: "interfacestring", F2: "interfaceString"}},
		{name: "Match String != variable contains invalid type", isErrExpected: true, args: map[string]interface{}{"v1": 10}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F1: "args.v1", F2: "interfaceString"}},
		{name: "Error Match String != passing only one value", isErrExpected: true, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F1: "interfacestring"}},
		{name: "Error Match String != passing float value", isErrExpected: true, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F1: "interfacestring", F2: 3.0}},
		{name: "Error Match String == value as an interface", isErrExpected: true, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "==", Type: "string", F1: "interfacestring", F2: map[string]interface{}{"val": "interfacestring"}}},
		{name: "Match String ==", isErrExpected: false, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "==", Type: "string", F1: "interfaceString1", F2: "interfaceString1", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{name: "Match String in-Success", isErrExpected: false, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "in", Type: "string", F1: "suyash", F2: []interface{}{"suyash", "interface"}}},
		{name: "Match String notin-Fail", isErrExpected: true, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "notin", Type: "string", F1: "suyash", F2: []interface{}{"interface", "interface", "suyash"}}},
		{name: "Match String in-Fail", isErrExpected: true, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "in", Type: "string", F1: "suyash", F2: []interface{}{"Suyash", "interface"}}},
		{name: "Error Match String eval is not provided !=", isErrExpected: true, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{name: "Error Match String !=", isErrExpected: true, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{name: "Error Match String !=", isErrExpected: true, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F1: "interfaceString1", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			gotErr := matchString(testCase.rule, map[string]interface{}{"args": testCase.args})
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
		{testName: "Match Number ==", isErrExpected: false, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "==", Type: "number", F1: 1.0, F2: 1}},
		{testName: "Match Number in-Success", isErrExpected: false, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "in", Type: "number", F1: 1.0, F2: []interface{}{2, 1}}},
		{testName: "Match Number in-Fail", isErrExpected: true, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "in", Type: "number", F1: 1.0, F2: []interface{}{2, 3}}},
		{testName: "Match Number in loaded from state-Success", isErrExpected: false, args: map[string]interface{}{"v1": []interface{}{2, 1}}, rule: &config.Rule{Rule: "Rule", Eval: "in", Type: "number", F1: 1.0, F2: "args.v1"}},
		{testName: "Match Number in loaded from state-Fail", isErrExpected: true, args: map[string]interface{}{"v1": []interface{}{2, 3}}, rule: &config.Rule{Rule: "Rule", Eval: "in", Type: "number", F1: 1.0, F2: "args.v1"}},
		{testName: "Match Number notin-Fail", isErrExpected: true, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "notin", Type: "number", F1: 1.0, F2: []interface{}{2, 1}}},
		{testName: "Match Number notin-Success", isErrExpected: false, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "notin", Type: "number", F1: 1.0, F2: []interface{}{2, 3}}},
		{testName: "Match Number notin loaded from state-Fail", isErrExpected: true, args: map[string]interface{}{"v1": []interface{}{2, 1}}, rule: &config.Rule{Rule: "Rule", Eval: "notin", Type: "number", F1: 1.0, F2: "args.v1"}},
		{testName: "Match Number notin loaded from state-Success", isErrExpected: false, args: map[string]interface{}{"v1": []interface{}{2, 3}}, rule: &config.Rule{Rule: "Rule", Eval: "notin", Type: "number", F1: 1.0, F2: "args.v1"}},
		{testName: "Match Number loaded from state ==", isErrExpected: false, args: map[string]interface{}{"num1": 12.0, "num2": 12}, rule: &config.Rule{Rule: "Rule", Eval: "==", Type: "number", F1: "args.num1", F2: "args.num2"}},
		{testName: "Match Number <=", isErrExpected: false, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "<=", Type: "number", F1: 2.54, F2: 3.67}},
		{testName: "Match Number >=", isErrExpected: false, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: ">=", Type: "number", F1: 3.67, F2: 2.54}},
		{testName: "Match Number <", isErrExpected: false, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "<", Type: "number", F1: 1, F2: 2.0}},
		{testName: "Match Number >", isErrExpected: false, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: ">", Type: "number", F1: 3.0, F2: 2.99}},
		{testName: "Match Number !=", isErrExpected: false, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "number", F1: 10.1, F2: 10}},
		{testName: "Match Number loaded from state !=", isErrExpected: false, args: map[string]interface{}{"num1": 12.34, "num2": 11}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "number", F1: "args.num1", F2: "args.num2"}},
		{testName: "Error Match Number eval is not provided", isErrExpected: true, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "", Type: "number", F1: 1.0, F2: 100}},
		{testName: "Error Match number !=(single field F2)", isErrExpected: true, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "number", F2: 1.0}},
		{testName: "Error Match number !=(single field F1)", isErrExpected: true, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "number", F1: 1.0}},
		{testName: "Error Match number != field does not exist", isErrExpected: true, args: map[string]interface{}{}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "number", F1: 1.0, F2: "args.num1"}},
		{testName: "Error Match number != field is of incorrect type", isErrExpected: true, args: map[string]interface{}{"num1": "wrong type"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "number", F1: 1.0, F2: "args.num1"}},
	}
	for _, test := range authMatchNumber {
		t.Run(test.testName, func(t *testing.T) {
			err := matchNumber(test.rule, map[string]interface{}{"args": test.args})
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
		{testName: "Match bool ==", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "==", Type: "bool", F1: true, F2: true}},
		{testName: "Match bool ==/invalid match error", isErrExpected: true, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "==", Type: "bool", F1: true, F2: false}},
		{testName: "Match bool !=", isErrExpected: false, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "bool", F1: false, F2: true}},
		{testName: "Match bool !=", isErrExpected: true, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "bool", F1: false, F2: false}},
		{testName: "Error Match bool eval is not provided", isErrExpected: true, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "", Type: "bool", F1: true, F2: true}},
		{testName: "Error Match bool !=(only F2 provided)", isErrExpected: true, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F2: true}},
		{testName: "Error Match bool !=(only F1 provided)", isErrExpected: true, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F1: true}},
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
