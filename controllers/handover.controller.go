package controllers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"server-v2/config"
	"server-v2/models"
	"strconv"
	"time"
	"github.com/dranikpg/dto-mapper"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
)

func uploadsToS3(sess *session.Session, bucket string, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	fileExt := filepath.Ext(fileHeader.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), fileExt)

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		return "", err
	}

	s3BucketURL := os.Getenv("S3_BUCKET_URL")
	fileURL := fmt.Sprintf("%s%s", s3BucketURL, fileName)

	return fileURL, nil
}

func CreateHandover(c *gin.Context) {
	var handover models.Handover

	// Bind data from form-data
	handover.WorkerBeforeId, _ = strconv.Atoi(c.PostForm("worker_before_id"))
	handover.WorkerAfterId, _ = strconv.Atoi(c.PostForm("worker_after_id"))
	handover.Condition = c.PostForm("condition")
	handover.Status = "pending" 

	id := c.MustGet("id").(string)
	idInt, _ := strconv.Atoi(id)
	handover.EmployeeId = idInt

	fileTangki, err := c.FormFile("handover_tangki")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read handover_tangki from request",
		})
		return
	}

	fileKebersihan, err := c.FormFile("handover_kebersihan")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read handover_kebersihan from request",
		})
		return
	}

	fileLevelGauge, err := c.FormFile("handover_level_gauge")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read handover_level_gauge from request",
		})
		return
	}

	filePetugas, err := c.FormFile("handover_petugas")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read handover_petugas from request",
		})
		return
	}

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to initialize AWS session",
		})
		return
	}

	handoverTangkiURL, err := uploadsToS3(sess, os.Getenv("S3_BUCKET_NAME"), fileTangki)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload handover_tangki to S3",
		})
		return
	}
	handover.HandoverTangki = handoverTangkiURL

	handoverKebersihanURL, err := uploadsToS3(sess, os.Getenv("S3_BUCKET_NAME"), fileKebersihan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload handover_kebersihan to S3",
		})
		return
	}
	handover.HandoverKebersihan = handoverKebersihanURL

	handoverLevelGaugeURL, err := uploadsToS3(sess, os.Getenv("S3_BUCKET_NAME"), fileLevelGauge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload handover_level_gauge to S3",
		})
		return
	}
	handover.HandoverLevelGauge = handoverLevelGaugeURL

	handoverPetugasURL, err := uploadsToS3(sess, os.Getenv("S3_BUCKET_NAME"), filePetugas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload handover_petugas to S3",
		})
		return
	}
	handover.HandoverPetugas = handoverPetugasURL

	// Save the handover record to the database
	if err := config.DB.Create(&handover).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create handover",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": "Successfully created handover",
	})
}

func GetHandoverById(c *gin.Context) {
	id := c.MustGet("id").(string)
	idInt, _ := strconv.Atoi(id)

	var handover models.Handover
	var handoverResponse models.HandoverResponse

	if err := config.DB.
		Where("employee_id = ?", idInt).
		Preload("WorkerBefore").
		Preload("WorkerAfter").
		Joins("Employee").
		Order("id desc").
		First(&handover).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Error here in the get handover by id",
		})
		return
	}

	err := dto.Map(&handoverResponse, &handover)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Error here in the get handover by id",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": handoverResponse,
	})
}
