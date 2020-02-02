package auth

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/config"
	"github.com/spaceuptech/space-cloud/modules/crud"
	"github.com/spaceuptech/space-cloud/modules/functions"
)

func TestMatch(t *testing.T) {
	var authMatch = []struct {
		testName string
		err      error
		rule     *config.Rule
		args     map[string]interface{}
	}{
		{testName: "Match String !=", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match String ==", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "==", Type: "string", F1: "interfaceString1", F2: "interfaceString1", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match Number ==", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "==", Type: "number", F1: 1.0, F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match Number <=", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "<=", Type: "number", F1: 2.0, F2: 3.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match Number >=", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: ">=", Type: "number", F1: 3.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match Number <", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "<", Type: "number", F1: 1.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match Number >", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: ">", Type: "number", F1: 3.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match Number !=", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "number", F1: 1.0, F2: 10.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match bool ==", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "==", Type: "bool", F1: true, F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match bool !=", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "Rule", Eval: "!=", Type: "bool", F1: false, F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		// error
		{testName: "Error Test", err: ErrIncorrectMatch, rule: &config.Rule{}},
		{testName: "Match String !=", err: ErrIncorrectMatch, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "match", Eval: "", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match Number ==", err: ErrIncorrectRuleFieldType, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "match", Eval: "", Type: "number", F1: 1.0, F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match bool ==", err: ErrIncorrectRuleFieldType, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "match", Eval: "", Type: "bool", F1: true, F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match String !=", err: ErrIncorrectRuleFieldType, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "match", Eval: "!=", Type: "string", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match String !=", err: ErrIncorrectRuleFieldType, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "match", Eval: "!=", Type: "string", F1: "interfaceString1", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match number !=", err: errors.New("Store: Cloud not load value"), args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "match", Eval: "!=", Type: "number", F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match number !=", err: errors.New("Store: Cloud not load value"), args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "match", Eval: "!=", Type: "number", F1: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match bool !=", err: errors.New("Store: Cloud not load value"), args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "match", Eval: "!=", Type: "bool", F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match bool !=", err: errors.New("Store: Cloud not load value"), args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "match", Eval: "!=", Type: "bool", F1: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
	}
	testcases := 9
	for i, test := range authMatch {
		t.Run(test.testName, func(t *testing.T) {
			err := match(test.rule, test.args)
			if i <= testcases {
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

func TestMatchAnd(t *testing.T) {
	var authMatchAnd = []struct {
		testName string
		err      error
		rule     *config.Rule
		args     map[string]interface{}
	}{
		{testName: "Success Match Test", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test and !=", rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "and", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test and ==", rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "and", Eval: "==", Type: "number", F1: 1.0, F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test and <=", rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "and", Eval: "<=", Type: "number", F1: 2.0, F2: 3.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test and ==", rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "and", Eval: "==", Type: "string", F1: "interfaceString1", F2: "interfaceString1", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test and >=", rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: ">=", Type: "number", F1: 3.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test and <", rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "<", Type: "number", F1: 1.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test and >", rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: ">", Type: "number", F1: 3.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test and !=", rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "number", F1: 3.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test and ==", rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "==", Type: "bool", F1: true, F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test and !=", rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "bool", F1: true, F2: false, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		// error
		{testName: "Error Test and != type is not provided", err: ErrIncorrectMatch, rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "", F1: true, F2: false, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Error Match String eval is not provided !=", err: ErrIncorrectMatch, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Error Match String f1 is not provided !=", err: ErrIncorrectRuleFieldType, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "string", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Error Match String f2 is not provided !=", err: ErrIncorrectRuleFieldType, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "string", F1: "interfaceString1", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},

		{testName: "Error Match number eval is not provided !=", err: ErrIncorrectRuleFieldType, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "", Type: "number", F1: 1.0, F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Error Match number f1 is not provided !=", err: errors.New("Store: Cloud not load value"), args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "number", F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Error Match number f2 is not provided !=", err: errors.New("Store: Cloud not load value"), args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "number", F1: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},

		{testName: "Error Match bool eval is not provided !=", err: ErrIncorrectRuleFieldType, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "", Type: "bool", F1: true, F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Error Match bool f1 is not provided !=", err: errors.New("Store: Cloud not load value"), args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "bool", F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Error Match bool f2 is not provided !=", err: errors.New("Store: Cloud not load value"), args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "bool", F1: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
	}
	testcases := 10
	for i, test := range authMatchAnd {
		t.Run(test.testName, func(t *testing.T) {
			err := matchAnd(test.rule, test.args)
			if i <= testcases {
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

func TestMatchOr(t *testing.T) {
	var authMatchOr = []struct {
		testName string
		err      error
		rule     *config.Rule
		args     map[string]interface{}
	}{
		{testName: "Success Match Test", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test or !=", rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "and", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test or ==", rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "and", Eval: "==", Type: "number", F1: 1.0, F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test or <=", rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "and", Eval: "<=", Type: "number", F1: 2.0, F2: 3.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test or ==", rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "and", Eval: "==", Type: "string", F1: "interfaceString1", F2: "interfaceString1", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test or >=", rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: ">=", Type: "number", F1: 3.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test or <", rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: "<", Type: "number", F1: 1.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test or >", rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: ">", Type: "number", F1: 3.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test or !=", rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "number", F1: 3.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test or ==", rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: "==", Type: "bool", F1: true, F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test or !=", rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "bool", F1: true, F2: false, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		// error
		{testName: "Error Test and != type is not provided", err: ErrIncorrectMatch, rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "", F1: true, F2: false, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Error Match String eval is not provided !=", err: ErrIncorrectMatch, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: "", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Error Match String f1 is not provided !=", err: ErrIncorrectMatch, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "string", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Error Match String f2 is not provided !=", err: ErrIncorrectMatch, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "string", F1: "interfaceString1", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},

		{testName: "Error Match number eval is not provided !=", err: ErrIncorrectMatch, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: "", Type: "number", F1: 1.0, F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Error Match number f1 is not provided !=", err: ErrIncorrectMatch, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "number", F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Error Match number f2 is not provided !=", err: ErrIncorrectMatch, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "number", F1: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},

		{testName: "Error Match bool eval is not provided !=", err: ErrIncorrectMatch, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: "", Type: "bool", F1: true, F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Error Match bool f1 is not provided !=", err: ErrIncorrectMatch, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "bool", F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Error Match bool f2 is not provided !=", err: ErrIncorrectMatch, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "bool", F1: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
	}
	testcases := 0
	for i, test := range authMatchOr {
		t.Run(test.testName, func(t *testing.T) {
			err := matchOr(test.rule, test.args)
			if i <= testcases {
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

func TestMatchRule(t *testing.T) {
	var authMatchRule = []struct {
		testName string
		err      error
		rule     *config.Rule
		args     map[string]interface{}
	}{
		{testName: "Success Match allow", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "allow", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Success Match authenticated", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "authenticated", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Success Match match", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "match", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},

		{testName: "Success Match And", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test and !=", rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "and", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test and ==", rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "and", Eval: "==", Type: "number", F1: 1.0, F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test and <=", rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "and", Eval: "<=", Type: "number", F1: 2.0, F2: 3.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test and ==", rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "and", Eval: "==", Type: "string", F1: "interfaceString1", F2: "interfaceString1", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test and >=", rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: ">=", Type: "number", F1: 3.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test and <", rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "<", Type: "number", F1: 1.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test and >", rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: ">", Type: "number", F1: 3.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test and !=", rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "number", F1: 3.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test and ==", rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "==", Type: "bool", F1: true, F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test and !=", rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "bool", F1: true, F2: false, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},

		{testName: "Success Match Or", args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test or !=", rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "and", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test or ==", rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "and", Eval: "==", Type: "number", F1: 1.0, F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test or <=", rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "and", Eval: "<=", Type: "number", F1: 2.0, F2: 3.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test or ==", rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "and", Eval: "==", Type: "string", F1: "interfaceString1", F2: "interfaceString1", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test or >=", rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: ">=", Type: "number", F1: 3.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test or <", rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: "<", Type: "number", F1: 1.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test or >", rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: ">", Type: "number", F1: 3.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test or !=", rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "number", F1: 3.0, F2: 2.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test or ==", rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: "==", Type: "bool", F1: true, F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Successful Test or !=", rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "Rule", Eval: "!=", Type: "bool", F1: true, F2: false, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		// errors
		{testName: "Error Test and != deny", err: ErrIncorrectMatch, rule: &config.Rule{Rule: "deny", Clauses: []*config.Rule{{Rule: "and", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Error Test and != default", err: ErrIncorrectMatch, rule: &config.Rule{Rule: "", Clauses: []*config.Rule{{Rule: "and", Eval: "!=", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Error Test and != incorrect type", err: ErrIncorrectMatch, rule: &config.Rule{Rule: "and", Clauses: []*config.Rule{{Rule: "and", Eval: "!=", Type: "", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},
		{testName: "Error Test or != incorrect type", err: ErrIncorrectMatch, rule: &config.Rule{Rule: "or", Clauses: []*config.Rule{{Rule: "and", Eval: "!=", Type: "", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}}}},

		{testName: "Match String !=", err: ErrIncorrectMatch, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "match", Eval: "", Type: "string", F1: "interfaceString1", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match Number ==", err: ErrIncorrectRuleFieldType, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "match", Eval: "", Type: "number", F1: 1.0, F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Match bool ==", err: ErrIncorrectRuleFieldType, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "match", Eval: "", Type: "bool", F1: true, F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match String !=", err: ErrIncorrectRuleFieldType, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "match", Eval: "!=", Type: "string", F2: "interfaceString2", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match String !=", err: ErrIncorrectRuleFieldType, args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "match", Eval: "!=", Type: "string", F1: "interfaceString1", DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match number !=", err: errors.New("Store: Cloud not load value"), args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "match", Eval: "!=", Type: "number", F2: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match number !=", err: errors.New("Store: Cloud not load value"), args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "match", Eval: "!=", Type: "number", F1: 1.0, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match bool !=", err: errors.New("Store: Cloud not load value"), args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "match", Eval: "!=", Type: "bool", F2: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
		{testName: "Error Match bool !=", err: errors.New("Store: Cloud not load value"), args: map[string]interface{}{"string1": "interface1", "string2": "interface2"}, rule: &config.Rule{Rule: "match", Eval: "!=", Type: "bool", F1: true, DB: "DB", Col: "Col", Find: map[string]interface{}{"findstring1": "inteface1", "findstring2": "interface2"}}},
	}
	testcases := 24
	authModule := Init(&crud.Module{}, &functions.Module{})
	authModule.project = "project"
	for i, test := range authMatchRule {
		t.Run(test.testName, func(t *testing.T) {
			err := authModule.matchRule("project", test.rule, test.args, map[string]interface{}{})
			if i <= testcases {
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

// todo implement
func TestMatchQuery(t *testing.T) {
	var authMatchQuery = []struct {
		testName string
		err      error
		rule     *config.Rule
		crud     *crud.Module
		args     map[string]interface{}
	}{}
	testcases := 4
	for i, test := range authMatchQuery {
		t.Run(test.testName, func(t *testing.T) {
			err := matchQuery("project", test.rule, test.crud, test.args)
			if i <= testcases {
				if err != test.err {
					t.Error("Success Got Err", err, "Want Error", test.err)
				}
			} else {
				if err == test.err {
					t.Error("Error : Got Err", err, "Want Error", test.err)
				}
			}
		})
	}
}

// func TestModule_matchAnd(t *testing.T) {
// 	type args struct {
// 		ctx       context.Context
// 		projectID string
// 		rule      *config.Rule
// 		args      map[string]interface{}
// 		auth      map[string]interface{}
// 	}
// 	tests := []struct {
// 		name    string
// 		m       *Module
// 		args    args
// 		want    *PostProcess
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "Success Match Test",
// 			args: args{
// 				rule: &config.Rule{
// 					Rule: "Rule",
// 					Eval: "!=",
// 					Type: "string",
// 					F1:   "interfaceString1",
// 					F2:   "interfaceString2",
// 					DB:   "DB",
// 					Col:  "Col",
// 					Find: map[string]interface{}{
// 						"findstring1": "inteface1",
// 						"findstring2": "interface2",
// 					},
// 				},
// 			},
// 			m: &Module{
// 				project: "test",
// 			},
// 			wantErr: false,
// 			want:    &PostProcess{},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := tt.m.matchAnd(tt.args.ctx, tt.args.projectID, tt.args.rule, tt.args.args, tt.args.auth)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Module.matchAnd() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Module.matchAnd() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestModule_matchForce(t *testing.T) {
	type args struct {
		ctx       context.Context
		projectID string
		rule      *config.Rule
		args      map[string]interface{}
		auth      map[string]interface{}
	}
	tests := []struct {
		name    string
		m       *Module
		args    args
		want    *PostProcess
		wantErr bool
	}{
		// Test cases
		{
			name:    "Successful Test Force Rule",
			wantErr: false,
			args: args{
				projectID: "test1",
				rule: &config.Rule{
					Rule:  "force",
					Field: "res.field",
					Value: "someValue",
				},
			},
			m: &Module{
				project: "project",
			},
			want: &PostProcess{
				postProcessAction: []PostProcessAction{
					{
						Action: "force",
						Field:  "res.field",
						Value:  "someValue",
					},
				},
			},
		},
		{
			name:    "Failed Test Force Rule -> prefix != res/args",
			wantErr: true,
			args: args{
				projectID: "test1",
				rule: &config.Rule{
					Rule:  "force",
					Field: "field",
					Value: "someValue",
				},
			},
			m: &Module{
				project: "project",
			},
			want: nil,
		},
		{
			name:    "Failed Test Force Rule -> prefix == args",
			wantErr: true,
			args: args{
				projectID: "test1",
				rule: &config.Rule{
					Rule:  "force",
					Field: "args.field",
					Value: "someValue",
				},
			},
			m: &Module{
				project: "project",
			},
			want: &PostProcess{},
		},
		{
			name:    "Failed Test Force Rule -> empty clause",
			wantErr: true,
			args: args{
				projectID: "test1",
				rule: &config.Rule{
					Rule:   "force",
					Field:  "args.field",
					Value:  "someValue",
					Clause: &config.Rule{},
				},
			},
			m: &Module{
				project: "project",
			},
			want: &PostProcess{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.matchForce(tt.args.ctx, tt.args.projectID, tt.args.rule, tt.args.args, tt.args.auth)
			if (err != nil) != tt.wantErr {
				t.Errorf("Module.matchForce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Module.matchForce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_matchRemove(t *testing.T) {
	type args struct {
		ctx       context.Context
		projectID string
		rule      *config.Rule
		args      map[string]interface{}
		auth      map[string]interface{}
	}
	tests := []struct {
		name    string
		m       *Module
		args    args
		want    *PostProcess
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "Successful Test Remove Rule - res",
			wantErr: false,
			args: args{
				projectID: "test1",
				rule: &config.Rule{
					Rule:   "remove",
					Fields: []string{"res.field"},
					Value:  "someValue",
				},
			},
			m: &Module{},
			want: &PostProcess{
				postProcessAction: []PostProcessAction{
					{
						Action: "remove",
						Field:  "res.field",
						Value:  nil,
					},
				},
			},
		},
		{
			name:    "Failed Test Remove Rule - args",
			wantErr: true,
			args: args{
				projectID: "test1",
				rule: &config.Rule{
					Rule:   "remove",
					Fields: []string{"args.field"},
					Value:  "someValue",
				},
			},
			m:    &Module{},
			want: nil,
		},
		{
			name:    "Failed Test Remove Rule - invalid prefix",
			wantErr: true,
			args: args{
				projectID: "test1",
				rule: &config.Rule{
					Rule:   "remove",
					Fields: []string{"field"},
					Value:  "someValue",
				},
			},
			m:    &Module{},
			want: nil,
		},
		{
			name:    "Successful Test Remove Rule -> empty clause/rule",
			wantErr: false,
			args: args{
				projectID: "test1",
				rule: &config.Rule{
					Rule:  "remove",
					Field: "res.field",
					Value: "someValue",
					Clause: &config.Rule{
						Rule: "",
					},
				},
			},
			m: &Module{
				project: "project",
			},
			want: &PostProcess{},
		},
		{
			name:    "Failed Test Remove Rule -> with clause/rule",
			wantErr: true,
			args: args{
				rule: &config.Rule{
					Rule:  "remove",
					Field: "res.field",
					Value: "someValue",
					Clause: &config.Rule{
						Rule: "match",
					},
				},
			},
			m: &Module{
				project: "project",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.matchRemove(tt.args.ctx, tt.args.projectID, tt.args.rule, tt.args.args, tt.args.auth)
			if (err != nil) != tt.wantErr {
				t.Errorf("Module.matchRemove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Module.matchRemove() = %v, want %v", got, tt.want)
			}
		})
	}
}
