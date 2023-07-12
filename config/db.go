package config

import (
	"github.com/golang-jwt/jwt"
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

	migrateError := db.AutoMigrate(modelsToMigrate...)
	if migrateError != nil {
		log.Fatalln(err)
		return nil
	}

	DB = db
	return nil
}

func DropDatabase(url string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	dropError := db.Migrator().DropTable(modelsToMigrate...)
	if dropError != nil {
		log.Fatalln(err)
		return nil
	}

	return nil
}

var modelsToMigrate = []interface{}{
	&models.Role{},
	&models.DeliveryOrder{},
	&models.TravelOrder{},
	&models.Warehouse{},
	&models.WarehouseDetail{},
	&models.User{},
	&models.Oil{},
	&models.VehicleType{},
	&models.Vehicle{},
	&models.Transaction{},
	&models.TransactionDetail{},
	&models.Officer{},
	&models.Company{},
	&models.HistoryOut{},
	&models.HistoryIn{},
	&models.Detail{},
	&models.Storage{},
	&models.Driver{},
	&models.Proof{},
	&models.DeliveryOrderRecipientDetail{},
	&models.DeliveryOrderWarehouseDetail{},
	&models.Handover{},
}

type JWTClaim struct {
	Username string `json:"username"`
	UserId   uint   `json:"userId"`
	Role     string `json:"role"`
	jwt.StandardClaims
}
