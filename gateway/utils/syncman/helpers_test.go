package syncman

import "testing"

func Test_calcTokens(t *testing.T) {
	type args struct {
		n      int
		tokens int
		i      int
	}
	tests := []struct {
		name      string
		args      args
		wantStart int
		wantEnd   int
	}{
		{name: "test1", args: args{n: 7, tokens: 100, i: 0}, wantStart: 0, wantEnd: 14},
		{name: "test2", args: args{n: 7, tokens: 100, i: 4}, wantStart: 60, wantEnd: 74},
		{name: "test3", args: args{n: 7, tokens: 100, i: 5}, wantStart: 75, wantEnd: 89},
		{name: "test4", args: args{n: 7, tokens: 100, i: 6}, wantStart: 90, wantEnd: 99},
		{name: "test5", args: args{n: 1, tokens: 100, i: 0}, wantStart: 0, wantEnd: 99},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd := calcTokens(tt.args.n, tt.args.tokens, tt.args.i)
			if gotStart != tt.wantStart {
				t.Errorf("calcTokens() gotStart = %v, want %v", gotStart, tt.wantStart)
			}
			if gotEnd != tt.wantEnd {
				t.Errorf("calcTokens() gotEnd = %v, want %v", gotEnd, tt.wantEnd)
			}
		})
	}
}
