package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type AWS struct {
	REGIONS string `env:"AWS_S3_REGION,default=eu-central-1"`
	COGNITO struct {
		CLIENT_ID   string `env:"AWS_COGNITO_CLIENT_ID"`
		USERPOOL_ID string `env:"AWS_COGNITO_USERPOOL_ID"`
	}
	CLOUDFRONT struct {
		BASE_URL string `env:"AWS_CLOUDFRONT_BASE_URL"`
	}
}

func (a *AWS) GetAwsSession() *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Region: aws.String(a.REGIONS),
	}))
}
