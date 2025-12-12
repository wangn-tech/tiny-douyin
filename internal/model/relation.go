package model

type Relation struct {
	ID         uint `gorm:"primaryKey;autoIncrement"`
	FollowerID uint `gorm:"not null;uniqueIndex:uk_follow_pair;priority:1"`
	FolloweeID uint `gorm:"not null;uniqueIndex:uk_follow_pair;priority:2"`
}

func (Relation) TableName() string { return "relations" }
