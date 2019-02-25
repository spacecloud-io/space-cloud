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
type LoadNumberStub struct {
	key  interface{}
	args map[string]interface{}
	ret  float64
}
type LoadBoolStub struct {
	key  interface{}
	args map[string]interface{}
	ret  bool
}
type AdjustStub struct {
	obj   interface{}
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
	trueCases := 4
	m := map[string]interface{}{
		"args": map[string]interface{}{
			"auth":   "key",
			"nested": "key2",
			"group1": map[string]interface{}{
				"key":  "nested",
				"key2": 10,
				"id":   "value",
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
		&LoadValueStub{
			value: "args.group1.key2",
			state: m,
			ret:   10,
		},
		// False/Error Cases.
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
		&LoadValueStub{
			value: "args.group1.key1",
			state: m,
			ret:   "",
		},
		&LoadValueStub{
			value: "args.group1[abc]",
			state: m,
			ret:   "",
		},
		&LoadValueStub{
			value: "args.group1.id2[args.auth]",
			state: m,
			ret:   "",
		},
		&LoadValueStub{
			value: "args.group3[args.auth].abc",
			state: m,
			ret:   "",
		},
		&LoadValueStub{
			value: "args.group1[args.auth].abc",
			state: m,
			ret:   "",
		},
		&LoadValueStub{
			value: "args.group1[args.auth1].abc",
			state: m,
			ret:   "",
		},
		&LoadValueStub{
			value: "abc",
			state: m,
			ret:   "",
		},
		&LoadValueStub{
			value: "args.group1[args.group1.key2].id",
			state: m,
			ret:   "",
		},
		&LoadValueStub{
			value: "args.group1[args.group1.key2]",
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

func TestLoadStringIfExists(t *testing.T) {
	trueCases := 2
	m := map[string]interface{}{
		"args": map[string]interface{}{
			"auth": map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}
	test := []*LoadStringIfExistStub{
		&LoadStringIfExistStub{
			value: "args.auth.key1",
			state: m,
			ret:   "value1",
		},
		&LoadStringIfExistStub{
			value: "args.auth.key3",
			state: m,
			ret:   "args.auth.key3",
		},
		//False Cases :
		&LoadStringIfExistStub{
			value: "args.auth.key",
			state: m,
			ret:   "value1",
		},
		&LoadStringIfExistStub{
			value: "args.auth.key1",
			state: m,
			ret:   "args.auth.key1",
		},
	}
	for i, eachTest := range test {
		res := LoadStringIfExists(eachTest.value, eachTest.state)
		eq := reflect.DeepEqual(eachTest.ret, res)
		if i < trueCases {
			if !eq {
				t.Error(i+1, ":", "Incorrect Match 1")
			}
			continue
		}
		if eq {
			t.Error(i+1, ":", "Incorrect Match 2")
		}
	}
}

func TestLoadNumber(t *testing.T) {
	trueCases := 3
	m := map[string]interface{}{
		"args": map[string]interface{}{
			"auth": map[string]interface{}{
				"key1": int64(10),
				"key2": float64(20),
				"key3": int(30),
			},
		},
	}
	test := []*LoadNumberStub{
		&LoadNumberStub{
			key:  "args.auth.key1",
			args: m,
			ret:  10,
		},
		&LoadNumberStub{
			key:  "args.auth.key2",
			args: m,
			ret:  20,
		},
		&LoadNumberStub{
			key:  int64(10),
			args: m,
			ret:  float64(10),
		},

		// False Cases :
		&LoadNumberStub{
			key:  "args.auth.key4",
			args: m,
			ret:  30,
		},
		&LoadNumberStub{
			key:  "args.auth.key3",
			args: m,
			ret:  30,
		},
		&LoadNumberStub{
			key:  int(10),
			args: m,
			ret:  float64(10),
		},
	}
	for i, eachTest := range test {
		res, err := LoadNumber(eachTest.key, eachTest.args)
		if i < trueCases {
			if err != nil || res != eachTest.ret {
				t.Error(i+1, ":", "Incorrect Match 1", err)
			}
			continue
		} else if err == nil && res == eachTest.ret {
			t.Error(i+1, ":", "Incorrect Match 2")
		}

	}
}

func TestLoadBool(t *testing.T) {
	trueCases := 2
	m := map[string]interface{}{
		"args": map[string]interface{}{
			"auth": map[string]interface{}{
				"key1": true,
			},
		},
	}
	test := []*LoadBoolStub{
		&LoadBoolStub{
			key:  "args.auth.key1",
			args: m,
			ret:  true,
		},
		&LoadBoolStub{
			key:  true,
			args: m,
			ret:  true,
		},
		// False Cases:
		&LoadBoolStub{
			key:  "args.auth.key2",
			args: m,
			ret:  true,
		},
		&LoadBoolStub{
			key:  false,
			args: m,
			ret:  true,
		},
		&LoadBoolStub{
			key:  int(10),
			args: m,
			ret:  true,
		},
	}
	for i, eachTest := range test {
		res, err := LoadBool(eachTest.key, eachTest.args)
		if i < trueCases {
			if err != nil || res != eachTest.ret {
				t.Error(i+1, ":", "Incorrect Match 1", err)
			}
			continue
		} else if res == eachTest.ret && err == nil {
			t.Error(i+1, ":", "Incorrect Match 2")
		}
	}
}

func TestAdjust(t *testing.T) {
	trueCases := 6
	m := map[string]interface{}{
		"args": map[string]interface{}{
			"auth": map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}

	test := []*AdjustStub{
		&AdjustStub{
			obj:   "args.auth.key1",
			state: m,
			ret:   "value1",
		},
		&AdjustStub{
			obj:   "args.auth.key3",
			state: m,
			ret:   "args.auth.key3",
		},
		&AdjustStub{
			obj:   int(10),
			state: m,
			ret:   int(10),
		},
		&AdjustStub{
			obj: map[string]interface{}{
				"args1": map[string]interface{}{
					"auth1": "args.auth.key1",
				},
			},
			state: m,
			ret: map[string]interface{}{
				"args1": map[string]interface{}{
					"auth1": "value1",
				},
			},
		},
		&AdjustStub{
			obj: map[string]interface{}{
				"args": map[string]interface{}{
					"auth": "args.auth1",
				},
			},
			state: m,
			ret: map[string]interface{}{
				"args": map[string]interface{}{
					"auth": "args.auth1",
				},
			},
		},
		&AdjustStub{
			obj: []interface{}{
				map[string]interface{}{
					"args": map[string]interface{}{
						"auth": "args.auth.key1",
					},
				},
				map[string]interface{}{
					"args": map[string]interface{}{
						"auth": "args.auth.key2",
					},
				},
			},
			state: m,
			ret: []interface{}{
				map[string]interface{}{
					"args": map[string]interface{}{
						"auth": "value1",
					},
				},
				map[string]interface{}{
					"args": map[string]interface{}{
						"auth": "value2",
					},
				},
			},
		},
		// False/Error Cases:
		&AdjustStub{
			obj:   "args.auth.key1",
			state: m,
			ret:   "val",
		},
		&AdjustStub{
			obj:   "args.auth.key3",
			state: m,
			ret:   "value1",
		},
		&AdjustStub{
			obj:   int(10),
			state: m,
			ret:   float64(10),
		},
		&AdjustStub{
			obj: map[string]interface{}{
				"args1": map[string]interface{}{
					"auth1": "args.auth.key3",
				},
			},
			state: m,
			ret: map[string]interface{}{
				"args1": map[string]interface{}{
					"auth1": "",
				},
			},
		},
		&AdjustStub{
			obj: []interface{}{
				map[string]interface{}{
					"args": map[string]interface{}{
						"auth": "args.auth",
					},
				},
				map[string]interface{}{
					"args": map[string]interface{}{
						"auth": "args.auth",
					},
				},
			},
			state: m,
			ret: []interface{}{
				map[string]interface{}{
					"args": map[string]interface{}{
						"auth": "value1",
					},
				},
				map[string]interface{}{
					"args": map[string]interface{}{
						"auth": "value2",
					},
				},
			},
		},
	}
	for i, eachTest := range test {
		res := Adjust(eachTest.obj, eachTest.state)
		eq := reflect.DeepEqual(eachTest.ret, res)
		if i < trueCases {
			if !eq {
				t.Error(i+1, ":", "Incorrect Match 1")
			}
			continue
		}
		if eq {
			t.Error(i+1, ":", "Incorrect Match 2")
		}
	}
}
