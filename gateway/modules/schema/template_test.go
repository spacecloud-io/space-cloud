package schema

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spaceuptech/space-cloud/gateway/model"
)

func Test_generateSDL(t *testing.T) {
	type args struct {
		schemaCol model.Collection
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Successful test",
			args: args{
				schemaCol: model.Collection{"table1": model.Fields{"col2": &model.FieldType{FieldName: "col2", Kind: model.TypeID}}},
			},
			want:    "type  table1 { \n\tcol2: ID     \n}",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateSDL(tt.args.schemaCol)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateSDL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("generateSDL() = %v-%v, want %v-%v", got, reflect.TypeOf(got), tt.want, reflect.TypeOf(tt.want))
			}
		})
	}
}
