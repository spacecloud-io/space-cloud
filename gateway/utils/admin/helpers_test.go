package admin

import (
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
)

func TestManager_parseToken(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		m       *Manager
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			m:    &Manager{admin: &config.Admin{Secret: "abcde"}},
			args : args{token:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.EKd8mqslm3YqO5cdfIF7mAkP6mdXrazy-hGK_SkJJDc"},
			want: map[string]interface{}
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.parseToken(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.parseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.parseToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
