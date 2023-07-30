package controllers

import (
	"fmt"
	"github.com/dranikpg/dto-mapper"
	"github.com/gin-gonic/gin"
	"net/http"
	"server-v2/config"
	"server-v2/models"
	"server-v2/utils"
	"strconv"
	"time"
)

func CreateTransactions(c *gin.Context) {

	var inputTransaction models.TransactionInput

	if err := c.ShouldBindJSON(&inputTransaction); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	userId := c.Param("id")

	user := models.User{}
	if err := config.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "User not found",
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

	fmt.Println("user id", intUserId)

	transaction := models.Transaction{
		CustomerId: intUserId,
		Email:      user.Email,
		OfficerId:  inputTransaction.OfficerId,
		Date:       inputTransaction.Date,
		CityId:     inputTransaction.CityId,
		ProvinceId: inputTransaction.ProvinceId,
		StatusId:   inputTransaction.StatusId,
	}

	fmt.Println("status", inputTransaction.StatusId)

	if err := config.DB.Create(&transaction).Error; err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	var transactionDetails []models.TransactionDetail

	if inputTransaction.TransactionDetail == nil {
		c.JSON(400, gin.H{
			"message": "transaction detail is required",
		})
		return
	}

	for _, detail := range inputTransaction.TransactionDetail {
		transactionDetail := models.TransactionDetail{
			Quantity:      detail.Quantity,
			TransactionID: int64(transaction.ID),
			StorageID:     detail.StorageId,
			OilID:         detail.OilID,
		}
		transactionDetails = append(transactionDetails, transactionDetail)
	}

	if err := config.DB.Create(&transactionDetails).Error; err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	qrData := strconv.Itoa(int(transaction.ID))
	qrFile, err := utils.GenerateQRCode(qrData)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	key := fmt.Sprintf("qrcodes/%v.png", transaction.ID)
	qrURL, err := utils.UploadToS3(qrFile, key)
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
		err = utils.SendEmail(email, subject, body, qrFile)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
	}()

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
		statusId     = 1
	)

	pageParam := c.Query("page")
	if pageParam != "" {
		page, _ = strconv.Atoi(pageParam)
	}

	typeTransactionQuery := c.Query("type")

	statusIdParam := c.Query("status")
	statusIdParamInt, _ := strconv.Atoi(statusIdParam)
	if statusIdParamInt != 0 {
		statusId = statusIdParamInt
	}

	limitParam := c.Query("limit")
	if limitParam != "" {
		pageSize, _ = strconv.Atoi(limitParam)
	}

	db := config.DB

	var count int64
	if err := db.Model(&models.Transaction{}).
		Where("status_id = ?", statusId).
		Count(&count).
		Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	offset := (page - 1) * pageSize
	fmt.Println(statusId)

	fmt.Println("typeTransactionQuery", typeTransactionQuery, "halo")

	err := db.Offset(offset).Limit(pageSize).
		Where("status_id = ?", statusId).
		Preload("User.Company").
		Preload("Vehicle.VehicleType").
		Preload("Officer").
		Preload("Province").
		Preload("City").
		Preload("TransactionDetail.Oil").
		Order("updated_at desc").
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

func GetUserTransaction(c *gin.Context) {

	userId := c.Query("user_id")
	dateQuery := c.Query("date")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	dateNow := time.Now()
	var fromDate time.Time

	switch dateQuery {
	case "month":
		fromDate = dateNow.AddDate(0, -1, 0)
	case "week":
		fromDate = dateNow.AddDate(0, 0, -7)
	case "today":
		fromDate = dateNow
	case "all":
	default:
	}

	dbQuery := config.DB.
		Where("date >= ?", fromDate.Format("2006-01-02")).
		Preload("User.Company").
		Preload("TransactionDetail").
		Preload("Status.StatusType").
		Order("date desc")

	if userId != "" {
		dbQuery = dbQuery.Where("user_id = ?", userId)
	}

	var totalRecords int64
	if err := dbQuery.Model(&models.Transaction{}).Count(&totalRecords).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var totalPage int64
	if totalRecords > 0 {
		totalPage = (totalRecords + int64(pageSize) - 1) / int64(pageSize)
	}

	offset := (page - 1) * pageSize
	dbQuery = dbQuery.Limit(pageSize).Offset(offset)

	var userTransactions []models.Transaction
	if err := dbQuery.Find(&userTransactions).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data":        userTransactions,
		"totalCount":  totalRecords,
		"currentPage": page,
		"pageSize":    pageSize,
		"totalPage":   totalPage,
	})
}

func UpdateTransactionBatch(c *gin.Context) {
	type IdToUpdate struct {
		Id        uint64 `json:"id"`
		VehicleId *int   `json:"vehicle_id"`
		DriverId  int    `json:"driver_id"`
	}

	var ids []IdToUpdate
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	for _, id := range ids {
		transaction := models.Transaction{}
		if err := config.DB.Find(&transaction, id.Id).Error; err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		status := c.Query("status")

		statusInt, err := strconv.Atoi(status)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Please provide status",
			})
			return
		}

		transaction.VehicleId = id.VehicleId

		transaction.StatusId = statusInt

		transactionDate := transaction.Date.Format("2006-01-02")
		now := time.Now().Format("2006-01-02")

		if transactionDate == now {
			if statusInt == 4 {
				transaction.StatusId = 7
			}
		}

		if err := config.DB.Save(&transaction).Error; err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"message": "success update transaction",
	})

}

