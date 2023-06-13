package models

type Transaction struct {
	ID        uint64 `gorm:"primary_key:auto_increment" json:"id"`
	UserId    int    `gorm:"not null" json:"user_id"`
	User      User
	VehicleId int `gorm:"not null" json:"vehicle_id" `
	Vehicle   Vehicle
	OilId     int `gorm:"not null" json:"oil_id" `
	Oil       Oil
	QrCodeUrl string `gorm:"not null" json:"qr_code_url" `
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
}

type TransactionInput struct {
	UserId    int `gorm:"not null" json:"user_id"`
	VehicleId int `gorm:"not null" json:"vehicle_id" `
	OilId     int `gorm:"not null" json:"oil_id" `
	QrCodeUrl string
}
