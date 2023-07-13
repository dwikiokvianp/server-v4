package models

type Handover struct {
	Id            	 int    `json:"id"`
	WorkerBeforeId 	 int    `json:"worker_before_id"`
	WorkerBefore  	 User   `json:"worker_before" gorm:"foreignKey:WorkerBeforeId"`
	WorkerAfterId 	 int    `json:"worker_after_id"`
	WorkerAfter   	 User   `json:"worker_after" gorm:"foreignKey:WorkerAfterId"`
	OfficerId      	 int    `json:"officer_id"`
	Officer       	 User   `json:"officer" gorm:"foreignKey:OfficerId"`
	Condition      	 string `json:"condition"`
	Status         	 string `json:"status"`
	PhotoTangki    	 string `json:"photo_tangki"`
	PhotoKebersihan  string `json:"photo_kebersihan"`
	PhotoLevelGauge  string `json:"photo_level_gauge"`
	PhotoPetugas     string `json:"photo_petugas"`
}

type HandoverResponse struct {
	Id            	 int                 `json:"id"`
	WorkerBeforeId	 int                 `json:"worker_before_id"`
	WorkerBefore  	 UserMinimumResponse `json:"worker_before" gorm:"foreignKey:WorkerBeforeId"`
	WorkerAfterId 	 int                 `json:"worker_after_id"`
	WorkerAfter   	 UserMinimumResponse `json:"worker_after" gorm:"foreignKey:WorkerAfterId"`
	OfficerId     	 int                 `json:"officer_id"`
	Officer       	 UserMinimumResponse `json:"officer" gorm:"foreignKey:OfficerId"`
	Condition     	 string              `json:"condition"`
	Status        	 string              `json:"status"`
	PhotoTangki    	 string 		     `json:"photo_tangki"`
	PhotoKebersihan  string 		     `json:"photo_kebersihan"`
	PhotoLevelGauge  string              `json:"photo_level_gauge"`
	PhotoPetugas     string              `json:"photo_petugas"`
}
