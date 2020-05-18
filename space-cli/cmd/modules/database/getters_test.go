package database

import (
	"reflect"
	"testing"

	"github.com/spaceuptech/space-cli/cmd/model"
)

func TestGetDbRule(t *testing.T) {
	type args struct {
		project     string
		commandName string
		params      map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    []*model.SpecObject
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDbRule(tt.args.project, tt.args.commandName, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDbRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDbRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDbConfig(t *testing.T) {
	type args struct {
		project     string
		commandName string
		params      map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    []*model.SpecObject
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDbConfig(tt.args.project, tt.args.commandName, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDbConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDbConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDbSchema(t *testing.T) {
	type args struct {
		project     string
		commandName string
		params      map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    []*model.SpecObject
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDbSchema(tt.args.project, tt.args.commandName, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDbSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDbSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}
