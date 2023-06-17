package models

type User struct {
	Id        int     `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string  `json:"username"`
	Password  string  `json:"password"`
	Email     string  `json:"email"`
	RoleId    int     `json:"role_id"`
	Role      Role    `json:"role" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DetailId  int     `json:"detail_id"`
	Detail    Detail  `json:"detail" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Company   Company `json:"company" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CompanyID int     `json:"company_id"`
	CreatedAt int64   `gorm:"autoCreateTime" json:"created_at"`
	Phone     string  `json:"phone"`
}

type UserInput struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	RoleId    int    `json:"role_id"`
	DetailId  int    `json:"detail_id"`
	CompanyID int    `json:"company_id"`
}

type UserResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	TokenString string `json:"access_token"`
}
