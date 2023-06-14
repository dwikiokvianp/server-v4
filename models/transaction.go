package models

import "time"

type Transaction struct {
	ID        uint64    `gorm:"primary_key:auto_increment" json:"id"`
	UserId    int       `gorm:"not null" json:"user_id"`
	Email     string    `gorm:"not null" json:"email"`
	User      User      `gorm:"foreignkey:UserId"`
	QrCodeUrl string    `gorm:"not null" json:"qr_code_url"`
	VehicleId int       `gorm:"not null" json:"vehicle_id"`
	Vehicle   Vehicle   `gorm:"foreignkey:VehicleId"`
	OilId     int       `gorm:"not null" json:"oil_id"`
	Oil       Oil       `gorm:"foreignkey:OilId"`
	CreatedAt int64     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64     `gorm:"autoUpdateTime" json:"updated_at"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	OfficerId int       `gorm:"not null" json:"officer_id"`
	Officer   Officer   `gorm:"foreignkey:OfficerId"`
	IssuedAt  time.Time `json:"issued_at"`
	Status    string    `gorm:"not null" json:"status"`
}

type TransactionInput struct {
	UserId    int    `gorm:"not null" json:"user_id"`
	VehicleId int    `gorm:"not null" json:"vehicle_id"`
	OilId     int    `gorm:"not null" json:"oil_id"`
	Email     string `gorm:"not null" json:"email"`
	Quantity  int    `gorm:"not null" json:"quantity"`
	OfficerId int    `gorm:"not null" json:"officer_id"`
	QrCodeUrl string
}
