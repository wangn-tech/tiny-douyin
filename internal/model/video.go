package model

import "time"

type Video struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	AuthorID    uint64    `gorm:"index;not null"`
	PlayURL     string    `gorm:"type:varchar(255);not null"`
	CoverURL    string    `gorm:"type:varchar(255)"`
	Title       string    `gorm:"type:varchar(128)"`
	Description string    `gorm:"type:varchar(255)"`
	DurationSec int       `gorm:"type:int"`
	Visibility  int8      `gorm:"type:tinyint;default:0"` // 0:public 1:private 2:friends
	CreatedAt   time.Time `gorm:"index"`
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
}

func (Video) TableName() string { return "videos" }
