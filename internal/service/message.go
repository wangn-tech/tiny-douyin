package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/wangn-tech/tiny-douyin/internal/api/dto"
	"github.com/wangn-tech/tiny-douyin/internal/common/constant"
	"github.com/wangn-tech/tiny-douyin/internal/dao"
	"github.com/wangn-tech/tiny-douyin/internal/global"
	"github.com/wangn-tech/tiny-douyin/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// IMessageService 消息服务接口
type IMessageService interface {
	// SendMessage 发送消息
	SendMessage(ctx context.Context, fromUserID uint, req *dto.MessageActionRequest) error
	// GetChatMessages 获取聊天记录
	GetChatMessages(ctx context.Context, currentUserID uint, req *dto.MessageChatRequest) (*dto.MessageChatResponse, error)
}

// MessageService 消息服务实现
type MessageService struct {
	messageDAO  dao.IMessageDAO
	userDAO     dao.IUserDAO
	relationSvc IRelationService
}

// NewMessageService 创建 MessageService 实例
func NewMessageService(
	messageDAO dao.IMessageDAO,
	userDAO dao.IUserDAO,
	relationSvc IRelationService,
) IMessageService {
	return &MessageService{
		messageDAO:  messageDAO,
		userDAO:     userDAO,
		relationSvc: relationSvc,
	}
}

// SendMessage 发送消息
func (s *MessageService) SendMessage(ctx context.Context, fromUserID uint, req *dto.MessageActionRequest) error {
	// 验证消息内容
	content := strings.TrimSpace(req.Content)
	if len(content) == 0 {
		global.Logger.Warn("service.SendMessage.content_empty",
			zap.Uint("from_user_id", fromUserID),
			zap.Uint("to_user_id", req.ToUserID),
		)
		return fmt.Errorf("消息内容不能为空")
	}

	if len(content) > constant.MessageMaxLength {
		global.Logger.Warn("service.SendMessage.content_too_long",
			zap.Uint("from_user_id", fromUserID),
			zap.Uint("to_user_id", req.ToUserID),
			zap.Int("length", len(content)),
		)
		return fmt.Errorf("消息内容过长，最多%d个字符", constant.MessageMaxLength)
	}

	// 验证接收者用户是否存在
	toUser, err := s.userDAO.GetUserByID(ctx, req.ToUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.Warn("service.SendMessage.to_user_not_found",
				zap.Uint("to_user_id", req.ToUserID),
			)
			return fmt.Errorf("接收者用户不存在")
		}
		global.Logger.Error("service.SendMessage.get_user_error",
			zap.Uint("to_user_id", req.ToUserID),
			zap.Error(err),
		)
		return fmt.Errorf("查询用户失败")
	}

	// 验证是否为好友关系（双向关注）
	isFriend, err := s.relationSvc.IsFriend(ctx, fromUserID, req.ToUserID)
	if err != nil {
		global.Logger.Error("service.SendMessage.check_friend_error",
			zap.Uint("from_user_id", fromUserID),
			zap.Uint("to_user_id", req.ToUserID),
			zap.Error(err),
		)
		return fmt.Errorf("检查好友关系失败")
	}

	if !isFriend {
		global.Logger.Warn("service.SendMessage.not_friend",
			zap.Uint("from_user_id", fromUserID),
			zap.Uint("to_user_id", req.ToUserID),
		)
		return fmt.Errorf("只能给好友发送消息")
	}

	// 创建消息记录
	message := &model.Message{
		FromUserID: fromUserID,
		ToUserID:   req.ToUserID,
		Content:    content,
	}

	err = s.messageDAO.CreateMessage(ctx, message)
	if err != nil {
		global.Logger.Error("service.SendMessage.create_message_error",
			zap.Uint("from_user_id", fromUserID),
			zap.Uint("to_user_id", req.ToUserID),
			zap.Error(err),
		)
		return fmt.Errorf("发送消息失败")
	}

	global.Logger.Info("service.SendMessage.success",
		zap.Uint("from_user_id", fromUserID),
		zap.Uint("to_user_id", req.ToUserID),
		zap.String("to_username", toUser.Username),
		zap.Uint("message_id", message.ID),
	)

	return nil
}

// GetChatMessages 获取聊天记录
func (s *MessageService) GetChatMessages(ctx context.Context, currentUserID uint, req *dto.MessageChatRequest) (*dto.MessageChatResponse, error) {
	// 验证对方用户是否存在
	toUser, err := s.userDAO.GetUserByID(ctx, req.ToUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.Warn("service.GetChatMessages.to_user_not_found",
				zap.Uint("to_user_id", req.ToUserID),
			)
			return nil, fmt.Errorf("对方用户不存在")
		}
		global.Logger.Error("service.GetChatMessages.get_user_error",
			zap.Uint("to_user_id", req.ToUserID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("查询用户失败")
	}

	// 查询聊天记录（双向消息）
	messages, err := s.messageDAO.GetChatMessages(ctx, currentUserID, req.ToUserID, req.PreMsgTime)
	if err != nil {
		global.Logger.Error("service.GetChatMessages.query_error",
			zap.Uint("current_user_id", currentUserID),
			zap.Uint("to_user_id", req.ToUserID),
			zap.Int64("pre_msg_time", req.PreMsgTime),
			zap.Error(err),
		)
		return nil, fmt.Errorf("查询聊天记录失败")
	}

	// 转换为 DTO
	messageList := make([]dto.Message, 0, len(messages))
	for _, msg := range messages {
		messageList = append(messageList, dto.Message{
			ID:         msg.ID,
			ToUserID:   msg.ToUserID,
			FromUserID: msg.FromUserID,
			Content:    msg.Content,
			CreateTime: msg.CreatedAt.Unix(), // 转换为 Unix 时间戳（秒）
		})
	}

	global.Logger.Info("service.GetChatMessages.success",
		zap.Uint("current_user_id", currentUserID),
		zap.Uint("to_user_id", req.ToUserID),
		zap.String("to_username", toUser.Username),
		zap.Int("message_count", len(messageList)),
	)

	return &dto.MessageChatResponse{
		StatusCode:  0,
		StatusMsg:   "success",
		MessageList: messageList,
	}, nil
}
