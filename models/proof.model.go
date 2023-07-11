package models

import (
	"time"
)

type Proof struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	PhotoKTPURL    string    `json:"photo_ktp_url"`
	Description    string    `json:"description"`
	InvoiceURL     string    `json:"invoice_url"`
	// SignatureURL   string    `json:"signature_url"`
	PhotoOrangURL  string    `json:"photo_orang_url"`
	PhotoTangkiURL string    `json:"photo_tangki_url"`
	CreatedAt      time.Time `json:"created_at"`
	TransactionID  int       `json:"transaction_id"`
	Transaction    Transaction
}
