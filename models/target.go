package models

import "gorm.io/gorm"

type Target struct {
	gorm.Model
	Name      string `json:"name" gorm:"type:varchar(100)"`
	Country   string `json:"country" gorm:"type:varchar(100)"`
	Notes     string `json:"notes" gorm:"type:text"`
	Complete  bool   `json:"complete"`
	MissionID *uint
}
