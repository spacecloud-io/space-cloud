package metrics

import (
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/utils"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

func TestModule_AddDBOperation(t *testing.T) {
	type args struct {
		project string
		dbType  string
		col     string
		count   int64
		op      utils.OperationType
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	m, _ := New("", "", false, admin.New("clusterID", &config.AdminUser{}), &syncman.Manager{}, false)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.AddDBOperation(tt.args.project, tt.args.dbType, tt.args.col, tt.args.count, tt.args.op)
		})
	}
}

func TestModule_AddEventingType(t *testing.T) {
	type args struct {
		project      string
		eventingType string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "valid case",
			args: args{
				project:      "projectID",
				eventingType: "type",
			},
		},
	}
	m, _ := New("", "", false, admin.New("clusterID", &config.AdminUser{}), &syncman.Manager{}, false)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.AddEventingType(tt.args.project, tt.args.eventingType)
		})
	}
}

func TestModule_AddFileOperation(t *testing.T) {
	type args struct {
		project   string
		storeType string
		op        utils.OperationType
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	m, _ := New("", "", false, admin.New("clusterID", &config.AdminUser{}), &syncman.Manager{}, false)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.AddFileOperation(tt.args.project, tt.args.storeType, tt.args.op)
		})
	}
}

func TestModule_AddFunctionOperation(t *testing.T) {
	type args struct {
		project  string
		service  string
		function string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	m, _ := New("", "", false, admin.New("clusterID", &config.AdminUser{}), &syncman.Manager{}, false)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.AddFunctionOperation(tt.args.project, tt.args.service, tt.args.function)
		})
	}
}

func TestModule_LoadMetrics(t *testing.T) {
	tests := []struct {
		name string
		want []interface{}
	}{
		// TODO: Add test cases.
	}
	m, _ := New("", "", false, admin.New("clusterID", &config.AdminUser{}), &syncman.Manager{}, false)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.LoadMetrics(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
