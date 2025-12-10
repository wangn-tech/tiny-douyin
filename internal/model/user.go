package model

import (
	"time"
)

type User struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	Username  string    `gorm:"type:varchar(32);uniqueIndex;not null"`
	Password  string    `gorm:"type:char(60);not null"` // bcrypt
	Nickname  string    `gorm:"type:varchar(32)"`
	Avatar    string    `gorm:"type:varchar(255)"`
	Signature string    `gorm:"type:varchar(255)"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

func (User) TableName() string { return "users" }
