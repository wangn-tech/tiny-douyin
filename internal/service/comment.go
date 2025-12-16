package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/wangn-tech/tiny-douyin/internal/api/dto"
	"github.com/wangn-tech/tiny-douyin/internal/common/constant"
	"github.com/wangn-tech/tiny-douyin/internal/dao"
	"github.com/wangn-tech/tiny-douyin/internal/global"
	"github.com/wangn-tech/tiny-douyin/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ICommentService 评论服务接口
type ICommentService interface {
	// CommentAction 评论操作（发布/删除）
	CommentAction(ctx context.Context, userID uint, req *dto.CommentActionRequest) (*dto.Comment, error)
	// GetCommentList 获取视频评论列表
	GetCommentList(ctx context.Context, videoID uint) ([]*dto.Comment, error)
}

// CommentService 评论服务实现
type CommentService struct {
	commentDAO dao.ICommentDAO
	videoDAO   dao.IVideoDAO
	userDAO    dao.IUserDAO
	db         *gorm.DB
}

// NewCommentService 创建 CommentService 实例
func NewCommentService(
	commentDAO dao.ICommentDAO,
	videoDAO dao.IVideoDAO,
	userDAO dao.IUserDAO,
	db *gorm.DB,
) ICommentService {
	return &CommentService{
		commentDAO: commentDAO,
		videoDAO:   videoDAO,
		userDAO:    userDAO,
		db:         db,
	}
}

