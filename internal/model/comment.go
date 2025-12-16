package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	VideoID uint   `gorm:"index;not null"`
	UserID  uint   `gorm:"index;not null"`
	Content string `gorm:"type:varchar(255);not null"`
}

func (Comment) TableName() string { return "comments" }
