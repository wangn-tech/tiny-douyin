package model

import "time"

type Favorite struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	VideoID   uint64    `gorm:"not null;uniqueIndex:uk_fav_user_video;priority:2"`
	UserID    uint64    `gorm:"not null;uniqueIndex:uk_fav_user_video;priority:1"`
	CreatedAt time.Time `gorm:"index"`
}

func (Favorite) TableName() string { return "favorites" }
