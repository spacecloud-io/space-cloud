package amazons3

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	uuid "github.com/satori/go.uuid"
	"github.com/spaceuptech/space-cloud/model"
)

// ListDir lists a directory in S3
func (a *AmazonS3) ListDir(req *model.ListFilesRequest) ([]*model.ListFilesResponse, error) {
	svc := s3.New(a.client)

	// Add a backslash if not there already
	if !strings.HasSuffix(req.Path, "/") {
		req.Path = req.Path + "/"
	}

	resp, _ := svc.ListObjects(&s3.ListObjectsInput{
		Bucket:    aws.String(a.bucket),
		Prefix:    aws.String(req.Path), //backslash at the end is important
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

// ReadFile reads a file from S3
func (a *AmazonS3) ReadFile(path string) (*model.File, error) {
	u2 := uuid.NewV4()

	tmpfile, err := ioutil.TempFile("", u2.String())
	if err != nil {
		return nil, err
	}

	downloader := s3manager.NewDownloader(a.client)

	_, err = downloader.Download(tmpfile,
		&s3.GetObjectInput{
			Bucket: aws.String(a.bucket),
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
