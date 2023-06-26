package models

type City struct {
	Id         int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string `json:"name"`
	ProvinceID int    `json:"province_id"`
}
