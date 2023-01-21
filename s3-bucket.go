package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	uuid "github.com/nu7hatch/gouuid"
	"log"
)

// s3BucketMain creates an S3 bucket with an uuid name appended to the bucket name
func s3BucketMain() error {

	ctx := context.Background()
	s3Client, err := initS3Client(ctx, "us-east-1")
	if err != nil {
		log.Printf("error initializing %s", err)
	}

	u, err := uuid.NewV4()
	s3Name := "go-aws-s3-" + u.String()

	err = createS3Bucket(ctx, s3Client, s3Name)
	if err != nil {
		log.Printf("error creating S3 bucket: %v", err)
		return err
	}

	return nil
}

// initS3Client initializes an S3 client
func initS3Client(context context.Context, region string) (*s3.Client, error) {

	cfg, err := config.LoadDefaultConfig(context, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}

// createS3Bucket creates an S3 bucket
func createS3Bucket(context context.Context, s3Client *s3.Client, bucketName string) error {

	_, err := s3Client.CreateBucket(context, &s3.CreateBucketInput{
		Bucket: &bucketName,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraintUsWest1,
		},
	})
	if err != nil {
		return err
	}

	log.Printf("Bucket created: %s", bucketName)

	return nil
}