// CommentAction 评论操作（发布/删除）
func (s *CommentService) CommentAction(ctx context.Context, userID uint, req *dto.CommentActionRequest) (*dto.Comment, error) {
	// 验证视频是否存在
	video, err := s.videoDAO.GetVideoByID(ctx, req.VideoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.Warn("service.CommentAction.video_not_found",
				zap.Uint("video_id", req.VideoID),
			)
			return nil, fmt.Errorf("视频不存在")
		}
		global.Logger.Error("service.CommentAction.get_video_error",
			zap.Uint("video_id", req.VideoID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("查询视频失败")
	}

	switch req.ActionType {
	case constant.CommentActionPublish:
		return s.publishComment(ctx, userID, req, video)
	case constant.CommentActionDelete:
		return s.deleteComment(ctx, userID, req, video)
	default:
		global.Logger.Warn("service.CommentAction.invalid_action_type",
			zap.Int32("action_type", req.ActionType),
		)
		return nil, fmt.Errorf("无效的操作类型")
	}
}

// publishComment 发布评论
func (s *CommentService) publishComment(ctx context.Context, userID uint, req *dto.CommentActionRequest, video *model.Video) (*dto.Comment, error) {
	// 验证评论内容
	if req.CommentText == "" {
		global.Logger.Warn("service.publishComment.empty_content",
			zap.Uint("user_id", userID),
		)
		return nil, fmt.Errorf("评论内容不能为空")
	}

	if len(req.CommentText) > constant.CommentMaxLength {
		global.Logger.Warn("service.publishComment.content_too_long",
			zap.Uint("user_id", userID),
			zap.Int("length", len(req.CommentText)),
		)
		return nil, fmt.Errorf("评论内容过长，最多%d个字符", constant.CommentMaxLength)
	}

	// 创建评论对象
	comment := &model.Comment{
		VideoID: req.VideoID,
		UserID:  userID,
		Content: req.CommentText,
	}

	// 使用事务：创建评论 + 增加视频评论数
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 创建评论
		if err := s.commentDAO.CreateComment(ctx, comment); err != nil {
			return err
		}

		// 增加视频评论数
		if err := s.videoDAO.IncrementCommentCount(ctx, req.VideoID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		global.Logger.Error("service.publishComment.transaction_error",
			zap.Uint("user_id", userID),
			zap.Uint("video_id", req.VideoID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("发布评论失败")
	}

	// 查询用户信息
	user, err := s.userDAO.GetUserByID(ctx, userID)
	if err != nil {
		global.Logger.Error("service.publishComment.get_user_error",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("查询用户信息失败")
	}

	// 构建返回的 DTO
	commentDTO := s.buildCommentDTO(comment, user)

	global.Logger.Info("service.publishComment.success",
		zap.Uint("comment_id", comment.ID),
		zap.Uint("user_id", userID),
		zap.Uint("video_id", req.VideoID),
	)

	return commentDTO, nil
}

// deleteComment 删除评论
func (s *CommentService) deleteComment(ctx context.Context, userID uint, req *dto.CommentActionRequest, video *model.Video) (*dto.Comment, error) {
	// 验证 CommentID
	if req.CommentID == 0 {
		global.Logger.Warn("service.deleteComment.missing_comment_id",
			zap.Uint("user_id", userID),
		)
		return nil, fmt.Errorf("评论ID不能为空")
	}

	// 查询评论
	comment, err := s.commentDAO.GetCommentByID(ctx, req.CommentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.Warn("service.deleteComment.comment_not_found",
				zap.Uint("comment_id", req.CommentID),
			)
			return nil, fmt.Errorf("评论不存在")
		}
		global.Logger.Error("service.deleteComment.get_comment_error",
			zap.Uint("comment_id", req.CommentID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("查询评论失败")
	}

	// 验证评论所属视频
	if comment.VideoID != req.VideoID {
		global.Logger.Warn("service.deleteComment.video_mismatch",
			zap.Uint("comment_id", req.CommentID),
			zap.Uint("comment_video_id", comment.VideoID),
			zap.Uint("request_video_id", req.VideoID),
		)
		return nil, fmt.Errorf("评论不属于该视频")
	}

	// 验证删除权限（只能删除自己的评论）
	if comment.UserID != userID {
		global.Logger.Warn("service.deleteComment.permission_denied",
			zap.Uint("comment_id", req.CommentID),
			zap.Uint("comment_user_id", comment.UserID),
			zap.Uint("current_user_id", userID),
		)
		return nil, fmt.Errorf("无权限删除该评论")
	}

	// 使用事务：删除评论 + 减少视频评论数
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 删除评论
		if err := s.commentDAO.DeleteComment(ctx, req.CommentID); err != nil {
			return err
		}

		// 减少视频评论数
		if err := s.videoDAO.DecrementCommentCount(ctx, req.VideoID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		global.Logger.Error("service.deleteComment.transaction_error",
			zap.Uint("user_id", userID),
			zap.Uint("comment_id", req.CommentID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("删除评论失败")
	}

	global.Logger.Info("service.deleteComment.success",
		zap.Uint("comment_id", req.CommentID),
		zap.Uint("user_id", userID),
		zap.Uint("video_id", req.VideoID),
	)

	return nil, nil
}

// GetCommentList 获取视频评论列表
func (s *CommentService) GetCommentList(ctx context.Context, videoID uint) ([]*dto.Comment, error) {
	// 验证视频是否存在
	_, err := s.videoDAO.GetVideoByID(ctx, videoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.Warn("service.GetCommentList.video_not_found",
				zap.Uint("video_id", videoID),
			)
			return nil, fmt.Errorf("视频不存在")
		}
		global.Logger.Error("service.GetCommentList.get_video_error",
			zap.Uint("video_id", videoID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("查询视频失败")
	}

	// 获取评论列表
	comments, err := s.commentDAO.GetVideoComments(ctx, videoID)
	if err != nil {
		global.Logger.Error("service.GetCommentList.get_comments_error",
			zap.Uint("video_id", videoID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("查询评论列表失败")
	}

	// 如果没有评论，返回空列表
	if len(comments) == 0 {
		return []*dto.Comment{}, nil
	}

	// 收集所有评论的用户ID
	userIDs := make([]uint, 0, len(comments))
	for _, comment := range comments {
		userIDs = append(userIDs, comment.UserID)
	}

	// 批量查询用户信息
	users, err := s.userDAO.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		global.Logger.Error("service.GetCommentList.get_users_error",
			zap.Uint("video_id", videoID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("查询用户信息失败")
	}

	// 构建用户信息映射
	userMap := make(map[uint]*model.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	// 构建评论 DTO 列表
	commentDTOs := make([]*dto.Comment, 0, len(comments))
	for _, comment := range comments {
		user := userMap[comment.UserID]
		if user == nil {
			global.Logger.Warn("service.GetCommentList.user_not_found",
				zap.Uint("user_id", comment.UserID),
			)
			continue
		}
		commentDTO := s.buildCommentDTO(comment, user)
		commentDTOs = append(commentDTOs, commentDTO)
	}

	global.Logger.Info("service.GetCommentList.success",
		zap.Uint("video_id", videoID),
		zap.Int("count", len(commentDTOs)),
	)

	return commentDTOs, nil
}

// buildCommentDTO 构建评论 DTO
func (s *CommentService) buildCommentDTO(comment *model.Comment, user *model.User) *dto.Comment {
	return &dto.Comment{
		ID: comment.ID,
		User: &dto.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Avatar:    user.Avatar,
			Signature: user.Signature,
		},
		Content:    comment.Content,
		CreateDate: comment.CreatedAt.Format("01-02"), // MM-DD 格式
	}
}
