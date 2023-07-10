package utils

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// UploadPDFToS3 mengunggah file PDF ke layanan S3
func UploadPDFToS3(file *multipart.FileHeader, key string) (string, error) {
	// Ambil konfigurasi AWS dari environment variables
	region := os.Getenv("AWS_REGION")    
	bucketName := os.Getenv("S3_BUCKET_NAME")   
	endpointURL := os.Getenv("S3_BUCKET_URL") 

	// Buat session AWS
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return "", err
	}

	// Buat layanan S3
	svc := s3.New(sess)

	// Buka file PDF
	pdfFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer pdfFile.Close()

	// Baca file PDF
	buffer := make([]byte, file.Size)
	_, err = pdfFile.Read(buffer)
	if err != nil {
		return "", err
	}

	// Tentukan metadata file
	contentType := "application/pdf"

	// Buat objek PutObjectInput untuk mengunggah file ke S3
	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		ACL:         aws.String("public-read"), // Opsi ini memberikan akses publik ke file yang diunggah, sesuaikan dengan kebutuhan keamanan Anda
		Body:        bytes.NewReader(buffer),
		ContentType: aws.String(contentType),
	}

	// Unggah file ke S3
	_, err = svc.PutObject(input)
	if err != nil {
		return "", err
	}

	// Dapatkan URL publik untuk file yang diunggah
	url := fmt.Sprintf("%s/%s", endpointURL, key)

	return url, nil
}
