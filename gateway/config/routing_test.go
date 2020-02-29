package config

import (
	"reflect"
	"testing"
)

func TestRoute_SelectTarget(t *testing.T) {
	type fields struct {
		ID      string
		Source  RouteSource
		Targets []RouteTarget
	}
	type args struct {
		weight float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    RouteTarget
		wantErr bool
	}{
		{
			name:    "valid case",
			fields:  fields{Targets: []RouteTarget{{Weight: 1}}},
			args:    args{weight: 0.3},
			want:    RouteTarget{Weight: 1},
			wantErr: false,
		}, {
			name:    "valid case - select 1st",
			fields:  fields{Targets: []RouteTarget{{Host: "1", Weight: 0.4}, {Host: "2", Weight: 0.3}, {Host: "3", Weight: 0.3}}},
			args:    args{weight: 0.2},
			want:    RouteTarget{Host: "1", Weight: 0.4},
			wantErr: false,
		}, {
			name:    "valid case - select 2nd",
			fields:  fields{Targets: []RouteTarget{{Host: "1", Weight: 0.4}, {Host: "2", Weight: 0.3}, {Host: "3", Weight: 0.3}}},
			args:    args{weight: 0.7},
			want:    RouteTarget{Host: "2", Weight: 0.3},
			wantErr: false,
		}, {
			name:    "valid case - select 2nd",
			fields:  fields{Targets: []RouteTarget{{Host: "1", Weight: 0.4}, {Host: "2", Weight: 0.2}, {Host: "3", Weight: 0.4}}},
			args:    args{weight: 1},
			want:    RouteTarget{Host: "3", Weight: 0.4},
			wantErr: false,
		}, {
			name:    "no routes provided",
			fields:  fields{Targets: []RouteTarget{}},
			wantErr: true,
		}, {
			name:    "weights don't add up to one",
			args:    args{weight: 1},
			fields:  fields{Targets: []RouteTarget{{Host: "1", Weight: 0.2}, {Host: "2", Weight: 0.1}, {Host: "3", Weight: 0.02}}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Route{
				ID:      tt.fields.ID,
				Source:  tt.fields.Source,
				Targets: tt.fields.Targets,
			}
			got, err := r.SelectTarget(tt.args.weight)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectTarget() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SelectTarget() got = %v, want %v", got, tt.want)
			}
		})
	}
}
