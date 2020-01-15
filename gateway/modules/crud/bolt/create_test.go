package bolt

import (
	"os"
	"testing"
)

func TestBolt_Create(t *testing.T) {

	b, err := Init(true, "embedded.db")
	if err != nil {
		t.Fatal("error initializing database")
	}

	for _, tt := range generateCreateTestCases() {
		t.Run(tt.name, func(t *testing.T) {

			got, err := b.Create(tt.args.ctx, tt.args.project, tt.args.col, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
	if err := os.Remove("embedded.db"); err != nil {
		t.Error("error removing database file")
	}
}
