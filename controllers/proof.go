package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"server-v2/config"
	"server-v2/models"
	"mime/multipart"
	"os"
	"path/filepath"
	"math/rand"
	"time"
	"io"
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

	// Simpan gambar ke direktori lokal dan perbarui URL gambar di model Proof
	if err := uploadProofImages(&proof, c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload images",
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
	// Direktori tempat menyimpan gambar
	imageDir := "uploads"

	// Buat direktori jika belum ada
	if err := os.MkdirAll(imageDir, os.ModePerm); err != nil {
		return err
	}

	// Upload gambar KTP
	photoKTP, err := c.FormFile("photo_ktp")
	if err != nil {
		return err
	}
	proof.PhotoKTPURL, err = uploadImage(imageDir, photoKTP)
	if err != nil {
		return err
	}

	// Upload gambar orang
	photoOrang, err := c.FormFile("photo_orang")
	if err != nil {
		return err
	}
	proof.PhotoOrangURL, err = uploadImage(imageDir, photoOrang)
	if err != nil {
		return err
	}

	// Upload gambar tangki
	photoTangki, err := c.FormFile("photo_tangki")
	if err != nil {
		return err
	}
	proof.PhotoTangkiURL, err = uploadImage(imageDir, photoTangki)
	if err != nil {
		return err
	}

	return nil
}

func uploadImage(imageDir string, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Generate nama file unik
	fileExt := filepath.Ext(fileHeader.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), fileExt)

	// Buat file baru di direktori tempat penyimpanan
	dstPath := filepath.Join(imageDir, fileName)
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dstFile.Close()

	// Salin data file ke file baru
	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}
	if _, err := io.Copy(dstFile, file); err != nil {
		return "", err
	}

	// Generate URL gambar dari path file lokal
	imageURL := fmt.Sprintf("/%s/%s", imageDir, fileName)

	return imageURL, nil
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
