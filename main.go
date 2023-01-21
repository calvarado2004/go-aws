package main

import "log"

func main() {

	err := ec2Instance()
	if err != nil {
		log.Fatalf("error creating EC2 instance: %v", err)
	}

	err = s3BucketMain()
	if err != nil {
		log.Printf("error creating S3 bucket: %v", err)
	}

}
