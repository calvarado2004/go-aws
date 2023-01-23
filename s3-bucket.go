package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"log"
	"os"
	"strings"
)

// s3BucketMain creates an S3 bucket with an uuid name appended to the bucket name
func s3BucketMain() error {

	ctx := context.Background()
	s3Client, err := initS3Client(ctx, "us-east-1")
	if err != nil {
		log.Printf("error initializing %s", err)
	}

	//u, err := uuid.NewV4()
	//s3Name := "go-aws-s3-" + u.String()

	s3Name := "go-aws-s3-bucket"

	err = createS3Bucket(ctx, s3Client, s3Name)
	if err != nil {
		log.Printf("error creating S3 bucket: %v", err)
		return err
	}

	err = uploadToS3Bucket(ctx, s3Client, s3Name, "test.txt", "./test.txt")
	if err != nil {
		log.Printf("error uploading to S3 bucket: %v", err)
		return err
	}

	fileReceived, err := downloadFromS3Bucket(ctx, s3Client, s3Name, "test.txt")
	if err != nil {
		log.Printf("error downloading from S3 bucket: %v", err)
		return err
	}

	err = os.WriteFile("./test-received.txt", fileReceived, 0644)
	if err != nil {
		log.Printf("error writing file: %v", err)
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

	allBuckets, err := s3Client.ListBuckets(context, &s3.ListBucketsInput{})
	if err != nil {
		return err
	}

	found := false

	for _, bucket := range allBuckets.Buckets {

		if *bucket.Name == bucketName {
			found = true
			break
		}
	}

	if !found {
		_, err := s3Client.CreateBucket(context, &s3.CreateBucketInput{
			Bucket: &bucketName,
		})
		if err != nil {
			return err
		}

		log.Printf("Bucket created: %s", bucketName)
	} else {
		log.Printf("Bucket already exists: %s", bucketName)
	}

	return nil
}

// uploadToS3Bucket uploads a file to an S3 bucket
func uploadToS3Bucket(ctx context.Context, s3Client *s3.Client, bucketName string, fileName string, filePath string) error {

	fileToSent, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("error reading file: %v", err)
		return err
	}

	uploader := manager.NewUploader(s3Client)

	uploadResult, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &fileName,
		Body:   io.NopCloser(strings.NewReader(string(fileToSent))),
	})

	if err != nil {
		log.Printf("error uploading file to S3 bucket: %v", err)
		return err
	}

	log.Printf("file uploaded to S3 bucket: %s", uploadResult.Location)

	return nil
}

// downloadFromS3Bucket downloads a file from an S3 bucket
func downloadFromS3Bucket(ctx context.Context, s3Client *s3.Client, bucketName string, fileName string) ([]byte, error) {

	downloader := manager.NewDownloader(s3Client)

	buffer := manager.NewWriteAtBuffer([]byte{})

	downloadResult, err := downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &fileName,
	})

	if err != nil {
		log.Printf("error downloading file from S3 bucket: %v", err)
		return nil, err
	}

	bytesReceived := int64(len(buffer.Bytes()))

	if bytesReceived != downloadResult {
		log.Printf("number of bytes received %d does not match the number of bytes downloaded %d", bytesReceived, downloadResult)
		return nil, nil
	}

	return buffer.Bytes(), nil

}
