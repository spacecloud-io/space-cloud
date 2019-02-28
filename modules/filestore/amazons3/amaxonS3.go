package amazons3

import (
	"github.com/spaceuptech/space-cloud/utils"
)

type AmazonS3 struct {
	region string
}

func Init(region string) (*AmazonS3, error) {
	return &AmazonS3{region}, nil
}

// GetStoreType returns the file store type
func (a *AmazonS3) GetStoreType() utils.FileStoreType {
	return utils.AmazonS3
}

// Close gracefully closed the local filestore module
func (a *AmazonS3) Close() error {
	return nil
}
