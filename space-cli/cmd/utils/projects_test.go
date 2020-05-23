package utils

import (
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cli/cmd/model"
)

func TestGetProjectsNamesFromArray(t *testing.T) {
	type args struct {
		projects []*model.Projects
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "get projects name properly",
			args: args{
				projects: []*model.Projects{
					{
						Name: "p1",
						ID:   "id1",
					},
					{
						Name: "p2",
						ID:   "id2",
					},
				},
			},
			want:    []string{"p1", "p2"},
			wantErr: false,
		},
		{
			name: "no projects provided",
			args: args{
				projects: []*model.Projects{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetProjectsNamesFromArray(tt.args.projects)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProjectsNamesFromArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProjectsNamesFromArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
