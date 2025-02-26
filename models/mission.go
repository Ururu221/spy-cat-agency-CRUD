package models

import (
	"gorm.io/gorm"
)

type Mission struct {
	gorm.Model
	Targets  []Target `gorm:"foreignKey:MissionID"`
	CatID    *uint
	Cat      Cat  `gorm:"foreignKey:CatID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Complete bool `json:"complete"`
}
