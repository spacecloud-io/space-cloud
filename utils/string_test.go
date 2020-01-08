package utils

import (
	"testing"
)

func TestSingleLeading(t *testing.T) {
	type args struct {
		s  string
		ch string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases for double slash
		{
			name: "only /",
			args: args{
				s:  "",
				ch: "/",
			},
			want: "/",
		}, /*
			{
				name: "/a/b/d",
				args: args{
					s:  "//a/b//d",
					ch: "/",
				},
			},*/
		{
			name: "/////a/b/d////",
			args: args{
				s:  "/////a/b/d////",
				ch: "/",
			},
			want: "/a/b/d",
		},
		{
			name: "a/b/d///",
			args: args{
				s:  "a/b/d////",
				ch: "/",
			},
			want: "/a/b/d",
		},
		{
			name: "///a/b/d",
			args: args{
				s:  "////a/b/d",
				ch: "/",
			},
			want: "/a/b/d",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SingleLeading(tt.args.s, tt.args.ch); got != tt.want {
				t.Errorf("SingleLeading() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSingleTrailing(t *testing.T) {
	type args struct {
		s  string
		ch string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "only /",
			args: args{
				s:  "",
				ch: "/",
			},
			want: "/",
		}, /*
			{
				name: "/a/b/d",
				args: args{
					s:  "//a/b//d",
					ch: "/",
				},
			},*/
		{
			name: "/////a/b/d////",
			args: args{
				s:  "/////a/b/d////",
				ch: "/",
			},
			want: "a/b/d/",
		},
		{
			name: "a/b/d///",
			args: args{
				s:  "a/b/d////",
				ch: "/",
			},
			want: "a/b/d/",
		},
		{
			name: "///a/b/d",
			args: args{
				s:  "////a/b/d",
				ch: "/",
			},
			want: "a/b/d/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SingleTrailing(tt.args.s, tt.args.ch); got != tt.want {
				t.Errorf("SingleTrailing() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSingleLeadingTrailing(t *testing.T) {
	type args struct {
		s  string
		ch string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "only /",
			args: args{
				s:  "",
				ch: "/",
			},
			want: "/",
		}, /*
			{
				name: "/a/b/d",
				args: args{
					s:  "//a/b//d",
					ch: "/",
				},
			},*/
		{
			name: "/////a/b/d////",
			args: args{
				s:  "/////a/b/d////",
				ch: "/",
			},
			want: "/a/b/d/",
		},
		{
			name: "a/b/d///",
			args: args{
				s:  "a/b/d////",
				ch: "/",
			},
			want: "/a/b/d/",
		},
		{
			name: "///a/b/d",
			args: args{
				s:  "////a/b/d",
				ch: "/",
			},
			want: "/a/b/d/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SingleLeadingTrailing(tt.args.s, tt.args.ch); got != tt.want {
				t.Errorf("SingleLeadingTrailing() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJoinLeading(t *testing.T) {
	type args struct {
		s1 string
		s2 string
		ch string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "only /",
			args: args{
				s1: "",
				s2: "",
				ch: "/",
			},
			want: "/",
		}, /*
			{
				name: "/a/b/d",
				args: args{
					s:  "//a/b//d",
					ch: "/",
				},
			},*/
		{
			name: "/////a/b/d////",
			args: args{
				s1: "/////a/b/d////",
				s2: "/////c/////",
				ch: "/",
			},
			want: "/a/b/d/c",
		},
		{
			name: "a/b/d///",
			args: args{
				s1: "/////a/b/d////",
				s2: "/////c/////",
				ch: "/",
			},
			want: "/a/b/d/c",
		},
		{
			name: "///a/b/d",
			args: args{
				s1: "/////a/b/d////",
				s2: "/////c/////",
				ch: "/",
			},
			want: "/a/b/d/c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JoinLeading(tt.args.s1, tt.args.s2, tt.args.ch); got != tt.want {
				t.Errorf("JoinLeading() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJoinTrailing(t *testing.T) {
	type args struct {
		s1 string
		s2 string
		ch string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "only /",
			args: args{
				s1: "",
				s2: "",
				ch: "/",
			},
			want: "/",
		}, /*
			{
				name: "/a/b/d",
				args: args{
					s:  "//a/b//d",
					ch: "/",
				},
			},*/
		{
			name: "/////a/b/d////",
			args: args{
				s1: "/////a/b/d////",
				s2: "/////c/////",
				ch: "/",
			},
			want: "a/b/d/c/",
		},
		{
			name: "a/b/d///",
			args: args{
				s1: "/////a/b/d////",
				s2: "/////c/////",
				ch: "/",
			},
			want: "a/b/d/c/",
		},
		{
			name: "///a/b/d",
			args: args{
				s1: "/////a/b/d////",
				s2: "/////c/////",
				ch: "/",
			},
			want: "a/b/d/c/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JoinTrailing(tt.args.s1, tt.args.s2, tt.args.ch); got != tt.want {
				t.Errorf("JoinTrailing() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJoinLeadingTrailing(t *testing.T) {
	type args struct {
		s1 string
		s2 string
		ch string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "only /",
			args: args{
				s1: "",
				s2: "",
				ch: "/",
			},
			want: "/",
		}, /*
			{
				name: "/a/b/d",
				args: args{
					s:  "//a/b//d",
					ch: "/",
				},
			},*/
		{
			name: "/////a/b/d////",
			args: args{
				s1: "/////a/b/d////",
				s2: "/////c/////",
				ch: "/",
			},
			want: "/a/b/d/c/",
		},
		{
			name: "a/b/d///",
			args: args{
				s1: "/////a/b/d////",
				s2: "/////c/////",
				ch: "/",
			},
			want: "/a/b/d/c/",
		},
		{
			name: "///a/b/d",
			args: args{
				s1: "/////a/b/d////",
				s2: "/////c/////",
				ch: "/",
			},
			want: "/a/b/d/c/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JoinLeadingTrailing(tt.args.s1, tt.args.s2, tt.args.ch); got != tt.want {
				t.Errorf("JoinLeadingTrailing() = %v, want %v", got, tt.want)
			}
		})
	}
}
