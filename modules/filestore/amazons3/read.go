package amazons3

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	uuid "github.com/satori/go.uuid"
	"github.com/spaceuptech/space-cloud/model"
)

func (a *AmazonS3) ListDir(ctx context.Context, project string, req *model.ListFilesRequest) ([]*model.ListFilesResponse, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(a.region),
	},
	)
	if err != nil {
		fmt.Println("AmazonS3 Couldn't Establish Connection ", err)
		return nil, err
	}
	svc := s3.New(sess)

	resp, _ := svc.ListObjects(&s3.ListObjectsInput{
		Bucket:    aws.String(project),
		Prefix:    aws.String(req.Path), //backslach at the end is important
		Delimiter: aws.String("/"),
	})

	result := []*model.ListFilesResponse{}
	resp.Contents = resp.Contents[1:]

	for _, key := range resp.Contents {
		t := &model.ListFilesResponse{Name: filepath.Base(*key.Key), Type: "file"}
		if req.Type == "all" || req.Type == t.Type {
			result = append(result, t)
		}
	}

	for _, key := range resp.CommonPrefixes {
		t := &model.ListFilesResponse{Name: filepath.Base(*key.Prefix), Type: "dir"}
		if req.Type == "all" || req.Type == t.Type {
			result = append(result, t)
		}
	}
	return result, nil
}

func (a *AmazonS3) ReadFile(ctx context.Context, project, path string) (*model.File, error) {
	u2 := uuid.NewV4()

	tmpfile, err := ioutil.TempFile("", u2.String())
	if err != nil {
		return nil, err
	}
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(a.region),
	},
	)
	if err != nil {
		fmt.Println("AmazonS3 Couldn't Establish Connection ", err)
		return nil, err
	}
	downloader := s3manager.NewDownloader(sess)

	_, err = downloader.Download(tmpfile,
		&s3.GetObjectInput{
			Bucket: aws.String(project),
			Key:    aws.String(path),
		})
	if err != nil {
		return nil, err
	}
	return &model.File{File: bufio.NewReader(tmpfile), Close: func() error {
		defer os.Remove(tmpfile.Name())
		return tmpfile.Close()
	}}, nil
}
