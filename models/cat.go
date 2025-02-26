package models

import (
	"gorm.io/gorm"
)

type Cat struct {
	gorm.Model
	Name     string  `json:"name" gorm:"type:varchar(100)"`
	YearsExp int     `json:"years_of_experience" gorm:"type:int"`
	Breed    string  `json:"breed" gorm:"type:varchar(100)"`
	Salary   float64 `json:"salary" gorm:"type:decimal(10,2)"`
}
