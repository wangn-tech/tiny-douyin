package model

import "gorm.io/gorm"

// Relation 关注关系
type Relation struct {
	gorm.Model
	FollowerID uint `gorm:"not null;uniqueIndex:uk_follow_pair;index:idx_follower_id;comment:关注者ID"`  // 关注者ID
	FolloweeID uint `gorm:"not null;uniqueIndex:uk_follow_pair;index:idx_followee_id;comment:被关注者ID"` // 被关注者ID
}

func (Relation) TableName() string { return "relations" }
