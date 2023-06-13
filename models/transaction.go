package models

type Transaction struct {
	ID        uint64 `gorm:"primary_key:auto_increment" json:"id"`
	UserId    int    `gorm:"not null" json:"user_id"`
	Email     string `gorm:"not null" json:"email"`
	User      User   `gorm:"foreignkey:UserId"`
	VehicleId int    `gorm:"not null" json:"vehicle_id"`
	Vehicle   Vehicle `gorm:"foreignkey:VehicleId"`
	OilId     int    `gorm:"not null" json:"oil_id"`
	Oil       Oil    `gorm:"foreignkey:OilId"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
}

type TransactionInput struct {
	UserId    int    `gorm:"not null" json:"user_id"`
	VehicleId int    `gorm:"not null" json:"vehicle_id"`
	OilId     int    `gorm:"not null" json:"oil_id"`
	Email     string `gorm:"not null" json:"email"`
}
