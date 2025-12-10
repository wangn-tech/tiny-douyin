package model

import "time"

type Message struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement"`
	SenderID   uint64    `gorm:"index;not null"`
	ReceiverID uint64    `gorm:"index;not null"`
	Content    string    `gorm:"type:varchar(1024);not null"`
	CreatedAt  time.Time `gorm:"index"`
}

func (Message) TableName() string { return "messages" }
