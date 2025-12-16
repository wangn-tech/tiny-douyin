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

// CommentHandler 评论处理器
type CommentHandler struct {
	commentService service.ICommentService
}

// NewCommentHandler 创建 CommentHandler 实例
func NewCommentHandler(commentService service.ICommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// CommentAction 评论操作（发布/删除）
// @Summary 评论操作
// @Description 用户对视频进行评论或删除评论
// @Tags 评论
// @Accept json
// @Produce json
// @Param token query string true "用户token"
// @Param video_id query uint true "视频ID"
// @Param action_type query int true "操作类型：1-发布评论，2-删除评论"
// @Param comment_text query string false "评论内容（发布评论时使用）"
// @Param comment_id query uint false "评论ID（删除评论时使用）"
// @Success 200 {object} dto.CommentActionResponse
// @Router /douyin/comment/action/ [post]
func (h *CommentHandler) CommentAction(c *gin.Context) {
	// 获取当前用户ID（从JWT中间件中获取）
	userID, exists := c.Get("user_id")
	if !exists {
		global.Logger.Warn("handler.CommentAction.user_id_not_found")
		response.Error(c, errc.ErrUnauthorized, "未授权")
		return
	}

	// 绑定请求参数
	var req dto.CommentActionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		global.Logger.Warn("handler.CommentAction.bind_error",
			zap.Error(err),
		)
		response.Error(c, errc.ErrInvalidParams, "参数错误")
		return
	}

	// 调用 Service 层
	comment, err := h.commentService.CommentAction(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		global.Logger.Error("handler.CommentAction.service_error",
			zap.Uint("user_id", userID.(uint)),
			zap.Uint("video_id", req.VideoID),
			zap.Int32("action_type", req.ActionType),
			zap.Error(err),
		)
		response.Error(c, errc.Failed, err.Error())
		return
	}

	// 返回成功响应
	resp := &dto.CommentActionResponse{
		StatusCode: errc.Success,
		StatusMsg:  "操作成功",
		Comment:    comment,
	}

	c.JSON(http.StatusOK, resp)

	global.Logger.Info("handler.CommentAction.success",
		zap.Uint("user_id", userID.(uint)),
		zap.Uint("video_id", req.VideoID),
		zap.Int32("action_type", req.ActionType),
	)
}

// GetCommentList 获取视频评论列表
// @Summary 获取视频评论列表
// @Description 获取指定视频的所有评论
// @Tags 评论
// @Accept json
// @Produce json
// @Param token query string true "用户token"
// @Param video_id query uint true "视频ID"
// @Success 200 {object} dto.CommentListResponse
// @Router /douyin/comment/list/ [get]
func (h *CommentHandler) GetCommentList(c *gin.Context) {
	// 绑定请求参数
	var req dto.CommentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		global.Logger.Warn("handler.GetCommentList.bind_error",
			zap.Error(err),
		)
		response.Error(c, errc.ErrInvalidParams, "参数错误")
		return
	}

	// 调用 Service 层
	comments, err := h.commentService.GetCommentList(c.Request.Context(), req.VideoID)
	if err != nil {
		global.Logger.Error("handler.GetCommentList.service_error",
			zap.Uint("video_id", req.VideoID),
			zap.Error(err),
		)
		response.Error(c, errc.Failed, err.Error())
		return
	}

	// 返回成功响应
	resp := &dto.CommentListResponse{
		StatusCode:  errc.Success,
		StatusMsg:   "success",
		CommentList: comments,
	}

	c.JSON(http.StatusOK, resp)

	global.Logger.Info("handler.GetCommentList.success",
		zap.Uint("video_id", req.VideoID),
		zap.Int("count", len(comments)),
	)
}
