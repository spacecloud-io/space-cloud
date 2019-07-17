package amazons3

import (
	"context"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spaceuptech/space-cloud/model"
)

// CreateFile creates a file in S3
func (a *AmazonS3) CreateFile(ctx context.Context, project string, req *model.CreateFileRequest, file io.Reader) error {
	uploader := s3manager.NewUploader(a.client)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(project),
		Key:    aws.String(req.Path + "/" + req.Name),
		Body:   file,
	})
	return err
}

// CreateDir creates a directory in S3
func (a *AmazonS3) CreateDir(ctx context.Context, project string, req *model.CreateFileRequest) error {
	path := req.Path
	// back slash at the end is important, if not then file will be created of that name
	if !strings.HasSuffix(path, "/") {
		path = req.Path + "/"
	}

	svc := s3.New(a.client)
	request := &s3.PutObjectInput{
		Bucket: aws.String(project),
		Key:    aws.String(req.Path),
	}
	_, err := svc.PutObject(request)
	return err
	// return errors.New("Not Implemented")
}
