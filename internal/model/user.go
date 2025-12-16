package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username      string `gorm:"type:varchar(32);uniqueIndex;not null"`
	Password      string `gorm:"type:char(60);not null"`
	Nickname      string `gorm:"type:varchar(32)"`
	Avatar        string `gorm:"type:varchar(255)"`
	Signature     string `gorm:"type:varchar(255)"`
	FollowCount   int64  `gorm:"type:bigint;default:0;not null;comment:关注数"` // 关注数
	FollowerCount int64  `gorm:"type:bigint;default:0;not null;comment:粉丝数"` // 粉丝数
}

func (User) TableName() string { return "users" }
