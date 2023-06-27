package models

type TravelOrder struct {
	ID             int64  `gorm:"primary_key;auto_increment" json:"id"`
	OfficerID      int64  `gorm:"not null" json:"officer_id"`
	PickupLocation string `gorm:"not null" json:"pickup_location"`
	DepartureDate  string `gorm:"not null" json:"departure"`
	Message        string `gorm:"not null" json:"message"`
	Status         string `gorm:"not null" json:"status"`
}

type TravelDeliveryInput struct {
	OfficerID      int64
	PickupLocation int64
}
