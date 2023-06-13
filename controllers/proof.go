package controllers

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/gin-gonic/gin"
	"net/http"
	"server-v2/config"
	"server-v2/models"
	"mime/multipart"
	"os"
	"path/filepath"
	"math/rand"
	"time"
)

func CreateProof(c *gin.Context) {
	var proof models.Proof

	// Bind form data ke struct Proof
	if err := c.ShouldBind(&proof); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Simpan gambar ke S3 dan perbarui URL gambar di model Proof
	if err := uploadProofImages(&proof, c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload images to S3",
		})
		return
	} else {
		// Perbarui URL gambar di model Proof
		if err := config.DB.Model(&proof).Updates(models.Proof{
			PhotoKTPURL:    proof.PhotoKTPURL,
			PhotoOrangURL:  proof.PhotoOrangURL,
			PhotoTangkiURL: proof.PhotoTangkiURL,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update image URLs",
			})
			return
		}
	}

	// Simpan model Proof ke database
	if err := config.DB.Create(&proof).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create proof",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Proof created successfully",
		"data":    proof,
	})
}

func uploadProofImages(proof *models.Proof, c *gin.Context) error {
	// Membaca konfigurasi S3 dari environment variables
	region := os.Getenv("AWS_REGION")
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	// Inisialisasi session AWS
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return err
	}

	// Membaca URL bucket S3 dari environment variable
	s3BucketURL := os.Getenv("S3_BUCKET_URL")

	// Upload gambar KTP
	photoKTP, err := c.FormFile("photo_ktp")
	if err != nil {
		return err
	}
	proof.PhotoKTPURL, err = uploadImageToS3(sess, s3BucketURL, photoKTP)
	if err != nil {
		return err
	}

	// Upload gambar orang
	photoOrang, err := c.FormFile("photo_orang")
	if err != nil {
		return err
	}
	proof.PhotoOrangURL, err = uploadImageToS3(sess, s3BucketURL, photoOrang)
	if err != nil {
		return err
	}

	// Upload gambar tangki
	photoTangki, err := c.FormFile("photo_tangki")
	if err != nil {
		return err
	}
	proof.PhotoTangkiURL, err = uploadImageToS3(sess, s3BucketURL, photoTangki)
	if err != nil {
		return err
	}

	return nil
}

func uploadImageToS3(sess *session.Session, bucket string, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Generate nama file unik
	fileExt := filepath.Ext(fileHeader.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), fileExt)

	// Upload file ke S3
	svc := s3.New(sess)
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		return "", err
	}

	// Generate URL gambar dari S3 bucket URL
	return fmt.Sprintf("%s/%s", bucket, fileName), nil
}

var randomGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateRandomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[randomGenerator.Intn(len(letterRunes))]
	}
	return string(b)
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
