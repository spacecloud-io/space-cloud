package amazons3

import "testing"

func TestAmazonS3_DoesExists(t *testing.T) {

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "test",
			args:    args{path: ""},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := Init("", "", "")
			if err != nil {
				t.Fatal(err)
			}
			err = a.DoesExists(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoesExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
