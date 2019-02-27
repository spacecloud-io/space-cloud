package amazons3

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spaceuptech/space-cloud/utils"
)

type AmazonS3 struct {
	tempFileName string
	session      *session.Session
}

func Init(region string) (*AmazonS3, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		fmt.Println("AmazonS3 Couldn't Establish Connection ", err)
		return nil, err
	}
	_, err = sess.Config.Credentials.Get()

	if err != nil {
		fmt.Println("AmazonS3 Credentials Not Found ", err)
		return nil, err
	}
	return &AmazonS3{"", sess}, nil
}

// GetStoreType returns the file store type
func (a *AmazonS3) GetStoreType() utils.FileStoreType {
	return utils.AmazonS3
}

// Close gracefully closed the local filestore module
func (a *AmazonS3) Close() error {
	os.Remove(a.tempFileName)
	return nil
}
