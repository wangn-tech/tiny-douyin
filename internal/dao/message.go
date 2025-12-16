package dao

import (
	"context"
	"time"

	"github.com/wangn-tech/tiny-douyin/internal/common/constant"
	"github.com/wangn-tech/tiny-douyin/internal/global"
	"github.com/wangn-tech/tiny-douyin/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// IMessageDAO 消息数据访问接口
type IMessageDAO interface {
	// CreateMessage 创建消息记录
	CreateMessage(ctx context.Context, message *model.Message) error
	// GetChatMessages 获取两个用户之间的聊天记录
	GetChatMessages(ctx context.Context, userID1, userID2 uint, preMsgTime int64) ([]*model.Message, error)
	// GetLatestMessage 获取两个用户之间的最新一条消息
	GetLatestMessage(ctx context.Context, userID1, userID2 uint) (*model.Message, error)
}

// MessageDAO 消息数据访问实现
type MessageDAO struct {
	db *gorm.DB
}

// NewMessageDAO 创建 MessageDAO 实例
func NewMessageDAO(db *gorm.DB) IMessageDAO {
	return &MessageDAO{
		db: db,
	}
}

// CreateMessage 创建消息记录
func (d *MessageDAO) CreateMessage(ctx context.Context, message *model.Message) error {
	err := d.db.WithContext(ctx).Create(message).Error

	if err != nil {
		global.Logger.Error("dao.CreateMessage.failed",
			zap.Uint("from_user_id", message.FromUserID),
			zap.Uint("to_user_id", message.ToUserID),
			zap.Error(err),
		)
	}

	return err
}

// GetChatMessages 获取两个用户之间的聊天记录（双向查询）
// preMsgTime: Unix时间戳（秒），查询比这个时间更早的消息（用于分页）
func (d *MessageDAO) GetChatMessages(ctx context.Context, userID1, userID2 uint, preMsgTime int64) ([]*model.Message, error) {
	var messages []*model.Message

	query := d.db.WithContext(ctx).
		Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)",
			userID1, userID2, userID2, userID1)

	// 如果提供了 preMsgTime，查询比这个时间更早的消息
	if preMsgTime > 0 {
		preTime := time.Unix(preMsgTime, 0)
		query = query.Where("created_at < ?", preTime)
	}

	err := query.
		Order("created_at DESC").
		Limit(constant.MessagePageSize).
		Find(&messages).Error

	if err != nil {
		global.Logger.Error("dao.GetChatMessages.failed",
			zap.Uint("user_id_1", userID1),
			zap.Uint("user_id_2", userID2),
			zap.Int64("pre_msg_time", preMsgTime),
			zap.Error(err),
		)
		return nil, err
	}

	return messages, nil
}

// GetLatestMessage 获取两个用户之间的最新一条消息
func (d *MessageDAO) GetLatestMessage(ctx context.Context, userID1, userID2 uint) (*model.Message, error) {
	var message model.Message

	err := d.db.WithContext(ctx).
		Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)",
			userID1, userID2, userID2, userID1).
		Order("created_at DESC").
		First(&message).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		global.Logger.Error("dao.GetLatestMessage.failed",
			zap.Uint("user_id_1", userID1),
			zap.Uint("user_id_2", userID2),
			zap.Error(err),
		)
		return nil, err
	}

	return &message, nil
}
