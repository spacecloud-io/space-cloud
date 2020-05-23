package utils

import (
	"testing"
)

func TestLogInfo(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{message: "msg1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LogInfo(tt.args.message)
		})
	}
}

func TestLogDebug(t *testing.T) {
	type args struct {
		message     string
		extraFields map[string]interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "no extra field",
			args: args{
				message: "msg",
			},
		},
		{
			name: "with extra field",
			args: args{
				message:     "msg",
				extraFields: map[string]interface{}{"key": "value"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LogDebug(tt.args.message, tt.args.extraFields)
		})
	}
}

func TestSetLogLevel(t *testing.T) {
	type args struct {
		loglevel string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "debug",
			args: args{
				loglevel: "debug",
			},
		},
		{
			name: "info",
			args: args{
				loglevel: "info",
			},
		},
		{
			name: "error",
			args: args{
				loglevel: "error",
			},
		},
		{
			name: "random",
			args: args{
				loglevel: "random",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLogLevel(tt.args.loglevel)
		})
	}
}
