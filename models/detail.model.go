package models

type Detail struct {
	Id      int `json:"id" gorm:"primaryKey;autoIncrement"`
	Balance int `json:"balance"`
	Credit  int `json:"credit"`
}
