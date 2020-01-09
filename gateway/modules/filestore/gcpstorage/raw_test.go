package gcpstorage

// func TestGCPStorage_DoesExists(t *testing.T) {
// 	type fields struct {
// 		client *storage.Client
// 		bucket string
// 	}
// 	type args struct {
// 		path string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   bool
// 	}{
// 		{
// 			name: "test",
// 			fields: fields{
// 				bucket: "gcpgolang",
// 			},
// 			args: args{path: "name/sirname/"},
// 			want: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			g, err := Init(tt.fields.bucket)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			if  err := g.DoesExists(tt.args.path); err != nil {
// 				t.Errorf("DoesExists() = %v, want %v " , tt.want, err)
// 			}
// 		})
// 	}
// }
