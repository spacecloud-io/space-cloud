package amazons3

import (
	"context"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spaceuptech/space-cloud/model"
)

func (s3 *AmazonS3) CreateFile(ctx context.Context, project string, req *model.CreateFileRequest, file io.Reader) error {
	uploader := s3manager.NewUploader(s3.session)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(project),
		Key:    aws.String(req.Path + "/" + req.Name),
		Body:   file,
	})
	return err
}

func (s3 *AmazonS3) CreateDir(ctx context.Context, project string, req *model.CreateFileRequest) error {
	return errors.New("Not Yet Created")
}
