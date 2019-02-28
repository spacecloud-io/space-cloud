package amazons3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func (a *AmazonS3) DeleteFile(ctx context.Context, project, path string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(a.region),
	},
	)
	if err != nil {
		fmt.Println("AmazonS3 Couldn't Establish Connection ", err)
		return err
	}
	svc := s3.New(sess)
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(project), Key: aws.String(path)})
	if err != nil {
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(project),
		Key:    aws.String(project + path),
	})
	return err
}

func (a *AmazonS3) DeleteDir(ctx context.Context, project, path string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(a.region),
	},
	)
	if err != nil {
		fmt.Println("AmazonS3 Couldn't Establish Connection ", err)
		return err
	}
	// TODO: Consider AWS operation limit
	svc := s3.New(sess)

	// Setup BatchDeleteIterator to iterate through a list of objects.
	iter := s3manager.NewDeleteListIterator(svc, &s3.ListObjectsInput{
		Bucket: aws.String(project),
		Prefix: aws.String(path),
	})

	// Traverse iterator deleting each object
	if err := s3manager.NewBatchDeleteWithClient(svc).Delete(aws.BackgroundContext(), iter); err != nil {
		fmt.Println("Unable to delete objects from bucket %q, %v", project, err)
	}
	return err
}
