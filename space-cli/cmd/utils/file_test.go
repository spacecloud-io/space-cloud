package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cli/cmd/model"
	"github.com/spaceuptech/space-cli/cmd/utils/file"
)

func TestAppendConfigToDisk(t *testing.T) {

	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}

	type args struct {
		specObj  *model.SpecObject
		filename string
	}

	tests := []struct {
		name           string
		schemaMockArgs []mockArgs
		args           args
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			name: "append config to non existing file",
			schemaMockArgs: []mockArgs{
				{
					method:         "Stat",
					args:           []interface{}{},
					paramsReturned: []interface{}{nil, nil},
				},
				{
					method:         "IsNotExist",
					args:           []interface{}{},
					paramsReturned: []interface{}{true},
				},
				{
					method:         "WriteFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{nil},
				},
			},
			args: args{
				specObj: &model.SpecObject{
					API:  "api",
					Type: "type",
					Meta: map[string]string{"key": "value"},
				},
				filename: "file123",
			},
			wantErr: false,
		},
		{
			name: "append config to existing file",
			schemaMockArgs: []mockArgs{
				{
					method:         "Stat",
					args:           []interface{}{},
					paramsReturned: []interface{}{nil, nil},
				},
				{
					method:         "IsNotExist",
					args:           []interface{}{},
					paramsReturned: []interface{}{false},
				},
				{
					method:         "IsDir",
					args:           []interface{}{},
					paramsReturned: []interface{}{false},
				},
				{
					method:         "OpenFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{&os.File{}, nil},
				},
				{
					method:         "Close",
					args:           []interface{}{},
					paramsReturned: []interface{}{nil},
				},
				{
					method:         "Write",
					args:           []interface{}{},
					paramsReturned: []interface{}{3, nil},
				},
				{
					method:         "WriteFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{nil},
				},
			},
			args: args{
				specObj: &model.SpecObject{
					API:  "api",
					Type: "type",
					Meta: map[string]string{"key": "value"},
				},
				filename: "file123",
			},
			wantErr: false,
		},
		{
			name: "error opening file",
			schemaMockArgs: []mockArgs{
				{
					method:         "Stat",
					args:           []interface{}{},
					paramsReturned: []interface{}{nil, nil},
				},
				{
					method:         "IsNotExist",
					args:           []interface{}{},
					paramsReturned: []interface{}{false},
				},
				{
					method:         "IsDir",
					args:           []interface{}{},
					paramsReturned: []interface{}{false},
				},
				{
					method:         "OpenFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{&os.File{}, fmt.Errorf("some-error")},
				},
			},
			args: args{
				specObj: &model.SpecObject{
					API:  "api",
					Type: "type",
					Meta: map[string]string{"key": "value"},
				},
				filename: "file123",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSchema := file.Mocket{}

			for _, m := range tt.schemaMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			file.File = &mockSchema

			if err := AppendConfigToDisk(tt.args.specObj, tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("AppendConfigToDisk() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReadSpecObjectsFromFile(t *testing.T) {

	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}

	type args struct {
		fileName string
	}

	tests := []struct {
		name           string
		schemaMockArgs []mockArgs
		args           args
		want           []*model.SpecObject
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			name: "proper",
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{[]byte("api: api1\ntype: type1\nmeta:\n  key: value\nspec:\n  key: value\n---\napi: api2\ntype: type2\nmeta:\n  key: value\nspec:\n  key: value"), nil},
				},
			},
			args: args{fileName: "test"},
			want: []*model.SpecObject{
				{
					API:  "api1",
					Type: "type1",
					Meta: map[string]string{"key": "value"},
					Spec: map[string]interface{}{"key": "value"},
				},
				{
					API:  "api2",
					Type: "type2",
					Meta: map[string]string{"key": "value"},
					Spec: map[string]interface{}{"key": "value"},
				},
			},
			wantErr: false,
		},
		{
			name: "cannot read data from file",
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{[]byte("api: api1\ntype: type1\nmeta:\n  key: value\nspec:\n  key: value\n---\napi: api2\ntype: type2\nmeta:\n  key: value\nspec:\n  key: value"), fmt.Errorf("some-error")},
				},
			},
			args:    args{fileName: "test"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "len of spec object is less than 5",
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{[]byte("a\n---\napi: api2\ntype: type2\nmeta:\n  key: value\nspec:\n  key: value"), nil},
				},
			},
			args: args{fileName: "test"},
			want: []*model.SpecObject{
				{
					API:  "api2",
					Type: "type2",
					Meta: map[string]string{"key": "value"},
					Spec: map[string]interface{}{"key": "value"},
				},
			},
			wantErr: false,
		},
		{
			name: "empty file content",
			schemaMockArgs: []mockArgs{
				{
					method:         "ReadFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{[]byte(""), nil},
				},
			},
			args:    args{fileName: "test"},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSchema := file.Mocket{}

			for _, m := range tt.schemaMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			file.File = &mockSchema

			got, err := ReadSpecObjectsFromFile(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadSpecObjectsFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				b, _ := json.MarshalIndent(got, "", " ")
				d, _ := json.MarshalIndent(tt.want, "", " ")
				t.Errorf("ReadSpecObjectsFromFile() = %s, want %v", string(b), string(d))

			}
		})
	}
}

