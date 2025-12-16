package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wangn-tech/tiny-douyin/internal/api/dto"
	"github.com/wangn-tech/tiny-douyin/internal/common/errc"
	"github.com/wangn-tech/tiny-douyin/internal/common/response"
	"github.com/wangn-tech/tiny-douyin/internal/global"
	"github.com/wangn-tech/tiny-douyin/internal/service"
	"go.uber.org/zap"
)

// MessageHandler 消息处理器
type MessageHandler struct {
	messageService service.IMessageService
}

// NewMessageHandler 创建 MessageHandler 实例
func NewMessageHandler(messageService service.IMessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

// SendMessage 发送消息
// @Summary 发送消息
// @Description 用户向好友发送消息
// @Tags 消息
// @Accept json
// @Produce json
// @Param token query string true "用户token"
// @Param to_user_id query uint true "接收者用户ID"
// @Param action_type query int true "操作类型：1-发送消息"
// @Param content query string true "消息内容"
// @Success 200 {object} dto.MessageActionResponse
// @Router /douyin/message/action/ [post]
func (h *MessageHandler) SendMessage(c *gin.Context) {
	// 获取当前用户ID（从JWT中间件中获取）
	userID, exists := c.Get("user_id")
	if !exists {
		global.Logger.Warn("handler.SendMessage.user_id_not_found")
		response.Error(c, errc.ErrUnauthorized, "未授权")
		return
	}

	// 绑定请求参数
	var req dto.MessageActionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		global.Logger.Warn("handler.SendMessage.bind_error",
			zap.Error(err),
		)
		response.Error(c, errc.ErrInvalidParams, "参数错误")
		return
	}

	// 调用 Service 层
	err := h.messageService.SendMessage(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		global.Logger.Error("handler.SendMessage.service_error",
			zap.Uint("from_user_id", userID.(uint)),
			zap.Uint("to_user_id", req.ToUserID),
			zap.Error(err),
		)
		response.Error(c, errc.Failed, err.Error())
		return
	}

	// 返回成功响应
	resp := &dto.MessageActionResponse{
		StatusCode: errc.Success,
		StatusMsg:  "发送成功",
	}

	c.JSON(http.StatusOK, resp)

	global.Logger.Info("handler.SendMessage.success",
		zap.Uint("from_user_id", userID.(uint)),
		zap.Uint("to_user_id", req.ToUserID),
	)
}

// GetChatMessages 获取聊天记录
// @Summary 获取聊天记录
// @Description 获取与某个用户的聊天记录
// @Tags 消息
// @Accept json
// @Produce json
// @Param token query string true "用户token"
// @Param to_user_id query uint true "对方用户ID"
// @Param pre_msg_time query int64 false "上次最新消息的时间戳（秒）"
// @Success 200 {object} dto.MessageChatResponse
// @Router /douyin/message/chat/ [get]
func (h *MessageHandler) GetChatMessages(c *gin.Context) {
	// 获取当前用户ID（从JWT中间件中获取）
	userID, exists := c.Get("user_id")
	if !exists {
		global.Logger.Warn("handler.GetChatMessages.user_id_not_found")
		response.Error(c, errc.ErrUnauthorized, "未授权")
		return
	}

	// 绑定请求参数
	var req dto.MessageChatRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		global.Logger.Warn("handler.GetChatMessages.bind_error",
			zap.Error(err),
		)
		response.Error(c, errc.ErrInvalidParams, "参数错误")
		return
	}

	// 调用 Service 层
	resp, err := h.messageService.GetChatMessages(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		global.Logger.Error("handler.GetChatMessages.service_error",
			zap.Uint("current_user_id", userID.(uint)),
			zap.Uint("to_user_id", req.ToUserID),
			zap.Error(err),
		)
		response.Error(c, errc.Failed, err.Error())
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, resp)

	global.Logger.Info("handler.GetChatMessages.success",
		zap.Uint("current_user_id", userID.(uint)),
		zap.Uint("to_user_id", req.ToUserID),
		zap.Int("message_count", len(resp.MessageList)),
	)
}
