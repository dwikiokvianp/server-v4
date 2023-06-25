package models

type Province struct {
	Id     int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name   string `json:"name"`
	CityId int    `json:"city_id"`
}
