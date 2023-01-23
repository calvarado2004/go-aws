package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"testing"
)

// MockS3Client is a mock S3 client
type MockS3Client struct {
	ListBucketOutput   *s3.ListBucketsOutput
	CreateBucketOutput *s3.CreateBucketOutput
}

// ListBuckets returns a list of buckets with two buckets already created
func (m MockS3Client) ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {

	ListBucketsOutput := &s3.ListBucketsOutput{
		Buckets: []types.Bucket{
			{
				Name: aws.String("testing-bucket-01"),
			},
			{
				Name: aws.String("testing-bucket-02"),
			},
		},
	}

	m.ListBucketOutput = ListBucketsOutput

	return m.ListBucketOutput, nil
}

// CreateBucket simulates creating a new bucket
func (m MockS3Client) CreateBucket(ctx context.Context, params *s3.CreateBucketInput, optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error) {
	return m.CreateBucketOutput, nil
}

// TestCreateS3Bucket tests the createS3Bucket function
func TestCreateS3Bucket(t *testing.T) {
	ctx := context.Background()
	err := createS3Bucket(ctx, MockS3Client{
		ListBucketOutput: &s3.ListBucketsOutput{},
	}, "testing-bucket-03")

	if err != nil {
		t.Fatalf("failed to create bucket: %v", err)
	}
}

// TestCreateS3BucketExistent tests if the bucket already exists
func TestCreateS3BucketExistent(t *testing.T) {
	ctx := context.Background()
	err := createS3Bucket(ctx, MockS3Client{
		ListBucketOutput: &s3.ListBucketsOutput{},
	}, "testing-bucket-02")

	if err != nil {
		t.Fatalf("failed to create bucket: %v", err)
	}
}
