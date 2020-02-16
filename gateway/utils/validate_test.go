package utils

import "testing"

func TestValidate(t *testing.T) {
	type args struct {
		where map[string]interface{}
		obj   interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "wrong object format",
			args: args{
				where: map[string]interface{}{"op1": 1},
				obj:   "yiu",
			},
			want: false,
		},
		{
			name: "wrong object format",
			args: args{
				where: map[string]interface{}{"op1": 1.5},
				obj:   1,
			},
			want: false,
		},
		{
			name: "wrong where",
			args: args{
				where: map[string]interface{}{"$or": []interface{}{map[string]interface{}{"op2": "1"}, map[string]interface{}{"op3": "2"}}},
				obj:   map[string]interface{}{"op1": 1},
			},
			want: false,
		},
		{
			name: "valid $or",
			args: args{
				where: map[string]interface{}{"$or": []interface{}{map[string]interface{}{"op2": map[string]interface{}{"$eq": 1}}}},
				obj:   map[string]interface{}{"op2": 1},
			},
			want: true,
		},
		{
			name: "valid $or",
			args: args{
				where: map[string]interface{}{"$or": map[string]interface{}{"op2": map[string]interface{}{"$eq": 1}}},
				obj:   map[string]interface{}{"op2": 1},
			},
			want: false,
		},
		{
			name: "test4",
			args: args{
				where: map[string]interface{}{"op2": []interface{}{map[string]interface{}{"op2": "1"}, map[string]interface{}{"op3": "2"}}},
				obj:   map[string]interface{}{"op2": 1},
			},
			want: false,
		},
		{
			name: "test5",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"op2": "1"}, "op3": []interface{}{map[string]interface{}{"op2": "1"}, map[string]interface{}{"op3": "2"}}},
				obj:   map[string]interface{}{"op2": 1},
			},
			want: false,
		},
		{
			name: "test6",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"op2": 1}, "op3": []interface{}{map[string]interface{}{"op2": "1"}, map[string]interface{}{"op3": "2"}}},
				obj:   map[string]interface{}{"op2": 1},
			},
			want: false,
		},
		{
			name: "test7",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$eq": 1}, "op3": []interface{}{map[string]interface{}{"op2": "1"}, map[string]interface{}{"op3": "2"}}},
				obj:   map[string]interface{}{"op2": 1},
			},
			want: false,
		},
		{
			name: "test8",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$eq": 1}, "op3": []interface{}{map[string]interface{}{"op2": "1"}, map[string]interface{}{"op3": "2"}}},
				obj:   map[string]interface{}{"op2": 1},
			},
			want: false,
		},
		{
			name: "valid eq",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$eq": 1}},
				obj:   map[string]interface{}{"op2": 1},
			},
			want: true,
		},
		{
			name: "compare",
			args: args{
				where: map[string]interface{}{"op2": "1"},
				obj:   map[string]interface{}{"op2": 1},
			},
			want: false,
		},
		{
			name: "compare2",
			args: args{
				where: map[string]interface{}{"op2": int64(1)},
				obj:   map[string]interface{}{"op2": int64(1)},
			},
			want: true,
		},
		{
			name: "invalid eq",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$eq": 3}},
				obj:   map[string]interface{}{"op2": 1},
			},
			want: false,
		},
		{
			name: "valid ne",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$ne": 2}},
				obj:   map[string]interface{}{"op2": 1},
			},
			want: true,
		},
		{
			name: "invalid ne",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$ne": 1}},
				obj:   map[string]interface{}{"op2": 1},
			},
			want: false,
		},
		{
			name: "invalid gt",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$gt": 2}},
				obj:   map[string]interface{}{"op2": 1},
			},
			want: false,
		},
		{
			name: "invalid gt",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$gt": "2"}},
				obj:   map[string]interface{}{"op2": "1"},
			},
			want: false,
		},
		{
			name: "valid gt(int64)",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$gt": 0}},
				obj:   map[string]interface{}{"op2": 1},
			},
			want: true,
		},
		{
			name: "valid gt(string)",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$gt": "0"}},
				obj:   map[string]interface{}{"op2": "1"},
			},
			want: true,
		},
		{
			name: "invalid gt(float)",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$gt": 0.7}},
				obj:   map[string]interface{}{"op2": 0.7},
			},
			want: false,
		},
		{
			name: "valid gte(float)",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$gte": 0.7}},
				obj:   map[string]interface{}{"op2": 0.7},
			},
			want: true,
		},
		{
			name: "valid gte(float)",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$gte": 0.6}},
				obj:   map[string]interface{}{"op2": 0.7},
			},
			want: true,
		},
		{
			name: "valid gte(string)",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$gte": "0.6"}},
				obj:   map[string]interface{}{"op2": "0.7"},
			},
			want: true,
		},
		{
			name: "invalid gte(string)",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$gte": "0.8"}},
				obj:   map[string]interface{}{"op2": "0.7"},
			},
			want: false,
		},
		{
			name: "invalid gte(string)",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$gte": 0.8}},
				obj:   map[string]interface{}{"op2": 0.7},
			},
			want: false,
		},
		{
			name: "invalid lt(string)",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$lt": "0.6"}},
				obj:   map[string]interface{}{"op2": "0.7"},
			},
			want: false,
		},
		{
			name: "valid lt(int)",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$lt": "9"}},
				obj:   map[string]interface{}{"op2": "7"},
			},
			want: true,
		},
		{
			name: "valid lt(float64)",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$lt": 9.5}},
				obj:   map[string]interface{}{"op2": 7.7},
			},
			want: true,
		},
		{
			name: "invalid lt(float64)",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$lt": 6.5}},
				obj:   map[string]interface{}{"op2": 7.7},
			},
			want: false,
		},
		{
			name: "valid lt(default)",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$lt": []interface{}{"ju5", "uiy"}}},
				obj:   map[string]interface{}{"op2": "j7jh"},
			},
			want: false,
		},
		{
			name: "valid lte(string)",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$lte": "9"}},
				obj:   map[string]interface{}{"op2": "7"},
			},
			want: true,
		},
		{
			name: "valid lte(string)",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$lte": int32(9)}},
				obj:   map[string]interface{}{"op2": int32(7)},
			},
			want: true,
		},
		{
			name: "invalid lte(float64)",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$lte": 6.9}},
				obj:   map[string]interface{}{"op2": 7},
			},
			want: false,
		},
		{
			name: "invalid lte(string)",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$lte": "6"}},
				obj:   map[string]interface{}{"op2": "7"},
			},
			want: false,
		},
		{
			name: "valid lte(int)",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$lte": "7"}},
				obj:   map[string]interface{}{"op2": "7"},
			},
			want: true,
		},
		{
			name: "valid lte(int)",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$lte": "7"}},
				obj:   []interface{}{map[string]interface{}{"op2": "7"}, map[string]interface{}{"op2": "6"}},
			},
			want: true,
		},
		{
			name: "invalid lte(int)",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$lte": 5.7}},
				obj:   []interface{}{map[string]interface{}{"op2": "7"}, map[string]interface{}{"op2": 6.0}},
			},
			want: false,
		},
		{
			name: "invalid parameter",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$ltfe": 5.7}},
				obj:   []interface{}{map[string]interface{}{"op2": "7"}, map[string]interface{}{"op2": 6.0}},
			},
			want: false,
		},
		{
			name: "valid lte(int)",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$lte": true}},
				obj:   map[string]interface{}{"op2": true},
			},
			want: true,
		},
		{
			name: "valid regex(prefix)",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$regex": "^sharad/"}},
				obj:   map[string]interface{}{"op2": "sharad/regoti"},
			},
			want: true,
		},
		{
			name: "invalid regex(prefix)",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$regex": "^sharad/"}},
				obj:   map[string]interface{}{"op2": "extra/sharad/regoti"},
			},
			want: false,
		},
		{
			name: "valid regex(contains)",
			args: args{
				where: map[string]interface{}{"op2": map[string]interface{}{"$regex": "/sharad/"}},
				obj:   map[string]interface{}{"op2": "extra/sharad/regoti"},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Validate(tt.args.where, tt.args.obj); got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
