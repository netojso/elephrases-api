package storage

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	configproject "github.com/netojso/elephrases-api/config"
	portrepository "github.com/netojso/elephrases-api/internal/core/ports/repository"
)

type S3Adapter struct {
	client *s3.Client
	bucket string
}

func NewS3Adapter(env *configproject.Env) (*S3Adapter, error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(env.AwsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			env.AwsAccessKeyID,
			env.AwsSecretAccessKey,
			"",
		)),
	)

	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &S3Adapter{
		client: client,
		bucket: env.AwsBucketName,
	}, nil
}

// Upload uploads a file to S3
func (a *S3Adapter) Upload(file portrepository.File) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(file.Name),
		Body:   bytes.NewReader(file.Data),
	}

	_, err := a.client.PutObject(context.TODO(), input)
	return err
}
