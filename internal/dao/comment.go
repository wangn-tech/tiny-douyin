package dao

import (
	"context"

	"github.com/wangn-tech/tiny-douyin/internal/global"
	"github.com/wangn-tech/tiny-douyin/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ICommentDAO 评论数据访问接口
type ICommentDAO interface {
	// CreateComment 创建评论
	CreateComment(ctx context.Context, comment *model.Comment) error
	// DeleteComment 删除评论（软删除）
	DeleteComment(ctx context.Context, commentID uint) error
	// GetCommentByID 根据ID查询评论
	GetCommentByID(ctx context.Context, commentID uint) (*model.Comment, error)
	// GetVideoComments 获取视频的所有评论（按时间倒序）
	GetVideoComments(ctx context.Context, videoID uint) ([]*model.Comment, error)
	// GetCommentCount 获取视频的评论数
	GetCommentCount(ctx context.Context, videoID uint) (int64, error)
}

// CommentDAO 评论数据访问实现
type CommentDAO struct {
	db *gorm.DB
}

// NewCommentDAO 创建 CommentDAO 实例
func NewCommentDAO(db *gorm.DB) ICommentDAO {
	return &CommentDAO{db: db}
}

// CreateComment 创建评论
func (d *CommentDAO) CreateComment(ctx context.Context, comment *model.Comment) error {
	err := d.db.WithContext(ctx).Create(comment).Error
	if err != nil {
		global.Logger.Error("dao.CreateComment.db_error",
			zap.Uint("user_id", comment.UserID),
			zap.Uint("video_id", comment.VideoID),
			zap.Error(err),
		)
		return err
	}

	global.Logger.Info("dao.CreateComment.success",
		zap.Uint("comment_id", comment.ID),
		zap.Uint("user_id", comment.UserID),
		zap.Uint("video_id", comment.VideoID),
	)

	return nil
}

// DeleteComment 删除评论（软删除）
func (d *CommentDAO) DeleteComment(ctx context.Context, commentID uint) error {
	result := d.db.WithContext(ctx).Delete(&model.Comment{}, commentID)

	if result.Error != nil {
		global.Logger.Error("dao.DeleteComment.db_error",
			zap.Uint("comment_id", commentID),
			zap.Error(result.Error),
		)
		return result.Error
	}

	global.Logger.Info("dao.DeleteComment.success",
		zap.Uint("comment_id", commentID),
		zap.Int64("rows_affected", result.RowsAffected),
	)

	return nil
}

// GetCommentByID 根据ID查询评论
func (d *CommentDAO) GetCommentByID(ctx context.Context, commentID uint) (*model.Comment, error) {
	var comment model.Comment
	err := d.db.WithContext(ctx).First(&comment, commentID).Error
	if err != nil {
		global.Logger.Error("dao.GetCommentByID.db_error",
			zap.Uint("comment_id", commentID),
			zap.Error(err),
		)
		return nil, err
	}

	return &comment, nil
}

// GetVideoComments 获取视频的所有评论（按时间倒序）
func (d *CommentDAO) GetVideoComments(ctx context.Context, videoID uint) ([]*model.Comment, error) {
	var comments []*model.Comment
	err := d.db.WithContext(ctx).
		Where("video_id = ?", videoID).
		Order("created_at DESC").
		Find(&comments).Error

	if err != nil {
		global.Logger.Error("dao.GetVideoComments.db_error",
			zap.Uint("video_id", videoID),
			zap.Error(err),
		)
		return nil, err
	}

	global.Logger.Info("dao.GetVideoComments.success",
		zap.Uint("video_id", videoID),
		zap.Int("count", len(comments)),
	)

	return comments, nil
}

// GetCommentCount 获取视频的评论数
func (d *CommentDAO) GetCommentCount(ctx context.Context, videoID uint) (int64, error) {
	var count int64
	err := d.db.WithContext(ctx).
		Model(&model.Comment{}).
		Where("video_id = ?", videoID).
		Count(&count).Error

	if err != nil {
		global.Logger.Error("dao.GetCommentCount.db_error",
			zap.Uint("video_id", videoID),
			zap.Error(err),
		)
		return 0, err
	}

	return count, nil
}
