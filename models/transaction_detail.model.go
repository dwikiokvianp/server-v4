package models

type TransactionDetail struct {
	ID            int64       `gorm:"primary_key;auto_increment" json:"id"`
	TransactionID int64       `gorm:"not null" json:"transaction_id"`
	Transaction   Transaction `gorm:"foreignKey:TransactionID" json:"transaction"`
	Quantity      int64       `gorm:"not null" json:"quantity"`
	StorageID     int64       `gorm:"not null" json:"storage_id"`
}
