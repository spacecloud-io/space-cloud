package metrics

import (
	"sync"
	"testing"

	api "github.com/spaceuptech/space-api-go"
	"github.com/spaceuptech/space-api-go/db"

	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

func TestNew(t *testing.T) {
	type fields struct {
		isProd           bool
		clusterID        string
		nodeID           string
		projects         sync.Map
		isMetricDisabled bool
		sink             *db.DB
		adminMan         *admin.Manager
		syncMan          *syncman.Manager
	}
	type args struct {
		clusterID        string
		nodeID           string
		isMetricDisabled bool
		adminMan         *admin.Manager
		syncMan          *syncman.Manager
		isProd           bool
	}
	tests := []struct {
		name    string
		args    args
		want    *fields
		wantErr bool
	}{
		{
			name: "valid config provided",
			args: args{
				clusterID: "clusterID",
				nodeID:    "nodeID",
				adminMan:  &admin.Manager{},
				syncMan:   &syncman.Manager{},
			},
			want: &fields{
				isProd:           false,
				clusterID:        "clusterID",
				nodeID:           "nodeID",
				isMetricDisabled: false,
				adminMan:         &admin.Manager{},
				syncMan:          &syncman.Manager{},
				sink:             api.New("spacecloud", "localhost:4123", false).DB("db"),
			},
			wantErr: false,
		},
		{
			name: "valid config provided metrics disabled",
			args: args{
				isMetricDisabled: true,
			},
			want:    new(fields),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.clusterID, tt.args.nodeID, tt.args.isMetricDisabled, tt.args.adminMan, tt.args.syncMan, tt.args.isProd)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			testFuncIsEqual(t, got.isProd, tt.want.isProd)
			testFuncIsEqual(t, got.isMetricDisabled, tt.want.isMetricDisabled)
			testFuncIsEqual(t, got.clusterID, tt.want.clusterID)
			testFuncIsEqual(t, got.nodeID, tt.want.nodeID)
			testFuncIsEqual(t, got.adminMan, tt.want.adminMan)
			testFuncIsEqual(t, got.syncMan, tt.want.syncMan)
			testFuncIsEqual(t, &got.projects, &tt.want.projects)
			// isEqual(t, got.sink, tt.want.sink) unable to compare sink field
		})
	}
}
