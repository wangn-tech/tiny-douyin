package model

import "time"

type Tag struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:varchar(64);uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

func (Tag) TableName() string { return "tags" }
