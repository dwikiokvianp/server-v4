package controllers

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"gopkg.in/gomail.v2"
	"os"
	"server-v2/config"
	"server-v2/models"
	"strconv"
	"time"
)

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
	var transactionToday models.Transaction
	if err := config.DB.Where("user_id = ? AND date = ?", userId, time.Now().Format("2006-01-02")).First(&transactionToday).Error; err != nil {
		fmt.Println("record not found")
	} else {
		c.JSON(400, gin.H{
			"message": "You have already made a transaction today",
		})
		return
	}

	intUserId, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	transaction := models.Transaction{
		UserId:     intUserId,
		Email:      inputTransaction.Email,
		VehicleId:  inputTransaction.VehicleId,
		OilId:      inputTransaction.OilId,
		Quantity:   inputTransaction.Quantity,
		OfficerId:  inputTransaction.OfficerId,
		Status:     "pending",
		Date:       inputTransaction.Date,
		CityId:     inputTransaction.CityId,
		ProvinceId: inputTransaction.ProvinceId,
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
	go func() {
		err = SendEmail(email, subject, body, qrFile)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
	}()

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
		Preload("User").
		Preload("Oil").
		Preload("City").
		Preload("Province").
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
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfToday := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	var transactions []models.Transaction
	err := config.DB.
		Preload("User").
		Where("date BETWEEN ? AND ?", startOfToday, endOfToday).
		Find(&transactions).Error

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func GetTomorrowTransactions() ([]models.Transaction, error) {
	now := time.Now()
	startOfTomorrow := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(24 * time.Hour)
	endOfTomorrow := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location()).Add(24 * time.Hour)

	var transactions []models.Transaction
	err := config.DB.
		Preload("User").
		Where("date BETWEEN ? AND ?", startOfTomorrow.Format("2006-01-02T15:04:05Z07:00"), endOfTomorrow.Format("2006-01-02T15:04:05Z07:00")).
		Find(&transactions).Error

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func GetTodayV2Transaction(c *gin.Context) {
	var transaction []models.Transaction

	today := time.Now().Format("2006-01-02")

	name := c.Query("name")

	if name != "" {
		err := config.DB.
			Preload("Vehicle.VehicleType").
			Preload("User.Company").Preload("Oil").
			Preload("User.Detail").
			Joins("JOIN users ON users.id = transactions.user_id").
			Where("users.username = ?", name).
			Where("date = ?", today).Find(&transaction).Error
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		}
	} else {
		err := config.DB.Preload("Vehicle.VehicleType").Preload("User.Company").Preload("Oil").
			Preload("User.Detail").
			Where("date = ?", today).Find(&transaction).Error
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		}
	}

	if len(transaction) == 0 {
		c.JSON(404, gin.H{
			"message": "Transaction Today Not Found",
		})
	}

	c.JSON(200, gin.H{
		"data": transaction,
	})

}
