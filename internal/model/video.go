package model

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	AuthorID    uint   `gorm:"index;not null"`
	PlayURL     string `gorm:"type:varchar(255);not null"`
	CoverURL    string `gorm:"type:varchar(255)"`
	Title       string `gorm:"type:varchar(128)"`
	Description string `gorm:"type:varchar(255)"`
}

func (Video) TableName() string { return "videos" }
