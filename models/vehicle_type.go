package models

type VehicleType struct {
	Id   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
	VehiclePhoto string `json:"vehicle_photo"`
}
