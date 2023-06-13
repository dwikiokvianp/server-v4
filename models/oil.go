package models

type Oil struct {
	Id   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}
