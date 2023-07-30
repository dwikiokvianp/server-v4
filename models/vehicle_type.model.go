package models

type VehicleType struct {
	Id   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

type VehicleIdentifier struct {
	Id         int    `json:"id" gorm:"primary_key"`
	Identifier string `json:"identifier"`
}
