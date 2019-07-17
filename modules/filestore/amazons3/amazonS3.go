package amazons3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spaceuptech/space-cloud/utils"
)

// AmazonS3 holds the S3 driver session
type AmazonS3 struct {
	client *session.Session
}

// Init initializes an amazon s3 driver
func Init(region, endpoint string) (*AmazonS3, error) {
	awsConf := &aws.Config{Region: aws.String(region)}
	if len(endpoint) > 0 {
		awsConf.Endpoint = aws.String(endpoint)
	}
	session, err := session.NewSession(awsConf)
	return &AmazonS3{client: session}, err
}

// GetStoreType returns the file store type
func (a *AmazonS3) GetStoreType() utils.FileStoreType {
	return utils.AmazonS3
}

// Close gracefully closed the local filestore module
func (a *AmazonS3) Close() error {
	return nil
}
