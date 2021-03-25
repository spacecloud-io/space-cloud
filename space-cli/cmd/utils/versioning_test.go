package utils

import (
	"testing"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
)

func TestGetChartDownloadURL(t *testing.T) {
	type args struct {
		url     string
		version string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Right URL and version",
			args: args{
				url:     model.HelmSpaceCloudChartDownloadURL,
				version: "0.21.2",
			},
			want: "https://storage.googleapis.com/space-cloud/helm/space-cloud/space-cloud-0.21.2.tgz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHelmChartDownloadURL(tt.args.url, tt.args.version); got != tt.want {
				t.Errorf("GetHelmChartDownloadURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
