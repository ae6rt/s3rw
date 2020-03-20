package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	bucketName  = flag.String("bucket-name", "", "Bucket name")
	objectKey   = flag.String("object-key", "", "Item key in bucket")
	fileName    = flag.String("filename", "", "file to get or put")
	profile     = flag.String("profile", "", "AWS profile from $HOME/.aws/credentials to use")
	region      = flag.String("region", "", "AWS region")
	op          = flag.String("op", "", "get or put")
	publicRead  = flag.Bool("make-public", false, "Make object publicly readable")
	contentType = flag.String("content-type", "", "MIME type")
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

	sess, err := session.NewSession(config)
	if err != nil {
		log.Fatal(err)
	}

	switch *op {
	case "get":
		file, err := os.Create(*fileName)
		if err != nil {
			log.Fatal(err)
		}

		downloader := s3manager.NewDownloader(sess)

		_, err = downloader.Download(file,
			&s3.GetObjectInput{
				Bucket: aws.String(*bucketName),
				Key:    aws.String(*objectKey),
			})
		if err != nil {
			log.Fatal(err)
		}
	case "put":
		data, err := ioutil.ReadFile(*fileName)
		if err != nil {
			log.Fatal(err)
		}

		uploader := s3manager.NewUploader(sess)

		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket:               aws.String(*bucketName),
			Key:                  aws.String(*objectKey),
			Body:                 bytes.NewReader(data),
			ServerSideEncryption: aws.String("AES256"),
		})
	default:
		log.Fatalf("Operation not supported %s\n", *op)
	}
}
