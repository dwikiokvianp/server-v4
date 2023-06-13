package controllers

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
	"server-v2/models"
	"server-v2/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"net/http"
)

func uploadImageToS3(sess *session.Session, bucket string, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Generate unique file name
	fileExt := filepath.Ext(fileHeader.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), fileExt)

	// Initialize S3 uploader
	uploader := s3manager.NewUploader(sess)

	// Upload file to S3
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		return "", err
	}

	// Generate image URL from S3 bucket URL
	s3BucketURL := os.Getenv("S3_BUCKET_URL")
	fileURL := fmt.Sprintf("%s/%s", s3BucketURL, fileName)

	return fileURL, nil
}

func CreateProof(c *gin.Context) {
	fileKTP, err := c.FormFile("photo_ktp")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read photo_ktp from request",
		})
		return
	}

	fileOrang, err := c.FormFile("photo_orang")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read photo_orang from request",
		})
		return
	}

	fileTangki, err := c.FormFile("photo_tangki")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read photo_tangki from request",
		})
		return
	}

	description := c.PostForm("description")

	proof := &models.Proof{
		Description: description,
	}

	// Initialize AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to initialize AWS session",
		})
		return
	}

	// Upload KTP image to S3
	photoKTPURL, err := uploadImageToS3(sess, os.Getenv("S3_BUCKET_NAME"), fileKTP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload KTP file to S3",
		})
		return
	}
	proof.PhotoKTPURL = photoKTPURL

	// Upload orang image to S3
	photoOrangURL, err := uploadImageToS3(sess, os.Getenv("S3_BUCKET_NAME"), fileOrang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload orang file to S3",
		})
		return
	}
	proof.PhotoOrangURL = photoOrangURL

	// Upload tangki image to S3
	photoTangkiURL, err := uploadImageToS3(sess, os.Getenv("S3_BUCKET_NAME"), fileTangki)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload tangki file to S3",
		})
		return
	}
	proof.PhotoTangkiURL = photoTangkiURL

	// Save proof to database
	err = config.DB.Create(&proof).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save proof to database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Proof created successfully",
		"data":    proof,
	})
}

func GetAllProofs(c *gin.Context) {
	var proofs []models.Proof

	if err := config.DB.Find(&proofs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch proofs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": proofs,
	})
}

func GetProofByID(c *gin.Context) {
	id := c.Param("id")

	var proof models.Proof

	if err := config.DB.First(&proof, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Proof not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": proof,
	})
}
