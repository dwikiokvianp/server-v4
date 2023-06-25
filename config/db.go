package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"server-v2/models"
)

var DB *gorm.DB

func InitDatabase(url string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	migrateError := db.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Province{},
		&models.City{},
		&models.Oil{},
		&models.VehicleType{},
		&models.Vehicle{},
		&models.Transaction{},
		&models.Officer{},
		&models.Company{},
		&models.HistoryOut{},
		&models.HistoryIn{},
		&models.Detail{},
		&models.Proof{})
	if migrateError != nil {
		log.Fatalln(err)
		return nil
	}

	DB = db

	return nil

}
