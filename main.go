package main

import "log"

func main() {

	//ec2Instance()
	err := s3BucketMain()
	if err != nil {
		log.Printf("error creating S3 bucket: %v", err)
	}

}