func TestCreateFileIfNotExist(t *testing.T) {

	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}

	type args struct {
		path    string
		content string
	}
	tests := []struct {
		name           string
		schemaMockArgs []mockArgs
		args           args
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			name: "non existing file create properly",
			schemaMockArgs: []mockArgs{
				{
					method:         "Stat",
					args:           []interface{}{},
					paramsReturned: []interface{}{nil, nil},
				},
				{
					method:         "IsNotExist",
					args:           []interface{}{},
					paramsReturned: []interface{}{true},
				},
				{
					method:         "WriteFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{nil},
				},
			},
			args: args{
				path:    "path",
				content: "some-content",
			},
			wantErr: false,
		},
		{
			name: "file exists",
			schemaMockArgs: []mockArgs{
				{
					method:         "Stat",
					args:           []interface{}{},
					paramsReturned: []interface{}{nil, nil},
				},
				{
					method:         "IsNotExist",
					args:           []interface{}{},
					paramsReturned: []interface{}{false},
				},
				{
					method:         "WriteFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{nil},
				},
			},
			args: args{
				path:    "path",
				content: "some-content",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSchema := file.Mocket{}

			for _, m := range tt.schemaMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			file.File = &mockSchema

			if err := CreateFileIfNotExist(tt.args.path, tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("CreateFileIfNotExist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateConfigFile(t *testing.T) {

	type mockArgs struct {
		method         string
		args           []interface{}
		paramsReturned []interface{}
	}

	type args struct {
		path string
	}

	tests := []struct {
		name           string
		schemaMockArgs []mockArgs
		args           args
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			name: "non existing file create properly",
			schemaMockArgs: []mockArgs{
				{
					method:         "Stat",
					args:           []interface{}{},
					paramsReturned: []interface{}{nil, nil},
				},
				{
					method:         "IsNotExist",
					args:           []interface{}{},
					paramsReturned: []interface{}{true},
				},
				{
					method:         "WriteFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{nil},
				},
			},
			args: args{
				path: "path",
			},
			wantErr: false,
		},
		{
			name: "file exists",
			schemaMockArgs: []mockArgs{
				{
					method:         "Stat",
					args:           []interface{}{},
					paramsReturned: []interface{}{nil, nil},
				},
				{
					method:         "IsNotExist",
					args:           []interface{}{},
					paramsReturned: []interface{}{false},
				},
				{
					method:         "WriteFile",
					args:           []interface{}{},
					paramsReturned: []interface{}{nil},
				},
			},
			args: args{
				path: "path",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSchema := file.Mocket{}

			for _, m := range tt.schemaMockArgs {
				mockSchema.On(m.method, m.args...).Return(m.paramsReturned...)
			}

			file.File = &mockSchema

			if err := CreateConfigFile(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("CreateConfigFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
