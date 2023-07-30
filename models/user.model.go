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
type UserResponse struct {
	Id        int             `json:"id"`
	Username  string          `json:"username"`
	Email     string          `json:"email"`
	Phone     string          `json:"phone"`
	CreatedAt int64           `gorm:"autoCreateTime" json:"created_at"`
	DetailId  int             `json:"detail_id"`
	Detail    Detail          `json:"detail"`
	CompanyID int             `json:"company_id"`
	Company   CompanyResponse `json:"company"`
}

type UserMinimumResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserInput struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	RoleId    int    `json:"role_id"`
	DetailId  int    `json:"detail_id"`
	CompanyID int    `json:"company_id"`
}

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	TokenString string `json:"access_token"`
	Name        string `json:"name"`
	Id          int    `json:"id"`
}

type Employee struct {
	Id     int `json:"id"`
	UserId int `json:"user_id"`
}

type Customer struct {
	Id     int `json:"id"`
	UserId int `json:"user_id"`
}
