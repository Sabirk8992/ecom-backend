package storage

import (
	"context"
	"fmt"
	"mime/multipart"

	appconfig "github.com/Sabirk8992/ecom-backend/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client *s3.Client
	bucket string
	region string
}

func NewS3Storage(cfg *appconfig.Config) (*S3Storage, error) {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.AWSRegion),
		// No credentials needed — EC2 IAM role handles it
	)
	if err != nil {
		return nil, err
	}

	return &S3Storage{
		client: s3.NewFromConfig(awsCfg),
		bucket: cfg.S3Bucket,
		region: cfg.AWSRegion,
	}, nil
}

func (s *S3Storage) Upload(file multipart.File, fileName string, contentType string) (string, error) {
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(fileName),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, fileName)
	return url, nil
}
