package model

import "time"

type Relation struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement"`
	FollowerID uint64    `gorm:"not null;uniqueIndex:uk_follow_pair;priority:1"` // 关注者
	FolloweeID uint64    `gorm:"not null;uniqueIndex:uk_follow_pair;priority:2"` // 被关注者
	CreatedAt  time.Time `gorm:"index"`
}

func (Relation) TableName() string { return "relations" }
