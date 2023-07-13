package models

type Status struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}
