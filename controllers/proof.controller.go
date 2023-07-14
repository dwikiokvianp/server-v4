package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"server-v2/config"
	"server-v2/models"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"net/http"
	"server-v2/utils"
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

	// proof.SignatureURL = signatureURL

	signatureFile, err := c.FormFile("signature")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read signature from request",
		})
		return
	}

	signaturePath := "/save/signature.png" // Specify the path where the signature file should be saved

	if err := c.SaveUploadedFile(signatureFile, signaturePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save signature file",
		})
		return
	}

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

	var company models.Company
	if err := config.DB.First(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve company information",
		})
		return
	}

	transaction.StatusId = 6

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
		historyOut := models.HistoryOut{
			Date:          time.Now(),
			UserId:        transaction.UserId,
			TransactionId: transactionIdInt,
			Quantity:      int(item.Quantity),
			OilId:         int(item.OilID),
		}

		if err := config.DB.Create(&historyOut).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save history to database",
			})
			return
		}
	}

	// Generate the invoice PDF with the signature
	invoicePDF, err := GenerateInvoicePDF(*proof, transaction, company, signaturePath)
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

	email := transaction.Email
	subject := "Bukti Transaction"
	body := `
	<html>
	<head>
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #f2f2f2;
				padding: 20px;
			}
	
			h1 {
				color: #333333;
				font-size: 24px;
				font-weight: bold;
				margin-bottom: 20px;
			}
	
			p {
				color: #666666;
				font-size: 16px;
				line-height: 1.5;
				margin-bottom: 10px;
			}
	
			.qr-code {
				display: block;
				text-align: center;
				margin-bottom: 20px;
			}
	
			.qr-code img {
				max-width: 200px;
				height: auto;
			}
		</style>
	</head>
	<body>
		<p>Bukti Transactions</p>
	</body>
	</html>
	`

	go func() {
		err = utils.SendPDFEmail(email, subject, body, invoicePDF)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
	}()

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
	})
}

func GenerateInvoicePDF(proof models.Proof, transaction models.Transaction, company models.Company, signaturePath string) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle("Invoice", true)
	pdf.SetAuthor("Your Company", true)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 10, "Invoice")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Company:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, company.CompanyName)
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Address:")

	cleanAddress := strings.Replace(company.Address, "%68", "", -1)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, cleanAddress)
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Email:")
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, transaction.Email)
	pdf.Ln(12)

	var status models.Status
    if err := config.DB.First(&status, transaction.StatusId).Error; err != nil {
	// Penanganan jika terjadi kesalahan
    }
	
	// Menggunakan nama status pada PDF
	pdf.Cell(40, 10, "Transaction Status:")
	pdf.Cell(0, 10, status.Name)
	pdf.Ln(8)

	pageWidth, pageHeight := pdf.GetPageSize()

	// Embed signature image from local file
	if signaturePath != "" {
		signatureWidth := 60.0
		signatureHeight := 30.0
		signatureX := (pageWidth - signatureWidth) / 2
		signatureY := pageHeight - signatureHeight - 20

		err := embedImageFromFile(pdf, signaturePath, signatureX, signatureY, signatureWidth, signatureHeight)
		if err != nil {
			return nil, err
		}
	}

	// Output the PDF as bytes
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// embedImageFromFile embeds an image from a local file in the PDF at the specified position and dimensions.
func embedImageFromFile(pdf *gofpdf.Fpdf, imagePath string, x, y, width, height float64) error {
	imageData, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return err
	}

	// Check the image type based on the file extension
	imageType := filepath.Ext(imagePath)[1:]
	if imageType == "" {
		return errors.New("unknown image type")
	}

	// Embed the image in the PDF
	pdf.RegisterImageOptionsReader(imagePath, gofpdf.ImageOptions{ImageType: imageType}, bytes.NewReader(imageData))
	pdf.ImageOptions(imagePath, x, y, width, height, false, gofpdf.ImageOptions{}, 0, "")

	return nil
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
