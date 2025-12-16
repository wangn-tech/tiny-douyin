package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string `gorm:"type:varchar(32);uniqueIndex;not null"`
	Password  string `gorm:"type:char(60);not null"`
	Nickname  string `gorm:"type:varchar(32)"`
	Avatar    string `gorm:"type:varchar(255)"`
	Signature string `gorm:"type:varchar(255)"`
}

func (User) TableName() string { return "users" }
