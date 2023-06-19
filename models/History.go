package models

import "time"

type HistoryOut struct {
	Id            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Date          time.Time `json:"date"`
	Quantity      int       `json:"quantity"`
	OilId         int       `json:"oil_id"`
	Oil           Oil
	User          User
	UserId        int `json:"user_id"`
	TransactionId int `json:"transaction_id"`
}

type HistoryIn struct {
	Date     time.Time `json:"date"`
	Quantity int       `json:"quantity"`
	OilId    int       `json:"oil_id"`
	Oil      Oil
}
