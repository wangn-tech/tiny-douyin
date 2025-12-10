package model

import "time"

type Comment struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	VideoID   uint64    `gorm:"index;not null"`
	UserID    uint64    `gorm:"index;not null"`
	Content   string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

func (Comment) TableName() string { return "comments" }
