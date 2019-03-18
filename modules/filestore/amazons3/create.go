package amazons3

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spaceuptech/space-cloud/model"
)

func (a *AmazonS3) CreateFile(ctx context.Context, project string, req *model.CreateFileRequest, file io.Reader) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(a.region),
	},
	)
	if err != nil {
		fmt.Println("AmazonS3 Couldn't Establish Connection ", err)
		return err
	}
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(project),
		Key:    aws.String(req.Path + "/" + req.Name),
		Body:   file,
	})
	return err
}

func (a *AmazonS3) CreateDir(ctx context.Context, project string, req *model.CreateFileRequest) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(a.region),
	},
	)
	if err != nil {
		fmt.Println("AmazonS3 Couldn't Establish Connection ", err)
		return err
	}

	path := req.Path
	// back slash at the end is important, if not then file will be created of that name
	if !strings.HasSuffix(path, "/") {
		path = req.Path + "/"
	}

	svc := s3.New(sess)
	request := &s3.PutObjectInput{
		Bucket: aws.String(project),
		Key:    aws.String(req.Path),
	}
	_, err = svc.PutObject(request)
	return err
	// return errors.New("Not Implemented")
}
