package bolt

import (
	"context"
	"fmt"

	// "go.etcd.io/bbolt"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/utils"
)

type args struct {
	ctx context.Context
	col string
	req *model.CreateRequest
}
type creatTestData struct {
	name    string
	args    args
	want    int64
	wantErr bool
}

func generateCreateTestCases() []creatTestData {
	tests := []creatTestData{
		{
			name: "insert single document",
			want: 1,
			args: args{
				ctx: context.Background(),
				col: "project_details",
				req: &model.CreateRequest{
					Document: map[string]interface{}{
						"_id":           "1",
						"name":          "sharad",
						"team":          "admin",
						"project_count": 15,
						"isPrimary":     false,
						"project_details": map[string]interface{}{
							"project_name": "project1",
						},
					},
					Operation: utils.One,
				},
			},
		},
		{
			name: "insert multiple document",
			want: 3,
			args: args{
				ctx: context.Background(),
				col: "project_details",
				req: &model.CreateRequest{
					Document: []interface{}{
						map[string]interface{}{
							"_id":           "2",
							"name":          "jayesh",
							"project_count": float64(10),
							"team":          "admin",
							"isPrimary":     true,
							"project_details": map[string]interface{}{
								"project_name": "project1",
							},
						},
						map[string]interface{}{
							"_id":           "3",
							"name":          "noorain",
							"team":          "admin",
							"project_count": float64(52),
							"isPrimary":     true,
							"project_details": map[string]interface{}{
								"project_name": "project1",
							},
						},
						map[string]interface{}{
							"_id":           "4",
							"name":          "ali",
							"team":          "admin",
							"project_count": float64(100),
							"isPrimary":     true,
							"project_details": map[string]interface{}{
								"project_name": "project1",
							},
						},
					},
					Operation: utils.All,
				},
			},
		},
	}
	return tests
}

func createDatabaseWithTestData(b *Bolt) error {
	for _, tt := range generateCreateTestCases() {
		got, err := b.Create(tt.args.ctx, tt.args.col, tt.args.req)
		if (err != nil) != tt.wantErr {
			return fmt.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
		}
		if got != tt.want {
			return fmt.Errorf("Create() got = %v, want %v", got, tt.want)
		}
	}
	return nil
}
