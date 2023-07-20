package models

type Status struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

type StatusType struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

type StatusTypeMapping struct {
	ID           int `json:"id" gorm:"primaryKey"`
	StatusID     int `json:"status_id"`
	Status       Status
	StatusTypeID int `json:"status_type_id"`
	StatusType   StatusType
}

type StatusTypeMappingResponse struct {
	ID         int    `json:"id"`
	Status     string `json:"name"`
	StatusType string `json:"status_type"`
}
