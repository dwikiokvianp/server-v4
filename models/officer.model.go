package models

type Officer struct {
	Id          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Username    string `gorm:"not null" json:"username"`
	Password    string `gorm:"not null" json:"password"`
	Email       string `gorm:"not null" json:"email"`
	CreatedAt   int64  `gorm:"autoCreateTime" json:"created_at"`
	WarehouseId int    `gorm:"not null" json:"warehouse_id"`
	PhoneNumber string `gorm:"not null" json:"phone_number"`
}
