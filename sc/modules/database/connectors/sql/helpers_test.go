package sql

import (
	"fmt"
	"testing"
)

func Test_replaceSQLOperationWithPlaceHolder(t *testing.T) {
	type args struct {
		replace     string
		sqlString   string
		replaceWith func(value string) string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name: "Operation without place holder",
			args: args{
				replace:   "limit",
				sqlString: "select * from posts limit",
				replaceWith: func(value string) string {
					return fmt.Sprintf("limit %s rows", value)
				},
			},
			want:  "",
			want1: "select * from posts limit",
		},
		{
			name: "Operation with place holder value",
			args: args{
				replace:   "limit",
				sqlString: "select * from posts limit $1",
				replaceWith: func(value string) string {
					return fmt.Sprintf("limit %s rows", value)
				},
			},
			want:  "$1",
			want1: "select * from posts limit $1 rows",
		},
		{
			name: "Operation is between some text",
			args: args{
				replace:   "limit",
				sqlString: "select * from posts limit $1 order by age ASC",
				replaceWith: func(value string) string {
					return fmt.Sprintf("limit %s rows", value)
				},
			},
			want:  "$1",
			want1: "select * from posts limit $1 rows order by age ASC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := replaceSQLOperationWithPlaceHolder(tt.args.replace, tt.args.sqlString, tt.args.replaceWith)
			if got != tt.want {
				t.Errorf("replaceSQLOperationWithPlaceHolder() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("replaceSQLOperationWithPlaceHolder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
