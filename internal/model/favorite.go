package model

import (
	"time"

	"gorm.io/gorm"
)

// Favorite 点赞记录模型
type Favorite struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	UserID    uint           `gorm:"not null;uniqueIndex:uk_user_video;priority:1;index:idx_user"`  // 用户ID
	VideoID   uint           `gorm:"not null;uniqueIndex:uk_user_video;priority:2;index:idx_video"` // 视频ID
	CreatedAt time.Time      // 点赞时间
	DeletedAt gorm.DeletedAt `gorm:"index"` // 软删除
}

func (Favorite) TableName() string { return "favorites" }
