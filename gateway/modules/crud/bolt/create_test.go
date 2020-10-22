package bolt

import (
	"context"
	"os"
	"testing"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

func TestBolt_Create(t *testing.T) {

	b, err := Init(true, "create.db", "bucketName", nil)
	if err != nil {
		t.Fatal("error initializing database")
	}

	for _, tt := range generateCreateTestCases() {
		t.Run(tt.name, func(t *testing.T) {

			got, err := b.Create(context.Background(), tt.args.col, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
	utils.CloseTheCloser(b)
	if err := os.Remove("create.db"); err != nil {
		t.Error("error removing database file:", err)
	}
}
