package s3

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kondohiroki/go-grpc-boilerplate/config"
)

type S3Service struct {
	awsRegion          string
	awsAccessKeyID     string
	awsSecretAccessKey string
	bucket             string
	path               string
}

func NewS3Service(conf *config.Config) S3Service {
	return S3Service{
		awsRegion:          conf.Services.S3.AwsRegion,
		awsAccessKeyID:     conf.Services.S3.AwsAccessKeyID,
		awsSecretAccessKey: conf.Services.S3.AwsSecretAccessKey,
		bucket:             conf.Services.S3.Bucket,
		path:               conf.Services.S3.Path,
	}

}

func (s *S3Service) PutObject(file *os.File, fileName string) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(s.awsRegion),
		Credentials: credentials.NewStaticCredentials(s.awsAccessKeyID, s.awsSecretAccessKey, ""),
	})
	if err != nil {
		return err
	}

	// Create an S3 service client
	s3Client := s3.New(sess)
	key := fmt.Sprintf("%s/%s", s.path, fileName)

	// fileInfo, _ := file.Stat()
	// fileSize := fileInfo.Size()
	// Create an S3 upload input
	input := &s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	// Perform the file upload
	_, err = s3Client.PutObject(input)
	return err
}
