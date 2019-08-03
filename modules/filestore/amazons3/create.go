package amazons3

import (
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spaceuptech/space-cloud/model"
)

// CreateFile creates a file in S3
func (a *AmazonS3) CreateFile(req *model.CreateFileRequest, file io.Reader) error {
	path := strings.Trim(req.Path, "/")
	name := strings.Trim(req.Name, "/")
	p := strings.Trim(path + "/" + name, "/")
	uploader := s3manager.NewUploader(a.client)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String("/" + p),
		Body:   file,
	})
	return err
}

// CreateDir creates a directory in S3
func (a *AmazonS3) CreateDir(req *model.CreateFileRequest) error {
	path := strings.Trim(req.Path, "/")
	name := strings.Trim(req.Name, "/")
	p := strings.Trim(path + "/" + name, "/")
	// back slash at the end is important, if not then file will be created of that name

	svc := s3.New(a.client)
	request := &s3.PutObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String("/" + p + "/"),
	}
	_, err := svc.PutObject(request)
	return err
}
