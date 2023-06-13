package models

type Vehicle struct {
	Id            int    `json:"id" gorm:"primary_key"`
	Name          string `json:"name"`
	VehicleTypeId int    `json:"vehicle_type_id" gorm:"foreignkey:VehicleTypeId"`
	VehicleType   VehicleType
}
