package models

type Vehicle struct {
	Id                  int    `json:"id" gorm:"primary_key"`
	Name                string `json:"name"`
	VehicleTypeId       int    `json:"vehicle_type_id" gorm:"foreignkey:VehicleTypeId"`
	VehicleType         VehicleType
	VehicleIdentifierId int `json:"vehicle_identifier_id" gorm:"foreignkey:VehicleIdentifierId"`
	VehicleIdentifier   VehicleIdentifier
}
