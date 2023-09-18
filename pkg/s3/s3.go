package storage

import (
	"backend/pkg/config"

	"github.com/aws/aws-sdk-go/service/s3"
)

type Client interface {
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
	GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error)
	DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
}

type S3 struct {
	Connect *s3.S3
}

func NewS3Client() (Client, error) {
	cfg := config.InitConfig()

	sess := cfg.AWS.GetAwsSession()

	s3Client := s3.New(sess)
	return &S3{Connect: s3Client}, nil
}

func (c *S3) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	s3Client := c.Connect
	return s3Client.PutObject(input)
}

func (c *S3) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	s3Client := c.Connect
	return s3Client.GetObject(input)
}

func (c *S3) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	s3Client := c.Connect
	return s3Client.DeleteObject(input)
}
