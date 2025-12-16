package model

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	AuthorID      uint   `gorm:"index;not null"`
	PlayURL       string `gorm:"type:varchar(255);not null"`
	CoverURL      string `gorm:"type:varchar(255)"`
	Title         string `gorm:"type:varchar(128)"`
	Description   string `gorm:"type:varchar(255)"`
	FavoriteCount int64  `gorm:"default:0;not null"` // 点赞数
	CommentCount  int64  `gorm:"default:0;not null"` // 评论数
}

func (Video) TableName() string { return "videos" }
