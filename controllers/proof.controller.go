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
	"strings"
	// "io/ioutil"
	// "image"
	// _ "image/jpeg"
	// _ "image/png"

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
	invoicePDF, err := GenerateInvoicePDF(*proof, transaction, company)
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
		err = utils.SendEmail(email, subject, body, invoiceURL)
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

func GenerateInvoicePDF(proof models.Proof, transaction models.Transaction, company models.Company) ([]byte, error) {
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

	pdf.Cell(40, 10, "Diterbitkan:")
	pdf.Cell(0, 10, proof.CreatedAt.Format("2006-01-02"))
	pdf.Ln(8)

	pdf.Cell(40, 10, "Description:")
	pdf.MultiCell(0, 10, proof.Description, "", "", false)
	pdf.Ln(8)

	// // Embed photo_ktp image from S3 URL
	// if proof.PhotoKTPURL != "" {
	// 	err := embedImageFromURL(pdf, proof.PhotoKTPURL, pdf.GetX(), pdf.GetY()+10, 0, 30)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	pdf.Ln(40)
	// }

	// Display transaction details
	pdf.Cell(40, 10, "Transaction Status:")
	pdf.Cell(0, 10, transaction.Status)
	pdf.Ln(8)

	// Output the PDF as bytes
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// func embedImageFromURL(pdf *gofpdf.Fpdf, imageURL string, x, y, w, h float64) error {
// 	// Get the image from URL
// 	response, err := http.Get(imageURL)
// 	if err != nil {
// 		return err
// 	}
// 	defer response.Body.Close()

// 	// Read the image data
// 	imageData, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		return err
// 	}

// 	// Embed the image into the PDF
// 	pdf.RegisterImageOptionsReader("", gofpdf.ImageOptions{}, bytes.NewReader(imageData))

// 	// Add a page and display the image
// 	pdf.AddPage()
// 	pdf.ImageOptions("", x, y, w, h, false, gofpdf.ImageOptions{}, 0, "")

// 	return nil
// }

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
