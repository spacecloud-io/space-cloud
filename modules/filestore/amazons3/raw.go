package amazons3

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (a *AmazonS3) DoesExists(path string) error {
	svc := s3.New(a.client)
	input := &s3.GetObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(path),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}
