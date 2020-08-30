package amazons3

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spaceuptech/helpers"
)

// DoesExists checks if path exists
func (a *AmazonS3) DoesExists(ctx context.Context, path string) error {
	path = strings.TrimPrefix(path, "/")

	svc := s3.New(a.client)
	input := &s3.GetObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(path),
	}

	_, err := svc.GetObject(input)
	if err != nil {
		return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to get specified object from Amazon s3", err, nil)
	}
	return nil
}

// GetState checks if sc is able to query s3
func (a *AmazonS3) GetState(ctx context.Context) error {
	err := a.DoesExists(ctx, "/")
	if err != nil {
		if v, ok := err.(awserr.Error); ok {
			if v.Code() != s3.ErrCodeNoSuchKey {
				return helpers.Logger.LogError(helpers.GetRequestID(ctx), "Unable to connect to Amazon s3", err, nil)
			}
		}
	}
	return nil
}
