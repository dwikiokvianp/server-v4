package models

type TransactionDelivery struct {
	ID             int64       `json:"id"`
	TransactionID  int64       `json:"transaction_id"`
	Transaction    Transaction `gorm:"foreignkey:TransactionID"`
	DeliveryStatus string      `json:"delivery_status"`
}