func UpdateTransaction(c *gin.Context) {
	transactionID := c.Param("id")

	var updateRequest models.TransactionUpdateInput
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	var transaction models.Transaction
	if err := config.DB.Preload("TransactionDetail").First(&transaction, transactionID).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "Transaksi tidak ditemukan",
		})
		return
	}

	transaction.OfficerId = updateRequest.OfficerID
	transaction.Date = updateRequest.Date
	transaction.ProvinceId = updateRequest.ProvinceID
	transaction.CityId = updateRequest.CityID
	transaction.VehicleId = updateRequest.VehicleID
	transaction.StatusId = updateRequest.StatusId

	// Update detail transaksi
	transaction.TransactionDetail = []models.TransactionDetail{}
	for _, detail := range updateRequest.TransactionDetails {
		transactionDetail := models.TransactionDetail{
			ID:            int64(detail.ID),
			TransactionID: int64(detail.TransactionID),
			OilID:         int64(detail.OilID),
			Quantity:      detail.Quantity,
			StorageID:     int64(detail.StorageID),
		}
		transaction.TransactionDetail = append(transaction.TransactionDetail, transactionDetail)
	}
	if err := config.DB.Save(&transaction).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Gagal memperbarui transaksi",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Berhasil memperbarui transaksi",
	})
}

func GetByIdTransaction(c *gin.Context) {
	var transaction models.Transaction
	id := c.Param("id")
	err := config.DB.
		Preload("TransactionDetail.Oil").
		Joins("Vehicle.VehicleType").
		Joins("User.Role").
		Joins("Officer").
		Joins("User.Detail").
		Joins("User.Company").
		Joins("Status.StatusType").
		Joins("Status.Status").
		Find(&transaction, id).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": transaction,
	})
}

func GetTransactionByUserId(c *gin.Context) {
	var transaction []models.Transaction
	var transactionResponse []models.TransactionResponse
	id := c.Param("id")
	fmt.Println(id)
	err := config.DB.
		Preload("TransactionDetail.Oil").
		Where("user_id = ?", id).
		Find(&transaction).Error
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

	err = dto.Map(&transactionResponse, &transaction)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"data": transactionResponse,
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

	today := time.Now().UTC().Format("2006-01-02")

	name := c.Query("name")

	if name != "" {
		err := config.DB.
			Preload("Vehicle.VehicleType").
			Preload("User.Company").
			Preload("User.Detail").
			Joins("JOIN users ON users.id = transactions.user_id").
			Where("users.username = ?", name).
			Where("date = ?", today).
			Find(&transaction).Error
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else {
		err := config.DB.Preload("Vehicle.VehicleType").
			Preload("User.Company").
			Preload("User.Detail").
			Where("date = ?", today).
			Find(&transaction).Error
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	if len(transaction) == 0 {
		c.JSON(404, gin.H{
			"message": "Transaction Today Not Found",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": transaction,
	})
}

func PostponeTransaction(c *gin.Context) {
	transactionID := c.Param("id")

	var postponeRequest struct {
		Reason string    `json:"reason"`
		Date   time.Time `json:"date"`
	}
	if err := c.ShouldBindJSON(&postponeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var transaction models.Transaction
	if err := config.DB.First(&transaction, transactionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi tidak ditemukan"})
		return
	}

	statusID := 7 // postponed status
	transaction.StatusId = statusID
	transaction.Date = postponeRequest.Date
	if err := config.DB.Save(&transaction).Error; err != nil {
		c.JSON(400, gin.H{"error": "Gagal mengubah status transaksi"})
	}

	postponeHistory := models.PostponeHistory{
		TransactionID: int(transaction.ID),
		Reason:        postponeRequest.Reason,
		Date:          postponeRequest.Date,
	}

	if err := config.DB.Create(&postponeHistory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat entri PostponeHistory"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction is successfully postponed"})
}

func UpdateTransactionType(c *gin.Context) {
	transactionID := c.Param("id")

	var typeRequest struct {
		StatusId int `json:"status_id"`
	}
	if err := c.ShouldBindJSON(&typeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var transaction models.Transaction
	if err := config.DB.First(&transaction, transactionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi tidak ditemukan"})
		return
	}

	transaction.StatusId = typeRequest.StatusId
	if err := config.DB.Save(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui tipe transaksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tipe transaksi berhasil diperbarui"})
}

func UpdateStatusTransactions(c *gin.Context) {
	transaction := models.Transaction{}

	id := c.Param("id")
	err := config.DB.Find(&transaction, id).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if transaction.ID == 0 {
		c.JSON(404, gin.H{
			"message": "Transaction Not Found",
		})
		return
	}

	status := c.Query("status")
	if status == "" {
		c.JSON(400, gin.H{
			"message": "Status is required",
		})
		return
	}

	statusQuery := c.Query("status")
	if statusQuery == "" {
		c.JSON(400, gin.H{
			"message": "Status is required",
		})
		return
	}

	statusInt, err := strconv.Atoi(statusQuery)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Status must be integer",
		})
		return
	}

	now := time.Now().Format("2006-01-02")
	transactionDate := transaction.Date.Format("2006-01-02")

	if transactionDate == now {
		if statusInt == 5 {
			transaction.StatusId = 6
		} else {
			transaction.StatusId = statusInt
		}
	} else {
		transaction.StatusId = statusInt
	}

	if err := config.DB.Save(&transaction).Error; err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": fmt.Sprintf("Success update status to %s transaction with id %s", status, id),
	})
}

func TestEndpoint(c *gin.Context) {
	transaction := models.Transaction{}
	err := config.DB.Find(&transaction, 2).Error
	if err != nil {
		fmt.Println(err)
		return
	}

	now := time.Now().Format("2006-01-02")
	transactionDate := transaction.Date.Format("2006-01-02")

	c.JSON(200, gin.H{
		"message": "success",
		"now":     now,
		"date":    transactionDate,
	})

}
