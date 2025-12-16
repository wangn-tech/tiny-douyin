package model

type VideoTag struct {
	VideoID uint `gorm:"primaryKey"`
	TagID   uint `gorm:"primaryKey"`
}

func (VideoTag) TableName() string { return "video_tags" }
