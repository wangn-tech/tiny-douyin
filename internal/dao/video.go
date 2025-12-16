package dao

import (
	"context"

	"github.com/wangn-tech/tiny-douyin/internal/global"
	"github.com/wangn-tech/tiny-douyin/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// IVideoDAO 视频数据访问接口
type IVideoDAO interface {
	// CreateVideo 创建视频
	CreateVideo(ctx context.Context, video *model.Video) error
	// GetVideoByID 根据ID查询视频
	GetVideoByID(ctx context.Context, id uint) (*model.Video, error)
	// GetVideosByUserID 根据用户ID查询视频列表
	GetVideosByUserID(ctx context.Context, userID uint) ([]*model.Video, error)
	// GetVideoFeed 获取视频流（按时间倒序）
	GetVideoFeed(ctx context.Context, latestTime int64, limit int) ([]*model.Video, error)
	// UpdateVideo 更新视频信息
	UpdateVideo(ctx context.Context, video *model.Video) error
	// GetVideosByIDs 批量查询视频（用于喜欢列表）
	GetVideosByIDs(ctx context.Context, videoIDs []uint) ([]*model.Video, error)
	// IncrementFavoriteCount 增加视频点赞数
	IncrementFavoriteCount(ctx context.Context, videoID uint) error
	// DecrementFavoriteCount 减少视频点赞数
	DecrementFavoriteCount(ctx context.Context, videoID uint) error
	// IncrementCommentCount 增加视频评论数
	IncrementCommentCount(ctx context.Context, videoID uint) error
	// DecrementCommentCount 减少视频评论数
	DecrementCommentCount(ctx context.Context, videoID uint) error
}

// VideoDAO 视频数据访问实现
type VideoDAO struct {
	db *gorm.DB
}

// NewVideoDAO 创建 VideoDAO 实例
func NewVideoDAO(db *gorm.DB) IVideoDAO {
	return &VideoDAO{db: db}
}

// CreateVideo 创建视频
func (d *VideoDAO) CreateVideo(ctx context.Context, video *model.Video) error {
	err := d.db.WithContext(ctx).Create(video).Error
	if err != nil {
		global.Logger.Error("dao.CreateVideo.db_error",
			zap.Uint("author_id", video.AuthorID),
			zap.String("title", video.Title),
			zap.Error(err),
		)
		return err
	}
	return nil
}

// GetVideoByID 根据ID查询视频
func (d *VideoDAO) GetVideoByID(ctx context.Context, id uint) (*model.Video, error) {
	var video model.Video
	err := d.db.WithContext(ctx).First(&video, id).Error
	if err != nil {
		global.Logger.Error("dao.GetVideoByID.db_error",
			zap.Uint("video_id", id),
			zap.Error(err),
		)
		return nil, err
	}
	return &video, nil
}

// GetVideosByUserID 根据用户ID查询视频列表
func (d *VideoDAO) GetVideosByUserID(ctx context.Context, userID uint) ([]*model.Video, error) {
	var videos []*model.Video
	err := d.db.WithContext(ctx).
		Where("author_id = ?", userID).
		Order("created_at DESC").
		Find(&videos).Error
	if err != nil {
		global.Logger.Error("dao.GetVideosByUserID.db_error",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return nil, err
	}
	return videos, nil
}

// GetVideoFeed 获取视频流（按时间倒序）
// latestTime: Unix时间戳（秒），返回比该时间更早的视频
// limit: 限制返回数量
func (d *VideoDAO) GetVideoFeed(ctx context.Context, latestTime int64, limit int) ([]*model.Video, error) {
	var videos []*model.Video
	query := d.db.WithContext(ctx).Order("created_at DESC")

	// 如果提供了 latestTime，则只返回比该时间更早的视频
	if latestTime > 0 {
		query = query.Where("UNIX_TIMESTAMP(created_at) < ?", latestTime)
	}

	err := query.Limit(limit).Find(&videos).Error
	if err != nil {
		global.Logger.Error("dao.GetVideoFeed.db_error",
			zap.Int64("latest_time", latestTime),
			zap.Int("limit", limit),
			zap.Error(err),
		)
		return nil, err
	}
	return videos, nil
}

// UpdateVideo 更新视频信息
func (d *VideoDAO) UpdateVideo(ctx context.Context, video *model.Video) error {
	err := d.db.WithContext(ctx).Save(video).Error
	if err != nil {
		global.Logger.Error("dao.UpdateVideo.db_error",
			zap.Uint("video_id", video.ID),
			zap.Error(err),
		)
		return err
	}
	return nil
}

// GetVideosByIDs 批量查询视频（用于喜欢列表）
func (d *VideoDAO) GetVideosByIDs(ctx context.Context, videoIDs []uint) ([]*model.Video, error) {
	if len(videoIDs) == 0 {
		return []*model.Video{}, nil
	}

	var videos []*model.Video
	err := d.db.WithContext(ctx).
		Where("id IN ?", videoIDs).
		Find(&videos).Error

	if err != nil {
		global.Logger.Error("dao.GetVideosByIDs.db_error",
			zap.Int("video_count", len(videoIDs)),
			zap.Error(err),
		)
		return nil, err
	}

	global.Logger.Info("dao.GetVideosByIDs.success",
		zap.Int("request_count", len(videoIDs)),
		zap.Int("found_count", len(videos)),
	)

	return videos, nil
}

// IncrementFavoriteCount 增加视频点赞数
func (d *VideoDAO) IncrementFavoriteCount(ctx context.Context, videoID uint) error {
	err := d.db.WithContext(ctx).
		Model(&model.Video{}).
		Where("id = ?", videoID).
		UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error

	if err != nil {
		global.Logger.Error("dao.IncrementFavoriteCount.db_error",
			zap.Uint("video_id", videoID),
			zap.Error(err),
		)
		return err
	}

	global.Logger.Info("dao.IncrementFavoriteCount.success",
		zap.Uint("video_id", videoID),
	)

	return nil
}

// DecrementFavoriteCount 减少视频点赞数
func (d *VideoDAO) DecrementFavoriteCount(ctx context.Context, videoID uint) error {
	err := d.db.WithContext(ctx).
		Model(&model.Video{}).
		Where("id = ? AND favorite_count > 0", videoID).
		UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error

	if err != nil {
		global.Logger.Error("dao.DecrementFavoriteCount.db_error",
			zap.Uint("video_id", videoID),
			zap.Error(err),
		)
		return err
	}

	global.Logger.Info("dao.DecrementFavoriteCount.success",
		zap.Uint("video_id", videoID),
	)

	return nil
}

// IncrementCommentCount 增加视频评论数
func (d *VideoDAO) IncrementCommentCount(ctx context.Context, videoID uint) error {
	err := d.db.WithContext(ctx).
		Model(&model.Video{}).
		Where("id = ?", videoID).
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error

	if err != nil {
		global.Logger.Error("dao.IncrementCommentCount.db_error",
			zap.Uint("video_id", videoID),
			zap.Error(err),
		)
		return err
	}

	global.Logger.Info("dao.IncrementCommentCount.success",
		zap.Uint("video_id", videoID),
	)

	return nil
}

// DecrementCommentCount 减少视频评论数
func (d *VideoDAO) DecrementCommentCount(ctx context.Context, videoID uint) error {
	err := d.db.WithContext(ctx).
		Model(&model.Video{}).
		Where("id = ? AND comment_count > 0", videoID).
		UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1)).Error

	if err != nil {
		global.Logger.Error("dao.DecrementCommentCount.db_error",
			zap.Uint("video_id", videoID),
			zap.Error(err),
		)
		return err
	}

	global.Logger.Info("dao.DecrementCommentCount.success",
		zap.Uint("video_id", videoID),
	)

	return nil
}
