package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/wangn-tech/tiny-douyin/internal/api/dto"
	"github.com/wangn-tech/tiny-douyin/internal/common/errc"
	"github.com/wangn-tech/tiny-douyin/internal/common/response"
	"github.com/wangn-tech/tiny-douyin/internal/global"
	"github.com/wangn-tech/tiny-douyin/internal/service"
)

// FavoriteHandler 点赞处理器
type FavoriteHandler struct {
	favoriteService service.IFavoriteService
}

// NewFavoriteHandler 创建 FavoriteHandler 实例（依赖注入）
func NewFavoriteHandler(favoriteService service.IFavoriteService) *FavoriteHandler {
	return &FavoriteHandler{
		favoriteService: favoriteService,
	}
}

// FavoriteAction 点赞操作
// POST /douyin/favorite/action/
// 参数：token（必填，JWT中间件处理），video_id（必填），action_type（必填，1-点赞，2-取消点赞）
func (h *FavoriteHandler) FavoriteAction(c *gin.Context) {
	ctx := c.Request.Context()

	// 从 JWT 中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		global.Logger.Warn("handler.FavoriteAction.missing_user_id")
		response.ErrorWithCode(c, errc.ErrUnauthorized)
		return
	}

	// 参数绑定和验证
	var req dto.FavoriteActionRequest
	if err := c.ShouldBind(&req); err != nil {
		global.Logger.Warn("handler.FavoriteAction.bind_error",
			zap.Error(err),
		)
		response.Error(c, errc.ErrInvalidParams, err.Error())
		return
	}

	global.Logger.Info("handler.FavoriteAction.request",
		zap.Uint("user_id", userID.(uint)),
		zap.Uint("video_id", req.VideoID),
		zap.Int32("action_type", req.ActionType),
		zap.String("client_ip", c.ClientIP()),
	)

	// 调用 Service 层处理业务逻辑
	err := h.favoriteService.FavoriteAction(ctx, userID.(uint), req.VideoID, req.ActionType)
	if err != nil {
		global.Logger.Error("handler.FavoriteAction.service_error",
			zap.Uint("user_id", userID.(uint)),
			zap.Uint("video_id", req.VideoID),
			zap.Int32("action_type", req.ActionType),
			zap.Error(err),
		)
		response.Error(c, errc.ErrInternalServer, err.Error())
		return
	}

	// 返回成功响应
	response.SuccessWithData(c, dto.FavoriteActionResponse{
		Response: response.Response{
			StatusCode: errc.Success,
			StatusMsg:  errc.GetMsg(errc.Success),
		},
	})
}

// GetFavoriteList 获取用户喜欢的视频列表
// GET /douyin/favorite/list/
// 参数：user_id（必填），token（可选）
func (h *FavoriteHandler) GetFavoriteList(c *gin.Context) {
	ctx := c.Request.Context()

	// 参数解析和验证
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		global.Logger.Warn("handler.GetFavoriteList.missing_param")
		response.Error(c, errc.ErrInvalidParams, "缺少user_id参数")
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		global.Logger.Warn("handler.GetFavoriteList.invalid_param",
			zap.String("user_id", userIDStr),
		)
		response.Error(c, errc.ErrInvalidParams, "user_id格式错误")
		return
	}

	// 获取当前用户ID（可选，用于判断 is_favorite）
	var currentUserID uint = 0
	if uid, exists := c.Get("user_id"); exists {
		currentUserID = uid.(uint)
	}

	global.Logger.Info("handler.GetFavoriteList.request",
		zap.Uint64("user_id", userID),
		zap.Uint("current_user_id", currentUserID),
		zap.String("client_ip", c.ClientIP()),
	)

	// 调用 Service 层处理业务逻辑
	data, err := h.favoriteService.GetFavoriteList(ctx, uint(userID), currentUserID)
	if err != nil {
		global.Logger.Error("handler.GetFavoriteList.service_error",
			zap.Uint64("user_id", userID),
			zap.Error(err),
		)
		response.Error(c, errc.ErrInternalServer, "获取喜欢列表失败")
		return
	}

	// 在 Handler 层封装响应
	resp := dto.FavoriteListResponse{
		Response: response.Response{
			StatusCode: errc.Success,
			StatusMsg:  errc.GetMsg(errc.Success),
		},
		Videos: data.Videos,
	}

	response.SuccessWithData(c, resp)
}
