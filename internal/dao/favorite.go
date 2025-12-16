package dao

import (
	"context"

	"github.com/wangn-tech/tiny-douyin/internal/global"
	"github.com/wangn-tech/tiny-douyin/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// IFavoriteDAO 点赞数据访问接口
type IFavoriteDAO interface {
	// CreateFavorite 创建点赞记录
	CreateFavorite(ctx context.Context, userID, videoID uint) error
	// DeleteFavorite 删除点赞记录
	DeleteFavorite(ctx context.Context, userID, videoID uint) error
	// IsFavorite 检查用户是否点赞了某个视频
	IsFavorite(ctx context.Context, userID, videoID uint) (bool, error)
	// GetUserFavoriteVideoIDs 获取用户点赞的所有视频ID列表
	GetUserFavoriteVideoIDs(ctx context.Context, userID uint) ([]uint, error)
	// GetFavoriteCount 获取视频的点赞数
	GetFavoriteCount(ctx context.Context, videoID uint) (int64, error)
	// BatchCheckFavorite 批量检查用户是否点赞了视频列表（返回 map[videoID]bool）
	BatchCheckFavorite(ctx context.Context, userID uint, videoIDs []uint) (map[uint]bool, error)
}

// FavoriteDAO 点赞数据访问实现
type FavoriteDAO struct {
	db *gorm.DB
}

// NewFavoriteDAO 创建 FavoriteDAO 实例
func NewFavoriteDAO(db *gorm.DB) IFavoriteDAO {
	return &FavoriteDAO{db: db}
}

// CreateFavorite 创建点赞记录
func (d *FavoriteDAO) CreateFavorite(ctx context.Context, userID, videoID uint) error {
	favorite := &model.Favorite{
		UserID:  userID,
		VideoID: videoID,
	}

	err := d.db.WithContext(ctx).Create(favorite).Error
	if err != nil {
		global.Logger.Error("dao.CreateFavorite.db_error",
			zap.Uint("user_id", userID),
			zap.Uint("video_id", videoID),
			zap.Error(err),
		)
		return err
	}

	global.Logger.Info("dao.CreateFavorite.success",
		zap.Uint("user_id", userID),
		zap.Uint("video_id", videoID),
	)

	return nil
}

// DeleteFavorite 删除点赞记录（软删除）
func (d *FavoriteDAO) DeleteFavorite(ctx context.Context, userID, videoID uint) error {
	result := d.db.WithContext(ctx).
		Where("user_id = ? AND video_id = ?", userID, videoID).
		Delete(&model.Favorite{})

	if result.Error != nil {
		global.Logger.Error("dao.DeleteFavorite.db_error",
			zap.Uint("user_id", userID),
			zap.Uint("video_id", videoID),
			zap.Error(result.Error),
		)
		return result.Error
	}

	global.Logger.Info("dao.DeleteFavorite.success",
		zap.Uint("user_id", userID),
		zap.Uint("video_id", videoID),
		zap.Int64("rows_affected", result.RowsAffected),
	)

	return nil
}

// IsFavorite 检查用户是否点赞了某个视频
func (d *FavoriteDAO) IsFavorite(ctx context.Context, userID, videoID uint) (bool, error) {
	var count int64
	err := d.db.WithContext(ctx).
		Model(&model.Favorite{}).
		Where("user_id = ? AND video_id = ?", userID, videoID).
		Count(&count).Error

	if err != nil {
		global.Logger.Error("dao.IsFavorite.db_error",
			zap.Uint("user_id", userID),
			zap.Uint("video_id", videoID),
			zap.Error(err),
		)
		return false, err
	}

	return count > 0, nil
}

// GetUserFavoriteVideoIDs 获取用户点赞的所有视频ID列表
func (d *FavoriteDAO) GetUserFavoriteVideoIDs(ctx context.Context, userID uint) ([]uint, error) {
	var favorites []model.Favorite
	err := d.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&favorites).Error

	if err != nil {
		global.Logger.Error("dao.GetUserFavoriteVideoIDs.db_error",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return nil, err
	}

	videoIDs := make([]uint, 0, len(favorites))
	for _, favorite := range favorites {
		videoIDs = append(videoIDs, favorite.VideoID)
	}

	global.Logger.Info("dao.GetUserFavoriteVideoIDs.success",
		zap.Uint("user_id", userID),
		zap.Int("count", len(videoIDs)),
	)

	return videoIDs, nil
}

// GetFavoriteCount 获取视频的点赞数
func (d *FavoriteDAO) GetFavoriteCount(ctx context.Context, videoID uint) (int64, error) {
	var count int64
	err := d.db.WithContext(ctx).
		Model(&model.Favorite{}).
		Where("video_id = ?", videoID).
		Count(&count).Error

	if err != nil {
		global.Logger.Error("dao.GetFavoriteCount.db_error",
			zap.Uint("video_id", videoID),
			zap.Error(err),
		)
		return 0, err
	}

	return count, nil
}

// BatchCheckFavorite 批量检查用户是否点赞了视频列表（优化性能）
func (d *FavoriteDAO) BatchCheckFavorite(ctx context.Context, userID uint, videoIDs []uint) (map[uint]bool, error) {
	if len(videoIDs) == 0 {
		return make(map[uint]bool), nil
	}

	var favorites []model.Favorite
	err := d.db.WithContext(ctx).
		Where("user_id = ? AND video_id IN ?", userID, videoIDs).
		Find(&favorites).Error

	if err != nil {
		global.Logger.Error("dao.BatchCheckFavorite.db_error",
			zap.Uint("user_id", userID),
			zap.Int("video_count", len(videoIDs)),
			zap.Error(err),
		)
		return nil, err
	}

	// 构建 map
	favoriteMap := make(map[uint]bool)
	for _, favorite := range favorites {
		favoriteMap[favorite.VideoID] = true
	}

	global.Logger.Info("dao.BatchCheckFavorite.success",
		zap.Uint("user_id", userID),
		zap.Int("total_videos", len(videoIDs)),
		zap.Int("favorited_videos", len(favoriteMap)),
	)

	return favoriteMap, nil
}
