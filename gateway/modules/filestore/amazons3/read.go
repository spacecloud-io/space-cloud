package amazons3

import (
	"bufio"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	uuid "github.com/satori/go.uuid"
	"github.com/spaceuptech/helpers"

	"github.com/spaceuptech/space-cloud/gateway/model"
)

// ListDir lists a directory in S3
func (a *AmazonS3) ListDir(ctx context.Context, req *model.ListFilesRequest) ([]*model.ListFilesResponse, error) {
	svc := s3.New(a.client)

	req.Path = strings.TrimPrefix(req.Path, "/")
	// Add a backslash if not there already
	if !strings.HasSuffix(req.Path, "/") && len(req.Path) != 0 {
		req.Path = req.Path + "/"
	}

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:    aws.String(a.bucket),
		Prefix:    aws.String(req.Path), //backslash at the end is important
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Contents) == 0 {
		helpers.Logger.LogDebug(helpers.GetRequestID(ctx), "AWS list response is empty", map[string]interface{}{"path": req.Path, "type": req.Type})
		return []*model.ListFilesResponse{}, nil
	}

	result := make([]*model.ListFilesResponse, 0)

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
func (a *AmazonS3) ReadFile(ctx context.Context, path string) (*model.File, error) {
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
		_ = tmpfile.Close()
		_ = os.Remove(tmpfile.Name())
		return nil, err
	}
	return &model.File{File: bufio.NewReader(tmpfile), Close: func() error {
		defer func() { _ = os.Remove(tmpfile.Name()) }()
		return tmpfile.Close()
	}}, nil
}
