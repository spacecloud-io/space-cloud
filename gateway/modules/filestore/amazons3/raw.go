package amazons3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

// DoesExists checks if path exists
func (a *AmazonS3) DoesExists(path string) error {
	svc := s3.New(a.client)
	input := &s3.GetObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(path),
	}

	_, err := svc.GetObject(input)
	if err != nil {
		return fmt.Errorf("error getting svc object: %s", err.Error())
	}
	return nil
}

// GetState checks if sc is able to query s3
func (a *AmazonS3) GetState(ctx context.Context) error {
	err := a.DoesExists("/")
	if err != nil {
		if v, ok := err.(awserr.Error); ok {
			if v.Code() != s3.ErrCodeNoSuchKey {
				return fmt.Errorf("cannot read state - %s", err.Error())
			}
		}
	}
	return nil
}
