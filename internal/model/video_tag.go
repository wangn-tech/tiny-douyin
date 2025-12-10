package model

type VideoTag struct {
	VideoID uint64 `gorm:"primaryKey"`
	TagID   uint64 `gorm:"primaryKey"`
}

func (VideoTag) TableName() string { return "video_tags" }
