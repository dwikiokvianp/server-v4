package models

type Employee struct {
	Id     int `json:"id"`
	UserId int `json:"user_id"`
	User   User
	RoleId int  `json:"role_id"`
	Role   Role `json:"role" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
