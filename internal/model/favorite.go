package model

type Favorite struct {
	ID      uint `gorm:"primaryKey;autoIncrement"`
	VideoID uint `gorm:"not null;uniqueIndex:uk_user_video;priority:2"`
	UserID  uint `gorm:"not null;uniqueIndex:uk_user_video;priority:1"`
}

func (Favorite) TableName() string { return "favorites" }
