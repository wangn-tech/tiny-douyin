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

// RelationHandler 关注处理器
type RelationHandler struct {
	relationService service.IRelationService
}

// NewRelationHandler 创建 RelationHandler 实例
func NewRelationHandler(relationService service.IRelationService) *RelationHandler {
	return &RelationHandler{
		relationService: relationService,
	}
}

// RelationAction 关注操作（关注/取消关注）
// @Summary 关注操作
// @Description 用户关注或取消关注其他用户
// @Tags 关注
// @Accept json
// @Produce json
// @Param token query string true "用户token"
// @Param to_user_id query uint true "目标用户ID"
// @Param action_type query int true "操作类型：1-关注，2-取消关注"
// @Success 200 {object} dto.RelationActionResponse
// @Router /douyin/relation/action/ [post]
func (h *RelationHandler) RelationAction(c *gin.Context) {
	// 获取当前用户ID（从JWT中间件中获取）
	userID, exists := c.Get("user_id")
	if !exists {
		global.Logger.Warn("handler.RelationAction.user_id_not_found")
		response.Error(c, errc.ErrUnauthorized, "未授权")
		return
	}

	// 绑定请求参数
	var req dto.RelationActionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		global.Logger.Warn("handler.RelationAction.bind_error",
			zap.Error(err),
		)
		response.Error(c, errc.ErrInvalidParams, "参数错误")
		return
	}

	// 调用 Service 层
	err := h.relationService.RelationAction(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		global.Logger.Error("handler.RelationAction.service_error",
			zap.Uint("user_id", userID.(uint)),
			zap.Uint("to_user_id", req.ToUserID),
			zap.Int32("action_type", req.ActionType),
			zap.Error(err),
		)
		response.Error(c, errc.Failed, err.Error())
		return
	}

	// 返回成功响应
	resp := &dto.RelationActionResponse{
		StatusCode: errc.Success,
		StatusMsg:  "操作成功",
	}

	c.JSON(http.StatusOK, resp)

	global.Logger.Info("handler.RelationAction.success",
		zap.Uint("user_id", userID.(uint)),
		zap.Uint("to_user_id", req.ToUserID),
		zap.Int32("action_type", req.ActionType),
	)
}

// GetFollowList 获取关注列表
// @Summary 获取关注列表
// @Description 获取指定用户的关注列表
// @Tags 关注
// @Accept json
// @Produce json
// @Param token query string true "用户token"
// @Param user_id query uint true "用户ID"
// @Success 200 {object} dto.RelationFollowListResponse
// @Router /douyin/relation/follow/list/ [get]
func (h *RelationHandler) GetFollowList(c *gin.Context) {
	// 获取当前用户ID（可选，用于判断关注状态）
	currentUserID := uint(0)
	if userID, exists := c.Get("user_id"); exists {
		currentUserID = userID.(uint)
	}

	// 绑定请求参数
	var req dto.RelationFollowListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		global.Logger.Warn("handler.GetFollowList.bind_error",
			zap.Error(err),
		)
		response.Error(c, errc.ErrInvalidParams, "参数错误")
		return
	}

	// 调用 Service 层
	userList, err := h.relationService.GetFollowList(c.Request.Context(), currentUserID, req.UserID)
	if err != nil {
		global.Logger.Error("handler.GetFollowList.service_error",
			zap.Uint("user_id", req.UserID),
			zap.Error(err),
		)
		response.Error(c, errc.Failed, err.Error())
		return
	}

	// 返回成功响应
	resp := &dto.RelationFollowListResponse{
		StatusCode: errc.Success,
		StatusMsg:  "success",
		UserList:   userList,
	}

	c.JSON(http.StatusOK, resp)

	global.Logger.Info("handler.GetFollowList.success",
		zap.Uint("user_id", req.UserID),
		zap.Int("count", len(userList)),
	)
}

// GetFollowerList 获取粉丝列表
// @Summary 获取粉丝列表
// @Description 获取指定用户的粉丝列表
// @Tags 关注
// @Accept json
// @Produce json
// @Param token query string true "用户token"
// @Param user_id query uint true "用户ID"
// @Success 200 {object} dto.RelationFollowerListResponse
// @Router /douyin/relation/follower/list/ [get]
func (h *RelationHandler) GetFollowerList(c *gin.Context) {
	// 获取当前用户ID（可选，用于判断关注状态）
	currentUserID := uint(0)
	if userID, exists := c.Get("user_id"); exists {
		currentUserID = userID.(uint)
	}

	// 绑定请求参数
	var req dto.RelationFollowerListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		global.Logger.Warn("handler.GetFollowerList.bind_error",
			zap.Error(err),
		)
		response.Error(c, errc.ErrInvalidParams, "参数错误")
		return
	}

	// 调用 Service 层
	userList, err := h.relationService.GetFollowerList(c.Request.Context(), currentUserID, req.UserID)
	if err != nil {
		global.Logger.Error("handler.GetFollowerList.service_error",
			zap.Uint("user_id", req.UserID),
			zap.Error(err),
		)
		response.Error(c, errc.Failed, err.Error())
		return
	}

	// 返回成功响应
	resp := &dto.RelationFollowerListResponse{
		StatusCode: errc.Success,
		StatusMsg:  "success",
		UserList:   userList,
	}

	c.JSON(http.StatusOK, resp)

	global.Logger.Info("handler.GetFollowerList.success",
		zap.Uint("user_id", req.UserID),
		zap.Int("count", len(userList)),
	)
}

// GetFriendList 获取好友列表
// @Summary 获取好友列表
// @Description 获取当前用户的好友列表（互相关注）
// @Tags 关注
// @Accept json
// @Produce json
// @Param token query string true "用户token"
// @Param user_id query uint true "用户ID"
// @Success 200 {object} dto.RelationFriendListResponse
// @Router /douyin/relation/friend/list/ [get]
func (h *RelationHandler) GetFriendList(c *gin.Context) {
	// 获取当前用户ID（从JWT中间件中获取）
	userID, exists := c.Get("user_id")
	if !exists {
		global.Logger.Warn("handler.GetFriendList.user_id_not_found")
		response.Error(c, errc.ErrUnauthorized, "未授权")
		return
	}

	// 绑定请求参数
	var req dto.RelationFriendListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		global.Logger.Warn("handler.GetFriendList.bind_error",
			zap.Error(err),
		)
		response.Error(c, errc.ErrInvalidParams, "参数错误")
		return
	}

	// 调用 Service 层（查询参数中的 user_id）
	friendList, err := h.relationService.GetFriendList(c.Request.Context(), req.UserID)
	if err != nil {
		global.Logger.Error("handler.GetFriendList.service_error",
			zap.Uint("user_id", userID.(uint)),
			zap.Error(err),
		)
		response.Error(c, errc.Failed, err.Error())
		return
	}

	// 返回成功响应
	resp := &dto.RelationFriendListResponse{
		StatusCode: errc.Success,
		StatusMsg:  "success",
		UserList:   friendList,
	}

	c.JSON(http.StatusOK, resp)

	global.Logger.Info("handler.GetFriendList.success",
		zap.Uint("user_id", userID.(uint)),
		zap.Int("count", len(friendList)),
	)
}
