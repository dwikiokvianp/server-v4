package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"server-v2/config"
	"server-v2/models"
	"server-v2/routes"
	"server-v2/utils"
	"strconv"
	"time"
)

func main() {
	failedLoadEnv := godotenv.Load(".env")
	if failedLoadEnv != nil {
		log.Fatal("Error loading .env file")
	}

	databaseUrl, port := os.Getenv("DB_URL"), os.Getenv("PORT")
	config.InitDatabase(databaseUrl)

	server := gin.New()
	myCorsConfig := cors.DefaultConfig()
	myCorsConfig.AllowAllOrigins = true
	myCorsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	server.Use(cors.New(myCorsConfig))
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Check github action",
		})
	})
	routes.Routes(server)

	c := cron.New()

	c.AddFunc("13 12 * * *", func() {
		var transactions []models.Transaction
		now := time.Now().Format("2006-01-02")
		if err := config.DB.
			Where("status_id = ? AND DATE(date) = ?", 3, now).
			Find(&transactions).Error; err != nil {
			log.Println("Gagal mendapatkan transaksi yang akan ditunda:", err.Error())
			return
		}

		if transactions == nil {
			log.Println("Tidak ada transaksi yang akan ditunda")
			return
		}

		for _, transaction := range transactions {
			transaction.StatusId = 7
			if err := config.DB.Save(&transaction).Error; err != nil {
				log.Println("Gagal memperbarui status transaksi:", err.Error())
				continue
			}

			postponeHistory := models.PostponeHistory{
				TransactionID: int(transaction.ID),
				Reason:        "automatic postponed",
			}
			if err := config.DB.Create(&postponeHistory).Error; err != nil {
				log.Println("Gagal membuat entri PostponeHistory:", err.Error())
				continue
			}

			var user models.User
			if err := config.DB.Where("role_id = ?", 1).First(&user).Error; err != nil {
				log.Println("Gagal mendapatkan email admin:", err.Error())
				continue
			}

			to := user.Email
			subject := "Transaksi Ditunda"
			body := "Transaksi dengan ID " + strconv.Itoa(int(transaction.ID)) + " telah ditunda."
			if err := utils.SendEmailNotification(to, subject, body); err != nil {
				log.Println("Gagal mengirim notifikasi email:", err.Error())
			}
		}

		log.Println("Cron job selesai: Transaksi berhasil diubah")
	})

	c.AddFunc("* * * * *", func() {
		var transactions []models.Transaction
		now := time.Now().Format("2006-01-02")
		if err := config.DB.
			Where("status_id = ? AND DATE(date) = ? AND type = ?", 2, now, "pickup").
			Find(&transactions).Error; err != nil {
			log.Println("Gagal mendapatkan transaksi yang akan diubah statusnya:", err.Error())
			return
		}

		if len(transactions) == 0 {
			log.Println("Tidak ada transaksi yang akan diubah statusnya")
			return
		}

		for _, transaction := range transactions {
			transaction.StatusId = 4
			if err := config.DB.Save(&transaction).Error; err != nil {
				log.Println("Gagal memperbarui status transaksi:", err.Error())
				continue
			}

			var user models.User
			if err := config.DB.Where("role_id = ?", 1).First(&user).Error; err != nil {
				log.Println("Gagal mendapatkan email admin:", err.Error())
				continue
			}

			to := user.Email
			subject := "Transaksi Diubah"
			body := "Transaksi dengan ID " + strconv.Itoa(int(transaction.ID)) + " telah diubah statusnya menjadi 'Pickup'."
			if err := utils.SendEmailNotification(to, subject, body); err != nil {
				log.Println("Gagal mengirim notifikasi email:", err.Error())
			}
		}

		log.Println("Cron job selesai: Transaksi berhasil diubah statusnya")
	})

	c.Start()

	err := server.Run(port)
	if err != nil {
		log.Fatal("Error running server")
	}
}
