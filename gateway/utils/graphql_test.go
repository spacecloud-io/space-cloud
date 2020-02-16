package utils

import (
	"reflect"
	"testing"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
)

func TestParseGraphqlValue(t *testing.T) {
	type args struct {
		value ast.Value
		store M
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "converting string",
			args: args{
				value: &ast.StringValue{Kind: kinds.StringValue, Value: "some-value"},
			},
			want: "some-value",
		}, {
			name: "converting float",
			args: args{
				value: &ast.FloatValue{Kind: kinds.FloatValue, Value: "12.32"},
			},
			want: 12.32,
		}, {
			name: "converting float - incorrect input",
			args: args{
				value: &ast.FloatValue{Kind: kinds.FloatValue, Value: "wrong-input"},
			},
			wantErr: true,
		}, {
			name: "converting boolean ",
			args: args{
				value: &ast.BooleanValue{Kind: kinds.BooleanValue, Value: true},
			},
			want: true,
		}, {
			name: "converting int",
			args: args{
				value: &ast.IntValue{Kind: kinds.IntValue, Value: "12"},
			},
			want: 12,
		}, {
			name: "converting int- wrong-input",
			args: args{
				value: &ast.IntValue{Kind: kinds.IntValue, Value: "wrong-input"},
			},
			wantErr: true,
		}, {
			name: "converting variable",
			args: args{
				value: &ast.Variable{Kind: kinds.Variable, Name: &ast.Name{Value: "abc"}},
				store: M{"vars": map[string]interface{}{"abc": 1}},
			},
			want: 1,
		}, {
			name: "converting list",
			args: args{
				value: &ast.ListValue{Kind: kinds.ListValue, Values: []ast.Value{
					&ast.StringValue{Kind: kinds.StringValue, Value: "1"}, &ast.StringValue{Kind: kinds.StringValue, Value: "1234"},
				}},
			},
			want: []interface{}{"1", "1234"},
		}, {
			name: "converting enum",
			args: args{
				value: &ast.EnumValue{Kind: kinds.EnumValue, Value: "a.b.c"},
				store: M{"a": map[string]interface{}{"b": map[string]interface{}{"c": 69}}},
			},
			want: 69,
		}, {
			name: "converting enum- using __ ",
			args: args{
				value: &ast.EnumValue{Kind: kinds.EnumValue, Value: "a_b_c"},
			},
			want: "a_b_c",
		}, {
			name: "converting obj",
			args: args{
				value: &ast.ObjectValue{Kind: kinds.ObjectValue, Fields: []*ast.ObjectField{
					{Name: &ast.Name{Value: "op1"}, Value: &ast.StringValue{Kind: kinds.StringValue, Value: "a.b.c"}},
				}},
			},
			want: map[string]interface{}{"op1": "a.b.c"},
		}, {
			name: "converting string error",
			args: args{
				value: &ast.StringValue{Kind: kinds.StringValue, Value: "a.b.c"},
				store: M{"a": map[string]interface{}{"b": map[string]interface{}{"c": 69}, "cdbb": 2}},
			},
			want: 69,
		}, {
			name: "converting string __",
			args: args{
				value: &ast.StringValue{Kind: kinds.StringValue, Value: "a__b__c"},
				store: M{"a": map[string]interface{}{"b": map[string]interface{}{"c": 69}, "cdbb": 2}},
			},
			want: 69,
		}, {
			name: "converting enum",
			args: args{
				value: &ast.EnumValue{Kind: kinds.EnumValue, Value: "a__b__c"},
				store: M{"a": map[string]interface{}{"b": map[string]interface{}{"c": 69}}},
			},
			want: 69,
		}, {
			name: "converting list error",
			args: args{
				value: &ast.ListValue{Kind: kinds.ListValue, Values: []ast.Value{
					&ast.StringValue{Kind: kinds.StringValue, Value: "a.b.c"}, &ast.FloatValue{Kind: kinds.FloatValue, Value: "123df.4"},
				}},
				store: M{"a": map[string]interface{}{"ba": map[string]interface{}{"c": 69}, "b": 8}},
			},
			wantErr: true,
		}, {
			name: "default case",
			args: args{
				value: &ast.ListValue{Kind: kinds.ListValue, Values: []ast.Value{
					&ast.StringValue{Kind: kinds.StringValue, Value: "a.b.c"}, &ast.FloatValue{Kind: kinds.FloatValue, Value: "123df.4"},
				}},
				store: M{"a": map[string]interface{}{"ba": map[string]interface{}{"c": 69}, "b": 8}},
			},
			wantErr: true,
		}, {
			name: "converting obj",
			args: args{
				value: &ast.ObjectValue{Kind: kinds.ObjectValue, Fields: []*ast.ObjectField{
					{Name: &ast.Name{Value: "_op1"}, Value: &ast.StringValue{Kind: kinds.StringValue, Value: "a.b.c"}},
				}},
			},
			want: map[string]interface{}{"$op1": "a.b.c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseGraphqlValue(tt.args.value, tt.args.store)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseGraphqlValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseGraphqlValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
