package controllers

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"server-v2/config"
	"server-v2/models"
	"strconv"
	"time"

	"server-v2/utils"
	"github.com/jung-kurt/gofpdf"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/jaytaylor/html2text"
	"net/http"
)

func uploadImageToS3(sess *session.Session, bucket string, fileHeader *multipart.FileHeader) (string, error) {
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
	// signature := c.PostForm("signature")

	transactionId := c.Param("id")

	transactionIdInt, err := strconv.Atoi(transactionId)

	proof := &models.Proof{
		TransactionID: transactionIdInt,
		Description:   description,
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to initialize AWS session",
		})
		return
	}

	photoKTPURL, err := uploadImageToS3(sess, os.Getenv("S3_BUCKET_NAME"), fileKTP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload KTP file to S3",
		})
		return
	}
	proof.PhotoKTPURL = photoKTPURL

	photoOrangURL, err := uploadImageToS3(sess, os.Getenv("S3_BUCKET_NAME"), fileOrang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload orang file to S3",
		})
		return
	}
	proof.PhotoOrangURL = photoOrangURL

	photoTangkiURL, err := uploadImageToS3(sess, os.Getenv("S3_BUCKET_NAME"), fileTangki)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload tangki file to S3",
		})
		return
	}
	proof.PhotoTangkiURL = photoTangkiURL

	err = config.DB.Create(&proof).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save proof to database",
		})
		return
	}

	var transaction models.Transaction
	if err := config.DB.
		Preload("TransactionDetail").
		First(&transaction, transactionIdInt).
		Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Transaction not found",
		})
		return
	}

	transaction.Status = "done"

	if err := config.DB.Save(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update transaction status",
		})
		return
	}

	transactionIdInt, err = strconv.Atoi(transactionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to convert transaction id to int",
		})
		return
	}

	var historyOut models.HistoryOut
	fmt.Println(transactionIdInt)

	historyOut.Date = time.Now()
	historyOut.UserId = transaction.UserId
	historyOut.TransactionId = transactionIdInt

	fmt.Println(transaction.TransactionDetail)

	for _, item := range transaction.TransactionDetail {
		historyOut.Quantity = int(item.Quantity)
		historyOut.OilId = int(item.OilID)

		if err := config.DB.Create(&historyOut).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save history to database",
			})
			return
		}
	}

	// Generate the invoice PDF
	invoicePDF, err := GenerateInvoicePDF(*proof, transaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate invoice PDF",
		})
		return
	}

	// Upload the PDF to S3
	invoiceKey := fmt.Sprintf("invoices/%v.pdf", proof.ID)
	invoiceURL, err := utils.UploadPdfToS3(invoicePDF, invoiceKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload invoice PDF to S3",
		})
		return
	}

	// Save the invoice URL in the proof
	proof.InvoiceURL = invoiceURL
	err = config.DB.Save(&proof).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save invoice URL in the proof",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Proof created successfully and transaction status updated to done",
		"data":    proof,
	})
}

func GenerateInvoicePDF(proof models.Proof, transaction models.Transaction) ([]byte, error) {
	// Convert description from HTML to plain text
	descriptionText, err := html2text.FromString(proof.Description)
	if err != nil {
		return nil, err
	}

	// Create a new PDF object
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle("Invoice", true)
	pdf.SetAuthor("Your Company", true)
	pdf.AddPage()

	// Set font and size for the title
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 10, "Invoice")
	pdf.Ln(12)

	// Set font and size for the content
	pdf.SetFont("Arial", "", 12)

	// Display proof details
	pdf.Cell(40, 10, "Transaction ID:")
	pdf.Cell(0, 10, strconv.Itoa(proof.TransactionID))
	pdf.Ln(8)

	pdf.Cell(40, 10, "Description:")
	pdf.MultiCell(0, 10, descriptionText, "", "", false)
	pdf.Ln(8)

	// Display photo URLs
	pdf.Cell(40, 10, "Photo KTP URL:")
	pdf.Cell(0, 10, proof.PhotoKTPURL)
	pdf.Ln(8)

	pdf.Cell(40, 10, "Photo Orang URL:")
	pdf.Cell(0, 10, proof.PhotoOrangURL)
	pdf.Ln(8)

	pdf.Cell(40, 10, "Photo Tangki URL:")
	pdf.Cell(0, 10, proof.PhotoTangkiURL)
	pdf.Ln(8)

	// Display transaction details
	pdf.Cell(40, 10, "Transaction Status:")
	pdf.Cell(0, 10, transaction.Status)
	pdf.Ln(8)

	pdf.Cell(40, 10, "Transaction Date:")
	pdf.Cell(0, 10, transaction.Date.Format("2006-01-02"))
	pdf.Ln(8)

	// Output the PDF as bytes
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
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

func GetProofByTransactionId(c *gin.Context) {
	var proofs models.Proof

	transactionId := c.Param("id")

	if err := config.DB.Where("transaction_id = ?", transactionId).First(&proofs).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Proof not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": proofs,
	})
}
