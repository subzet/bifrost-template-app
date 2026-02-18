package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client  *s3.Client
	bucket  string
	baseURL string
}

// NewS3Storage creates an S3-compatible storage client (works with Backblaze B2).
// endpoint: e.g. "https://s3.us-west-004.backblazeb2.com"
// baseURL:  public base URL for uploaded files, e.g. "https://bucket.s3.us-west-004.backblazeb2.com"
func NewS3Storage(endpoint, region, bucket, keyID, appKey, baseURL string) (*S3Storage, error) {
	cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(keyID, appKey, "")),
		awsconfig.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	return &S3Storage{client: client, bucket: bucket, baseURL: baseURL}, nil
}

func (s *S3Storage) Upload(ctx context.Context, key string, r io.Reader, size int64, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          r,
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(size),
	})
	if err != nil {
		return "", fmt.Errorf("S3 upload failed: %w", err)
	}
	return s.baseURL + "/" + key, nil
}
