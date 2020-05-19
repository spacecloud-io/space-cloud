package utils

import (
	"reflect"
	"testing"
	"time"
)

func TestStoreValue(t *testing.T) {
	type args struct {
		key   string
		value interface{}
		state map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "succesful test",
			args: args{
				key:   "a.b.c",
				value: 4,
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "succesful test [] in between",
			args: args{
				key:   "a.b[a.e].d",
				value: 4,
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": map[string]interface{}{
								"d": "ok",
							},
						},
						"e": "c",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "succesful test [] prefix",
			args: args{
				key:   "a.b[a.e]",
				value: 4,
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": map[string]interface{}{
								"d": "ok",
							},
						},
						"e": "c",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "len=0",
			args: args{
				key:   "h",
				value: 4,
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": map[string]interface{}{
								"d": "ok",
							},
						},
						"e": "c",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "succesful test",
			args: args{
				key:   "a.b[a.e].d",
				value: 4,
				state: map[string]interface{}{
					"aa": map[string]interface{}{
						"b": map[string]interface{}{
							"c": map[string]interface{}{
								"d": "ok",
							},
						},
						"e": "c",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "not map string interface error",
			args: args{
				key:   "a.b[a.e].d",
				value: 4,
				state: map[string]interface{}{
					"q": map[string]interface{}{
						"b": map[string]interface{}{
							"c": map[string]interface{}{
								"d": "ok",
							},
						},
						"e": "c",
					},
					"a": 1,
				},
			},
			wantErr: true,
		},
		{
			name: "convert create error 1",
			args: args{
				key:   "a.b[a.e].d",
				value: 4,
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"bw": map[string]interface{}{
							"c": map[string]interface{}{
								"d": "ok",
							},
						},
						"b": "c",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "convert create error 2",
			args: args{
				key:   "a.b[a.e]",
				value: 4,
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"bw": map[string]interface{}{
							"c": map[string]interface{}{
								"d": "ok",
							},
						},
						"b": "c",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "load error",
			args: args{
				key:   "a.b[.e]",
				value: 4,
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": map[string]interface{}{
								"d": "ok",
							},
						},
						"e": "c",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "load error",
			args: args{
				key:   "a.b[.e].d",
				value: 4,
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": map[string]interface{}{
								"d": "ok",
							},
						},
						"e": "c",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "subval not string",
			args: args{
				key:   "a.b[a.e].d",
				value: 4,
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": map[string]interface{}{
								"d": "ok",
							},
						},
						"e": 5,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "load error",
			args: args{
				key:   "a.b[a.e].d",
				value: 4,
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": 5,
						},
						"e": "c",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "convert create error 3",
			args: args{
				key:   "a.b.c.d",
				value: 4,
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": 6,
						},
						"e": "c",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "subval not string",
			args: args{
				key:   "a.b[a.e]",
				value: 4,
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": map[string]interface{}{
								"d": "ok",
							},
						},
						"e": 5,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StoreValue(tt.args.key, tt.args.value, tt.args.state); (err != nil) != tt.wantErr {
				t.Errorf("StoreValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_splitVariable(t *testing.T) {
	type args struct {
		key       string
		delimiter rune
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
		{
			name: "successful test",
			args: args{
				key:       "(op1).[op2]",
				delimiter: '.',
			},
			want: []string{"(op1)", "[op2]"},
		},
		{
			name: "test",
			args: args{
				key:       "(op1].(op2]",
				delimiter: '.',
			},
			want: []string{"(op1].(op2]"},
		},
		{
			name: "3op",
			args: args{
				key:       "args.abc[args.abc]",
				delimiter: '.',
			},
			want: []string{"args", "abc[args.abc]"},
		},
		{
			name: "3op",
			args: args{
				key:       "args.abc[args.abc].abc",
				delimiter: '.',
			},
			want: []string{"args", "abc[args.abc]", "abc"},
		},
		{
			name: "3op",
			args: args{
				key:       "utils.exist(args.abc)",
				delimiter: '.',
			},
			want: []string{"utils", "exist(args.abc)"},
		},
		{
			name: "3op",
			args: args{
				key:       "utils.exists(args.abc[args.abc].abc)",
				delimiter: '.',
			},
			want: []string{"utils", "exists(args.abc[args.abc].abc)"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitVariable(tt.args.key, tt.args.delimiter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitVariable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convert(t *testing.T) {
	type args struct {
		key string
		obj map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "successful test",
			args: args{
				key: "a",
				obj: map[string]interface{}{
					"a": map[string]interface{}{
						"b": 5,
					},
				},
			},
			want:    map[string]interface{}{"b": 5},
			wantErr: false,
		},
		{
			name: "key not present error",
			args: args{
				key: "a",
				obj: map[string]interface{}{
					"ab": map[string]interface{}{
						"b": 5,
					},
				},
			},
			//want:    map[string]interface{}{},
			wantErr: true,
		},
		{
			name: "wrong object",
			args: args{
				key: "a",
				obj: map[string]interface{}{
					"a": 3,
				},
			},
			//want:    map[string]interface{}{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convert(tt.args.key, tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadBool(t *testing.T) {
	type args struct {
		key  interface{}
		args map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "successful test",
			args: args{
				key: "a.b.c",
				args: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": true,
						},
					},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "load value error",
			args: args{
				key: "a.b.c",
				args: map[string]interface{}{
					"ab": map[string]interface{}{
						"b": map[string]interface{}{
							"c": true,
						},
					},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "successful test",
			args: args{
				key: "a.b.c",
				args: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": "true",
						},
					},
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadBool(tt.args.key, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LoadBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadNumber(t *testing.T) {
	type args struct {
		key  interface{}
		args map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "successful test",
			args: args{
				key: "a.b",
				args: map[string]interface{}{
					"a": map[string]interface{}{
						"b": 7,
					},
				},
			},
			want: 7,
		},
		{
			name: "successful test",
			args: args{
				key: "a.b",
				args: map[string]interface{}{
					"a": map[string]interface{}{
						"b": 7.0,
					},
				},
			},
			want: 7,
		},
		{
			name: "successful test",
			args: args{
				key: "a.b",
				args: map[string]interface{}{
					"a": map[string]interface{}{
						"b": int32(7),
					},
				},
			},
			want: 7,
		},
		{
			name: "successful test",
			args: args{
				key: "a.b",
				args: map[string]interface{}{
					"a": map[string]interface{}{
						"b": int64(7),
					},
				},
			},
			want: 7,
		},
		{
			name: "successful test",
			args: args{
				key: "a.b",
				args: map[string]interface{}{
					"a": map[string]interface{}{
						"b": "7",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "successful test",
			args: args{
				key: "a.b",
				args: map[string]interface{}{
					"ab": map[string]interface{}{
						"b": int32(7),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadNumber(tt.args.key, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LoadNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadValue(t *testing.T) {
	type args struct {
		key   string
		state map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "successful test",
			args: args{
				key: "a.b.c.d",
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": map[string]interface{}{
								"d": 5,
							},
						},
					},
				},
			},
			want:    5,
			wantErr: false,
		},
		{
			name: "successful test",
			args: args{
				key: "a.b.c.d",
				state: map[string]interface{}{
					"ab": map[string]interface{}{
						"b": map[string]interface{}{
							"c": map[string]interface{}{
								"d": 5,
							},
						},
					},
					"a": 3,
				},
			},
			//want:    5,
			wantErr: true,
		},
		{
			name: "utils testing",
			args: args{
				key: "utils.exists(a.b.c)",
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": 54,
						},
					},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "utils testing length arr",
			args: args{
				key: "utils.length(a.b.somearray)",
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"somearray": []interface{}{1, 2, 3},
						},
					},
				},
			},
			want:    int64(3),
			wantErr: false,
		},

		{
			name: "utils testing (not split)",
			args: args{
				key: "utils.exists(a.(b.c))",
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": 54,
						},
					},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "utils testing (not split)",
			args: args{
				key: "utils.exist(a.(b.c))",
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": 54,
						},
					},
				},
			},
			//want:    false,
			wantErr: true,
		},
		{
			name: "utils testing",
			args: args{
				key: "utils.exists(a.b[a.b.c])",
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": "c",
						},
					},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "pre:post of [ not present",
			args: args{
				key: "a.b[ab.b.c]",
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": "c",
						},
					},
				},
			},
			//want:    true,
			wantErr: true,
		},
		{
			name: "subkey not string",
			args: args{
				key: "a.b[a.b.c]",
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": 5,
						},
					},
				},
			},
			//want:    true,
			wantErr: true,
		},
		{
			name: "subkey not string",
			args: args{
				key: "a.b[a.b.c]",
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": "5",
						},
					},
				},
			},
			//want:    true,
			wantErr: true,
		},
		{
			name: "subkey not string",
			args: args{
				key: "a.b[a.b.c]",
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"bc": map[string]interface{}{
							"c": "5",
						},
						"b": 5,
					},
				},
			},
			//want:    true,
			wantErr: true,
		},
		{
			name: "utils testing",
			args: args{
				key: "utils.exists(a.b[a.e].d)",
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": map[string]interface{}{
								"d": "ok",
							},
						},
						"e": "c",
					},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "utils testing",
			args: args{
				key: "",
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": map[string]interface{}{
								"d": "ok",
							},
						},
						"e": "c",
					},
				},
			},
			//want:    true,
			wantErr: true,
		},
		{
			name: "utils testing",
			args: args{
				key: "utilsexists(a.b[a.e].d)",
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": map[string]interface{}{
								"d": "ok",
							},
						},
						"e": "c",
					},
				},
			},
			//want:    true,
			wantErr: true,
		},
		{
			name: "0:pre not map string interface",
			args: args{
				key: "a.b[a.e].d",
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"bv": map[string]interface{}{
							"c": map[string]interface{}{
								"d": "ok",
							},
						},
						"b": "c",
					},
				},
			},
			//want:    true,
			wantErr: true,
		},
		{
			name: "pre:post not map string interface",
			args: args{
				key: "a.b[ab.e].d",
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": map[string]interface{}{
								"d": "ok",
							},
						},
						"ab": "c",
					},
				},
			},
			//want:    true,
			wantErr: true,
		},
		{
			name: "subval not map string",
			args: args{
				key: "a.b[a.e].d",
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": map[string]interface{}{
								"d": "ok",
							},
						},
						"e": 5,
					},
				},
			},
			//want:    true,
			wantErr: true,
		},
		{
			name: "subval not map string",
			args: args{
				key: "a.b[a.e].d",
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"ch": map[string]interface{}{
								"d": "ok",
							},
						},
						"e": "c",
						"c": "hj",
					},
				},
			},
			// want:    true,
			wantErr: true,
		},
		{
			name: "utils.addDuration testing",
			args: args{
				key:   "utils.addDuration('utils.now()', '0h')",
				state: map[string]interface{}{},
			},
			want:    time.Now().UTC().Format(time.RFC3339),
			wantErr: false,
		},
		{
			name: "utils.addDuration testing",
			args: args{
				key:   "utils.addDuration('2020-01-01T00:00:00Z', '4h')",
				state: map[string]interface{}{},
			},
			want:    "2020-01-01T04:00:00Z",
			wantErr: false,
		},
		{
			name: "utils.roundUpDate testing",
			args: args{
				key:   "utils.roundUpDate('2020-03-25', 'month')",
				state: map[string]interface{}{},
			},
			want:    "2020-03-01T00:00:00Z",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadValue(tt.args.key, tt.args.state)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadStringIfExists(t *testing.T) {
	type args struct {
		value string
		state map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "sucesful test",
			args: args{
				value: "args.b",
				state: map[string]interface{}{
					"args": map[string]interface{}{
						"b": "5",
					},
				},
			},
			want:    "5",
			wantErr: false,
		},
		{
			name: "wrong prefix",
			args: args{
				value: "arggs.b",
				state: map[string]interface{}{
					"args": map[string]interface{}{
						"b": "5",
					},
				},
			},
			want: "arggs.b",
		},
		{
			name: "load value error",
			args: args{
				value: "args.",
				state: map[string]interface{}{},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "wrong prefix",
			args: args{
				value: "arggs.b",
				state: map[string]interface{}{
					"args": map[string]interface{}{
						"b": 5,
					},
				},
			},
			want: "arggs.b",
		},
		{
			name: "sucesful test",
			args: args{
				value: "args.b",
				state: map[string]interface{}{
					"args": map[string]interface{}{
						"b": 5,
					},
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadStringIfExists(tt.args.value, tt.args.state)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadStringIfExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LoadStringIfExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdjust(t *testing.T) {
	type args struct {
		obj   interface{}
		state map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		// TODO: Add test cases.
		{
			name: "successful string",
			args: args{
				obj: "a.b.c",
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": 5,
						},
					},
				},
			},
			want: 5,
		},
		{
			name: "unsuccessful string",
			args: args{
				obj: "a.b.d",
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": 5,
						},
					},
				},
			},
			want: "a.b.d",
		},
		{
			name: "successful map",
			args: args{
				obj: map[string]interface{}{
					"op1": "a.b.c",
				},
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": 5,
						},
					},
				},
			},
			want: map[string]interface{}{
				"op1": 5,
			},
		},
		{
			name: "successful []interface]",
			args: args{
				obj: []interface{}{"a.b.c"},
				state: map[string]interface{}{
					"a": map[string]interface{}{
						"b": map[string]interface{}{
							"c": 5,
						},
					},
				},
			},
			want: []interface{}{5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Adjust(tt.args.obj, tt.args.state); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Adjust() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertOrCreate(t *testing.T) {
	type args struct {
		k   string
		obj map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "successful test",
			args: args{
				k: "op1",
				obj: map[string]interface{}{
					"op1": map[string]interface{}{"op2": 2, "op3": 4},
				},
			},
			want:    map[string]interface{}{"op2": 2, "op3": 4},
			wantErr: false,
		},
		{
			name: "key not present",
			args: args{
				k: "op1",
				obj: map[string]interface{}{
					"op5": map[string]interface{}{"op2": 2, "op3": 4},
				},
			},
			want:    map[string]interface{}{},
			wantErr: false,
		},
		{
			name: "obj value not map string",
			args: args{
				k: "op1",
				obj: map[string]interface{}{
					"op1": 4,
				},
			},
			//want:    map[string]interface{}{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertOrCreate(tt.args.k, tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertOrCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertOrCreate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteValue(t *testing.T) {
	type args struct {
		key   string
		state map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "successful",
			args: args{
				key: "a.b.c",
				state: map[string]interface{}{
					"a": map[string]interface{}{"op2": 1, "b": map[string]interface{}{"c": 1}},
					"b": map[string]interface{}{"op2": 1},
				},
			},
			wantErr: false,
		},
		{
			name: "len 0",
			args: args{
				key: "abc",
				state: map[string]interface{}{
					"a": map[string]interface{}{"op2": 1, "b": map[string]interface{}{"c": 1}},
					"b": map[string]interface{}{"op2": 1},
				},
			},
			wantErr: true,
		},
		{
			name: "scope not present",
			args: args{
				key: "a.b.c",
				state: map[string]interface{}{
					"ab": map[string]interface{}{"op2": 1, "b": map[string]interface{}{"c": 1}},
					"b":  map[string]interface{}{"op2": 1},
				},
			},
			wantErr: true,
		},
		{
			name: "successful",
			args: args{
				key: "a.b.c",
				state: map[string]interface{}{
					"a": 1,
					"b": map[string]interface{}{"op2": 1},
				},
			},
			wantErr: true,
		},
		{
			name: "successful",
			args: args{
				key: "a.b.c",
				state: map[string]interface{}{
					"a": map[string]interface{}{"op2": 1, "bc": map[string]interface{}{"c": 1}},
					"b": map[string]interface{}{"op2": 1},
				},
			},
			wantErr: true,
		},
		{
			name: "successful",
			args: args{
				key: "a.b.c",
				state: map[string]interface{}{
					"a": map[string]interface{}{"op2": 1, "b": 1},
					"b": map[string]interface{}{"op2": 1},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteValue(tt.args.key, tt.args.state); (err != nil) != tt.wantErr {
				t.Errorf("DeleteValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStoreValueInObject(t *testing.T) {
	type args struct {
		key   string
		value interface{}
		obj   map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "successful",
			args: args{
				key:   "a.b.c",
				value: 1,
				obj: map[string]interface{}{
					"a": map[string]interface{}{"b": map[string]interface{}{}},
				},
			},
		},
		{
			name: "not present",
			args: args{
				key:   "a.b.c",
				value: 1,
				obj: map[string]interface{}{
					"a": map[string]interface{}{"bv": map[string]interface{}{}},
				},
			},
		},
		{
			name: "type erro",
			args: args{
				key:   "a.b.c",
				value: 1,
				obj: map[string]interface{}{
					"a": map[string]interface{}{"b": 5},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StoreValueInObject(tt.args.key, tt.args.value, tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("StoreValueInObject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
