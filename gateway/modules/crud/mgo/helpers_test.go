package mgo

import (
	"context"
	"reflect"
	"testing"
)

func Test_sanitizeWhereClause(t *testing.T) {
	type args struct {
		ctx  context.Context
		col  string
		find map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "Where clause with collection name separated by dot notation",
			args: args{
				ctx: context.Background(),
				col: "users",
				find: map[string]interface{}{
					"users.id":             1,
					"users.address.street": "washington DC",
					"posts.id":             1,
					"users.name":           "same",
					"users.height":         5.5,
					"users.isUnderAge":     true,
					"users.posts": map[string]interface{}{
						"users.postId": 11,
						"users.views": map[string]interface{}{
							"users.viewCount": 100,
							"views.viewCount": 100,
						},
					},
					"$or": []interface{}{
						map[string]interface{}{
							"users.posts": map[string]interface{}{
								"users.postId": 11,
								"users.views": map[string]interface{}{
									"users.viewCount": 100,
									"views.viewCount": 100,
								},
							},
						},
					},
				},
			},
			want: map[string]interface{}{
				"id":             1,
				"address.street": "washington DC",
				"posts.id":       1,
				"name":           "same",
				"height":         5.5,
				"isUnderAge":     true,
				"posts": map[string]interface{}{
					"postId": 11,
					"views": map[string]interface{}{
						"viewCount":       100,
						"views.viewCount": 100,
					},
				},
				"$or": []interface{}{
					map[string]interface{}{
						"posts": map[string]interface{}{
							"postId": 11,
							"views": map[string]interface{}{
								"viewCount":       100,
								"views.viewCount": 100,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if sanitizeWhereClause(tt.args.ctx, tt.args.col, tt.args.find); !reflect.DeepEqual(tt.args.find, tt.want) {
				t.Errorf("sanitizeWhereClause() = %v, want %v", tt.args.find, tt.want)
			}
		})
	}
}
