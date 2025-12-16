package dao

import (
	"context"

	"github.com/wangn-tech/tiny-douyin/internal/global"
	"github.com/wangn-tech/tiny-douyin/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// IRelationDAO 关注关系数据访问接口
type IRelationDAO interface {
	// CreateRelation 创建关注关系
	CreateRelation(ctx context.Context, relation *model.Relation) error
	// DeleteRelation 删除关注关系（软删除）
	DeleteRelation(ctx context.Context, followerID, followeeID uint) error
	// IsFollowing 判断是否已关注
	IsFollowing(ctx context.Context, followerID, followeeID uint) (bool, error)
	// GetFollowList 获取关注列表（某用户关注的所有用户ID）
	GetFollowList(ctx context.Context, userID uint) ([]uint, error)
	// GetFollowerList 获取粉丝列表（关注某用户的所有用户ID）
	GetFollowerList(ctx context.Context, userID uint) ([]uint, error)
	// GetFriendList 获取好友列表（互相关注的用户ID）
	GetFriendList(ctx context.Context, userID uint) ([]uint, error)
	// BatchCheckFollowing 批量检查关注状态
	BatchCheckFollowing(ctx context.Context, followerID uint, followeeIDs []uint) (map[uint]bool, error)
}

// RelationDAO 关注关系数据访问实现
type RelationDAO struct {
	db *gorm.DB
}

// NewRelationDAO 创建 RelationDAO 实例
func NewRelationDAO(db *gorm.DB) IRelationDAO {
	return &RelationDAO{db: db}
}

// CreateRelation 创建关注关系
func (d *RelationDAO) CreateRelation(ctx context.Context, relation *model.Relation) error {
	err := d.db.WithContext(ctx).Create(relation).Error
	if err != nil {
		global.Logger.Error("dao.CreateRelation.db_error",
			zap.Uint("follower_id", relation.FollowerID),
			zap.Uint("followee_id", relation.FolloweeID),
			zap.Error(err),
		)
		return err
	}

	global.Logger.Info("dao.CreateRelation.success",
		zap.Uint("relation_id", relation.ID),
		zap.Uint("follower_id", relation.FollowerID),
		zap.Uint("followee_id", relation.FolloweeID),
	)

	return nil
}

// DeleteRelation 删除关注关系（软删除）
func (d *RelationDAO) DeleteRelation(ctx context.Context, followerID, followeeID uint) error {
	result := d.db.WithContext(ctx).
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Delete(&model.Relation{})

	if result.Error != nil {
		global.Logger.Error("dao.DeleteRelation.db_error",
			zap.Uint("follower_id", followerID),
			zap.Uint("followee_id", followeeID),
			zap.Error(result.Error),
		)
		return result.Error
	}

	global.Logger.Info("dao.DeleteRelation.success",
		zap.Uint("follower_id", followerID),
		zap.Uint("followee_id", followeeID),
		zap.Int64("rows_affected", result.RowsAffected),
	)

	return nil
}

// IsFollowing 判断是否已关注
func (d *RelationDAO) IsFollowing(ctx context.Context, followerID, followeeID uint) (bool, error) {
	var count int64
	err := d.db.WithContext(ctx).
		Model(&model.Relation{}).
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Count(&count).Error

	if err != nil {
		global.Logger.Error("dao.IsFollowing.db_error",
			zap.Uint("follower_id", followerID),
			zap.Uint("followee_id", followeeID),
			zap.Error(err),
		)
		return false, err
	}

	return count > 0, nil
}

// GetFollowList 获取关注列表（某用户关注的所有用户ID）
func (d *RelationDAO) GetFollowList(ctx context.Context, userID uint) ([]uint, error) {
	var relations []*model.Relation
	err := d.db.WithContext(ctx).
		Select("followee_id").
		Where("follower_id = ?", userID).
		Find(&relations).Error

	if err != nil {
		global.Logger.Error("dao.GetFollowList.db_error",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return nil, err
	}

	followeeIDs := make([]uint, 0, len(relations))
	for _, rel := range relations {
		followeeIDs = append(followeeIDs, rel.FolloweeID)
	}

	global.Logger.Info("dao.GetFollowList.success",
		zap.Uint("user_id", userID),
		zap.Int("count", len(followeeIDs)),
	)

	return followeeIDs, nil
}

// GetFollowerList 获取粉丝列表（关注某用户的所有用户ID）
func (d *RelationDAO) GetFollowerList(ctx context.Context, userID uint) ([]uint, error) {
	var relations []*model.Relation
	err := d.db.WithContext(ctx).
		Select("follower_id").
		Where("followee_id = ?", userID).
		Find(&relations).Error

	if err != nil {
		global.Logger.Error("dao.GetFollowerList.db_error",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return nil, err
	}

	followerIDs := make([]uint, 0, len(relations))
	for _, rel := range relations {
		followerIDs = append(followerIDs, rel.FollowerID)
	}

	global.Logger.Info("dao.GetFollowerList.success",
		zap.Uint("user_id", userID),
		zap.Int("count", len(followerIDs)),
	)

	return followerIDs, nil
}

// GetFriendList 获取好友列表（互相关注的用户ID）
func (d *RelationDAO) GetFriendList(ctx context.Context, userID uint) ([]uint, error) {
	// 查询当前用户关注的人
	followList, err := d.GetFollowList(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(followList) == 0 {
		return []uint{}, nil
	}

	// 查询这些人中谁也关注了当前用户（互相关注）
	var relations []*model.Relation
	err = d.db.WithContext(ctx).
		Select("follower_id").
		Where("follower_id IN ? AND followee_id = ?", followList, userID).
		Find(&relations).Error

	if err != nil {
		global.Logger.Error("dao.GetFriendList.db_error",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return nil, err
	}

	friendIDs := make([]uint, 0, len(relations))
	for _, rel := range relations {
		friendIDs = append(friendIDs, rel.FollowerID)
	}

	global.Logger.Info("dao.GetFriendList.success",
		zap.Uint("user_id", userID),
		zap.Int("count", len(friendIDs)),
	)

	return friendIDs, nil
}

// BatchCheckFollowing 批量检查关注状态
func (d *RelationDAO) BatchCheckFollowing(ctx context.Context, followerID uint, followeeIDs []uint) (map[uint]bool, error) {
	if len(followeeIDs) == 0 {
		return make(map[uint]bool), nil
	}

	var relations []*model.Relation
	err := d.db.WithContext(ctx).
		Select("followee_id").
		Where("follower_id = ? AND followee_id IN ?", followerID, followeeIDs).
		Find(&relations).Error

	if err != nil {
		global.Logger.Error("dao.BatchCheckFollowing.db_error",
			zap.Uint("follower_id", followerID),
			zap.Any("followee_ids", followeeIDs),
			zap.Error(err),
		)
		return nil, err
	}

	// 构建关注状态映射
	followMap := make(map[uint]bool)
	for _, rel := range relations {
		followMap[rel.FolloweeID] = true
	}

	global.Logger.Info("dao.BatchCheckFollowing.success",
		zap.Uint("follower_id", followerID),
		zap.Int("total_count", len(followeeIDs)),
		zap.Int("following_count", len(followMap)),
	)

	return followMap, nil
}
