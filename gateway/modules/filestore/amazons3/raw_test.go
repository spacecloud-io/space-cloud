package amazons3

// For this test to pass save your credentials in file named credentials
// save the above file at .aws/ directory in home folder of linux
// *****File contents*******
// [default]
// aws_access_key_id=<your access key>
// aws_secret_access_key=<your secret key>

// func TestAmazonS3_DoesExists(t *testing.T) {
// 	type args struct {
// 		path string
// 	}
// 	tests := []struct {
// 		name          string
// 		args          args
// 		isErrExpected bool
// 	}{
// 		{
// 			name:          "Directory directory1 exists",
// 			args:          args{path: "directory1/"},
// 			isErrExpected: false,
// 		},
// 		{
// 			name:          "Directory directory5 does not exists",
// 			args:          args{path: "directory5/"},
// 			isErrExpected: true,
// 		},
// 		{
// 			name:          "Root file Screenshot from 2019-12-12 10-42-50.png exists",
// 			args:          args{path: "Screenshot from 2019-12-12 10-42-50.png"},
// 			isErrExpected: false,
// 		},
// 		{
// 			name:          "Root file Screenshot from sdfsdfds 2019-12-12 10-42-50.png does not exists",
// 			args:          args{path: "Screenshot from sdfsdfds 2019-12-12 10-42-50.png"},
// 			isErrExpected: true,
// 		},
// 	}

// 	// Asia Pacific (Mumbai) 	ap-south-1 	rds.ap-south-1.amazonaws.com 	HTTPS	region := "ap-south-1"
// 	region := "ap-south-1"
// 	bucketName := "space-cloud-s3-test"
// 	a, err := Init(region, "", bucketName)
// 	if err != nil {
// 		t.Fatal("Error initializing amazon s3 instance", err)
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			got := a.DoesExists(tt.args.path)
// 			if tt.isErrExpected {
// 				if got == nil {
// 					t.Errorf("DoesExists() isExpected = %v got = %v", tt.isErrExpected, got)
// 					return
// 				}
// 			} else {
// 				if got != nil {
// 					t.Errorf("DoesExists() isExpected = %v got = %v", tt.isErrExpected, got)
// 					return
// 				}
// 			}
// 		})
// 	}
// }
