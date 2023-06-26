package models

type TransactionDetail struct {
	ID            int64 `gorm:"primary_key;auto_increment" json:"id"`
	OilID         int64 `gorm:"not null" json:"oil_id"`
	TransactionID int64 `gorm:"not null" json:"transaction_id"`
}
