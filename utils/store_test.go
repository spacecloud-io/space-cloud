package utils

import (
	"reflect"
	"testing"
)

type LoadStringIfExistStub struct {
	value string
	state map[string]interface{}
	ret   string
}
type LoadValueStub struct {
	value string
	state map[string]interface{}
	ret   interface{}
}

func TestUtilsExists(t *testing.T) {
	trueCases := 1
	m := map[string]interface{}{
		"args": map[string]interface{}{
			"auth": "id",
		},
	}
	test := []*LoadValueStub{
		//1
		&LoadValueStub{
			value: "utils.exists(args.auth)",
			state: m,
			ret:   true,
		},
		//False
		&LoadValueStub{
			value: "utils.exists(args.auth.id)",
			state: m,
			ret:   true,
		},
		&LoadValueStub{
			value: "utils.abc",
			state: m,
			ret:   false,
		},
	}

	for i, eachTest := range test {
		res, err := LoadValue(eachTest.value, eachTest.state)
		eq := reflect.DeepEqual(eachTest.ret, res)
		if i < trueCases {
			if ((res == false || !eq) && err == nil) || err != nil {
				t.Error(i+1, ":", "Incorrect Match 1")
			}
			continue
		} else if (res == true || eq) && err == nil {
			t.Error(i+1, ":", "Incorrect Match 2")
		}
	}
}

func TestLoadValue(t *testing.T) {
	trueCases := 3
	m := map[string]interface{}{
		"args": map[string]interface{}{
			"auth":   "key",
			"nested": "key2",
			"group1": map[string]interface{}{
				"key": "nested",
				"id":  "value",
				"nested": map[string]interface{}{
					"id": "group1",
				},
			},
			"group2": map[string]interface{}{},
		},
	}
	empty := map[string]interface{}{}
	onelevel := map[string]interface{}{
		"args": "id",
	}
	test := []*LoadValueStub{
		&LoadValueStub{
			value: "args.auth",
			state: m,
			ret:   "key",
		},
		&LoadValueStub{
			value: "args.group1[args.auth]",
			state: m,
			ret:   "nested",
		},
		&LoadValueStub{
			value: "args.group1[args.group1.key].id",
			state: m,
			ret:   "group1",
		},
		//False / Error Cases.
		&LoadValueStub{
			value: "",
			state: m,
			ret:   "",
		},
		&LoadValueStub{
			value: "args",
			state: m,
			ret:   "args",
		},
		&LoadValueStub{
			value: "args.auth",
			state: empty,
			ret:   "id",
		},
		&LoadValueStub{
			value: "args.auth",
			state: onelevel,
			ret:   "id",
		},
		&LoadValueStub{
			value: "args.group2[args.auth]",
			state: m,
			ret:   "id",
		},
		&LoadValueStub{
			value: "args.group3.abc",
			state: m,
			ret:   "",
		},
	}
	for i, eachTest := range test {
		res, err := LoadValue(eachTest.value, eachTest.state)
		eq := reflect.DeepEqual(eachTest.ret, res)
		if i < trueCases {
			if err != nil || (err == nil && (!eq)) {
				t.Error(i+1, ":", "Incorrect Match", err)
			}
		} else if err == nil {
			t.Error(i+1, ":", "Incorrect Match")
		}
	}
}

// func TestLoadStringIfExists(t *testing.T) {

// 	m := map[string]interface{}{
// 		"args": map[string]interface{}{
// 			"auth": map[string]interface{}{
// 				"id":   "UserID",
// 				"role": "admin",
// 				"key1": "value1",
// 				"key2": "value2",
// 				"nested": map[string]interface{}{
// 					"value1": "v1",
// 				},
// 			},
// 		},
// 	}
// 	k := map[string]interface{}{
// 		"files": map[string]interface{}{
// 			"file1": "myfile",
// 		},
// 		"args": "",
// 	}
// 	l := map[string]interface{}{
// 		"": "Empty",
// 	}

// 	test := []*LoadStringIfExistStub{
// 		//1
// 		&LoadStringIfExistStub{
// 			value: "args.auth.id",
// 			state: m,
// 			ret:   "UserID",
// 		},
// 		//2
// 		&LoadStringIfExistStub{
// 			value: "args.auth.role",
// 			state: m,
// 			ret:   "admin",
// 		},
// 		//3
// 		&LoadStringIfExistStub{
// 			value: "",
// 			state: m,
// 			ret:   "",
// 		},
// 		//4
// 		&LoadStringIfExistStub{
// 			value: "args.auth.nested[args.auth.key1]",
// 			state: m,
// 			ret:   "v1",
// 		},
// 		//5
// 		&LoadStringIfExistStub{
// 			value: "args.auth.nested[args.auth.key2]",
// 			state: m,
// 			ret:   "args.auth.nested[args.auth.key2]",
// 		},
// 		//6
// 		&LoadStringIfExistStub{
// 			value: "(",
// 			state: m,
// 			ret:   "(",
// 		},
// 		//7
// 		&LoadStringIfExistStub{
// 			value: "args.auth",
// 			state: k,
// 			ret:   "args.auth",
// 		},
// 		//8
// 		&LoadStringIfExistStub{
// 			value: "files.myfile",
// 			state: k,
// 			ret:   "files.myfile",
// 		},
// 		//9
// 		&LoadStringIfExistStub{
// 			value: "files.file1",
// 			state: k,
// 			ret:   "myfile",
// 		},
// 		//10
// 		&LoadStringIfExistStub{
// 			value: "args.auth",
// 			state: l,
// 			ret:   "args.auth",
// 		},
// 		//False Cases
// 		&LoadStringIfExistStub{
// 			value: "args.auth.role",
// 			state: m,
// 			ret:   "args.auth.role",
// 		},
// 		&LoadStringIfExistStub{
// 			value: "args.auth.id",
// 			state: m,
// 			ret:   "args.auth.id",
// 		},
// 		&LoadStringIfExistStub{
// 			value: "",
// 			state: m,
// 			ret:   "Error",
// 		},
// 		&LoadStringIfExistStub{
// 			value: "args.auth.nested[args.auth.key1]",
// 			state: m,
// 			ret:   "args.auth.nested[args.auth.key1]",
// 		},
// 		&LoadStringIfExistStub{
// 			value: "(",
// 			state: m,
// 			ret:   "",
// 		},
// 		&LoadStringIfExistStub{
// 			value: "",
// 			state: l,
// 			ret:   "Empty",
// 		},
// 	}
// 	for i, eachTest := range test {
// 		res := LoadStringIfExists(eachTest.value, eachTest.state)
// 		eq := reflect.DeepEqual(eachTest.ret, res)
// 		if i < 10 {
// 			if !eq {
// 				t.Error(i+1, ":", "Incorrect Match")
// 			}
// 			continue
// 		}
// 		if eq {
// 			t.Error(i+1, ":", "Incorrect Match")
// 		}
// 	}
// }
