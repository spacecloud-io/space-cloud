package utils

import (
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cli/cmd/model"
)

func TestCreateSpecObject(t *testing.T) {
	type args struct {
		api     string
		objType string
		meta    map[string]string
		spec    interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *model.SpecObject
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "create spec object",
			args: args{
				api:     "api",
				objType: "type",
				meta:    map[string]string{"key": "value"},
				spec:    map[string]string{"key1": "value1"},
			},
			want: &model.SpecObject{
				API:  "api",
				Type: "type",
				Meta: map[string]string{"key": "value"},
				Spec: map[string]string{"key1": "value1"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateSpecObject(tt.args.api, tt.args.objType, tt.args.meta, tt.args.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSpecObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateSpecObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrintYaml(t *testing.T) {
	type args struct {
		objs []*model.SpecObject
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "print spec object",
			args: args{
				objs: []*model.SpecObject{
					{
						API:  "api",
						Type: "type",
						Meta: map[string]string{"key": "value"},
						Spec: map[string]string{"key1": "value1"},
					},
					{
						API:  "api1",
						Type: "type2",
						Meta: map[string]string{"key": "value"},
						Spec: map[string]string{"key1": "value1"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PrintYaml(tt.args.objs); (err != nil) != tt.wantErr {
				t.Errorf("PrintYaml() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
