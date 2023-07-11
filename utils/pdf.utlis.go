package utils

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
)

func UploadPdfToS3(pdf []byte, key string) (string, error) {
	// Retrieve AWS region from environment variable
	region := os.Getenv("AWS_REGION")
	if region == "" {
		return "", fmt.Errorf("AWS region not specified")
	}

	// Retrieve S3 bucket name from environment variable
	bucket := os.Getenv("S3_BUCKET_NAME")
	if bucket == "" {
		return "", fmt.Errorf("S3 bucket name not specified")
	}

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return "", err
	}

	// Create an S3 client
	svc := s3.New(sess)

	// Upload the PDF to S3
	_, err = svc.PutObject(&s3.PutObjectInput{
		Body:   bytes.NewReader(pdf),
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			return "", fmt.Errorf("S3 upload error: %v, %v", awsErr.Code(), awsErr.Message())
		}
		return "", err
	}

	// Generate the URL for the uploaded PDF
	pdfURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, key)

	return pdfURL, nil
}
