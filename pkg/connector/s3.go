package connector

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// GetS3Connector создает новое подключение к S3.
func GetS3Connector(s3endpoint, region, accessKey, secretKey string) (*s3.S3, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	svc := s3.New(
		sess,
		aws.NewConfig().
			WithEndpoint(s3endpoint).
			WithRegion(region).
			WithCredentials(credentials.NewStaticCredentials(accessKey, secretKey, "")),
	)
	return svc, nil
}
