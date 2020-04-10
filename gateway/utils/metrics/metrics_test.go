package metrics

import (
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

func TestNew(t *testing.T) {
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
		want    *Module
		wantErr bool
	}{
		{
			name: "valid config provided",
			args: args{
				isMetricDisabled: true,
			},
			want:    new(Module),
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}
