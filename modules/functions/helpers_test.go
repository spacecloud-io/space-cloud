package functions

import "testing"

func Test_adjustPath(t *testing.T) {
	type args struct {
		path   string
		params interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "no params", args: args{path: "/abc"}, want: "/abc", wantErr: false},
		{name: "single param", args: args{path: "/abc/{args.p1}",
			params: map[string]interface{}{"p1": "xyz"}}, want: "/abc/xyz", wantErr: false},
		{name: "double params", args: args{path: "/abc/{args.p1}/def/{args.p2}",
			params: map[string]interface{}{"p1": "xyz", "p2": "cba"}}, want: "/abc/xyz/def/cba", wantErr: false},
		{name: "double params with float", args: args{path: "/abc/{args.p1}/def/{args.p2}",
			params: map[string]interface{}{"p1": 10.23, "p2": "cba"}}, want: "/abc/10.23/def/cba", wantErr: false},
		{name: "double params with int", args: args{path: "/abc/{args.p1}/def/{args.p2}",
			params: map[string]interface{}{"p1": 10.23, "p2": 20}}, want: "/abc/10.23/def/20", wantErr: false},
		{name: "invalid params", args: args{path: "/abc/{args.p1}/def/{args.p2}",
			params: map[string]interface{}{}}, wantErr: true},
		{name: "nil params", args: args{path: "/abc/{args.p1}/def/{args.p2}"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := adjustPath(tt.args.path, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("adjustPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("adjustPath() got = %v, want %v", got, tt.want)
			}
		})
	}
}
