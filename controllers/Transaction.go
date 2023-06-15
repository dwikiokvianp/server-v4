package controllers

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/skip2/go-qrcode"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"time"
	"server-v2/config"
	"server-v2/models"
	"strconv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GenerateQRCode(data string) (string, error) {
	qrFile := "qrcode.png"
	err := qrcode.WriteFile(data, qrcode.Medium, 256, qrFile)
	if err != nil {
		return "", err
	}
	return qrFile, nil
}
func UploadToS3(file string, key string) (string, error) {
	region := os.Getenv("AWS_REGION")
	bucketName := os.Getenv("S3_BUCKET_NAME")
	bucketURL := os.Getenv("S3_BUCKET_URL")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return "", err
	}

	svc := s3.New(sess)

	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   f,
	})
	if err != nil {
		return "", err
	}

	url := bucketURL + key
	return url, nil
}

func SendEmail(to, subject, body, attachment string) error {
	from := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if attachment != "" {
		m.Attach(attachment)
	}

	d := gomail.NewDialer(smtpHost, smtpPort, from, password)

	err := d.DialAndSend(m)
	if err != nil {
		return err
	}
	return nil
}

func CreateTransactions(c *gin.Context) {

	var inputTransaction models.TransactionInput

	if err := c.ShouldBindJSON(&inputTransaction); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	userId := c.Param("id")

	intUserId, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	transaction := models.Transaction{
		UserId:    intUserId,
		Email:     inputTransaction.Email,
		VehicleId: inputTransaction.VehicleId,
		OilId:     inputTransaction.OilId,
		Quantity:  inputTransaction.Quantity,
		OfficerId: inputTransaction.OfficerId,
		Status:    "pending",
	}

	if err := config.DB.Create(&transaction).Error; err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Generate QR code
	qrData := strconv.Itoa(int(transaction.ID))
	qrFile, err := GenerateQRCode(qrData)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Upload QR code to S3
	key := "qrcodes/" + qrFile
	qrURL, err := UploadToS3(qrFile, key)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	email := inputTransaction.Email
	subject := "QR Code Transaction"
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
		<p>Tunjukkan QR kode ini kepada petugas untuk mendapatkan layanan.</p>
	</body>
	</html>
	`
	err = SendEmail(email, subject, body, qrFile)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Update transaction with QR code URL
	transaction.QrCodeUrl = qrURL
	if err := config.DB.Save(&transaction).Error; err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success create transaction",
	})
}

func GetAllTransactions(c *gin.Context) {
	var (
		transactions []models.Transaction
		pageSize     = 5
		page         = 1
	)

	pageParam := c.Query("page")
	if pageParam != "" {
		page, _ = strconv.Atoi(pageParam)
	}

	db := config.DB

	var count int64
	if err := db.Model(&models.Transaction{}).Count(&count).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	offset := (page - 1) * pageSize

	err := db.Offset(offset).Limit(pageSize).
		Preload("Vehicle.VehicleType").
		Preload("User.Role").
		Preload("Officer").
		Preload("User.Detail").
		Preload("User.Company").
		Preload("Oil").
		Find(&transactions).Error
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	totalPages := (int(count) + pageSize - 1) / pageSize

	c.JSON(200, gin.H{
		"data":     transactions,
		"page":     page,
		"pageSize": pageSize,
		"total":    totalPages,
	})
}

func GetByIdTransaction(c *gin.Context) {
	var transaction models.Transaction
	id := c.Param("id")
	err := config.DB.Preload("Vehicle.VehicleType").Preload("User.Role").Preload("Officer").Preload("User.Detail").Preload("User.Company").Preload("Oil").Find(&transaction, id).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(200, gin.H{
		"data": transaction,
	})
}

func GetTransactionByUserId(c *gin.Context) {
	var transaction []models.Transaction
	id := c.Param("id")
	fmt.Println(id)
	err := config.DB.Preload("Vehicle.VehicleType").Preload("Oil").Preload("User.Detail").Where("user_id = ?", id).Find(&transaction).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}

	if len(transaction) == 0 {
		c.JSON(404, gin.H{
			"message": "Not Found",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": transaction,
	})
}

func GetTodayTransactions() ([]models.Transaction, error) {
	var transactions []models.Transaction
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfToday := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	if err := config.DB.Where("created_at BETWEEN ? AND ?", startOfToday.Unix(), endOfToday.Unix()).Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

