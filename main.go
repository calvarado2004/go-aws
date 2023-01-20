package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"log"
	"os"
)

func main() {

	var (
		instanceId string
		err        error
	)
	ctx := context.Background()

	instanceId, err = createEC2Instance(ctx, "us-east-1")
	if err != nil {
		fmt.Printf("error creating EC2 instance: %v", err)
		os.Exit(1)
	}

	fmt.Printf("New EC2 instance created. Instance ID: %s", instanceId)

}

func createEC2Instance(ctx context.Context, region string) (instanceCreated string, err error) {

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
		return "", err
	}

	// Create an EC2 client from just a session.
	ec2client := ec2.NewFromConfig(cfg)

	// Check if the key pair already exists
	keyPair, err := ec2client.DescribeKeyPairs(ctx, &ec2.DescribeKeyPairsInput{
		KeyNames: []string{"go-aws-key"},
	})

	idExistent := *keyPair.KeyPairs[0].KeyName

	if idExistent == "" {
		_, err = ec2client.CreateKeyPair(ctx, &ec2.CreateKeyPairInput{
			KeyName: aws.String("go-aws-key"),
		})
		if err != nil {
			log.Printf("Unable to create key pair, %v", err)
			return "", err
		}
	}

	// Create an EC2 instance with the key pair
	imageOutput, err := ec2client.DescribeImages(ctx, &ec2.DescribeImagesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("name"),
				Values: []string{"ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"},
			},
		},
		Owners: []string{"099720109477"},
	})
	if err != nil {
		log.Fatalf("Unable to describe images, %v", err)
		return "", err
	}

	if len(imageOutput.Images) == 0 {
		log.Fatalf("Unable to find image, %v", err)
		return "", err
	}

	//imageOutput.Images[0].ImageId

	instance, err := ec2client.RunInstances(ctx, &ec2.RunInstancesInput{
		ImageId:      imageOutput.Images[0].ImageId,
		KeyName:      aws.String("go-aws-key"),
		InstanceType: types.InstanceTypeT3Micro,
		MaxCount:     aws.Int32(1),
		MinCount:     aws.Int32(1),
	})
	if err != nil {
		log.Fatalf("Unable to run instance, %v", err)
		return "", err
	}

	return *instance.Instances[0].InstanceId, err

}
