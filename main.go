package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	bucketName  = flag.String("bucket-name", "", "Bucket name")
	objectKey   = flag.String("object-key", "", "Item key in bucket")
	fileName    = flag.String("filename", "", "file to get or put")
	profile     = flag.String("profile", "", "AWS profile from $HOME/.aws/credentials to use")
	region      = flag.String("region", "", "AWS region")
	op          = flag.String("op", "", "get or put")
	versionFlag = flag.Bool("version", false, "Print version info and exit.")

	buildInfo string
)

func init() {
	flag.Parse()
	if *versionFlag {
		log.Printf("%s\n", buildInfo)
		os.Exit(0)
	}
}

func main() {
	home := os.Getenv("HOME")
	if home == "" {
		log.Fatalf("HOME environment variable not set")
	}
	credentialsFile := fmt.Sprintf("%s/.aws/credentials", home)
	config := aws.NewConfig().WithCredentials(credentials.NewSharedCredentials(credentialsFile, *profile)).WithRegion(*region).WithMaxRetries(3)

	svc := s3.New(config)

	switch *op {
	case "get":
		params := &s3.GetObjectInput{
			Bucket: aws.String(*bucketName),
			Key:    aws.String(*objectKey),
		}

		resp, err := svc.GetObject(params)
		if err != nil {
			log.Println(err.Error())
			log.Fatal(err)
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		err = ioutil.WriteFile(*fileName, data, 0644)
		if err != nil {
			log.Fatal(err)
		}
	case "put":
		data, err := ioutil.ReadFile(*fileName)
		if err != nil {
			log.Fatal(err)
		}
		params := &s3.PutObjectInput{
			Key:                  aws.String(*objectKey),
			Bucket:               aws.String(*bucketName),
			Body:                 bytes.NewReader(data),
			ServerSideEncryption: aws.String("AES256"),
		}

		_, err = svc.PutObject(params)

		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				log.Println(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
				if reqErr, ok := err.(awserr.RequestFailure); ok {
					log.Fatal(reqErr.Code(), reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())
				}
			} else {
				log.Fatal(err.Error())
			}
		}

	default:
		log.Fatalf("Operation not supported %s\n", *op)
	}
}
