package model

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	Name string `gorm:"type:varchar(64);uniqueIndex;not null"`
}

func (Tag) TableName() string { return "tags" }
