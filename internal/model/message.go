package model

import "gorm.io/gorm"

// Message 消息模型
type Message struct {
	gorm.Model
	FromUserID uint   `gorm:"index:idx_from_to;index:idx_from;not null;comment:发送者用户ID"`
	ToUserID   uint   `gorm:"index:idx_from_to;index:idx_to;not null;comment:接收者用户ID"`
	Content    string `gorm:"type:varchar(255);not null;comment:消息内容"`
}

// TableName 指定表名
func (Message) TableName() string {
	return "messages"
}
