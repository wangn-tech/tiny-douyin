package dao

import (
	"context"

	"github.com/wangn-tech/tiny-douyin/internal/global"
	"github.com/wangn-tech/tiny-douyin/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// IUserDAO 用户数据访问接口
type IUserDAO interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	ExistsUsername(ctx context.Context, username string) (bool, error)
	GetUsersByIDs(ctx context.Context, userIDs []uint) ([]*model.User, error)
	// IncrementFollowCount 增加用户关注数
	IncrementFollowCount(ctx context.Context, userID uint) error
	// DecrementFollowCount 减少用户关注数
	DecrementFollowCount(ctx context.Context, userID uint) error
	// IncrementFollowerCount 增加用户粉丝数
	IncrementFollowerCount(ctx context.Context, userID uint) error
	// DecrementFollowerCount 减少用户粉丝数
	DecrementFollowerCount(ctx context.Context, userID uint) error
}

// UserDAO 用户数据访问实现
type UserDAO struct {
	db *gorm.DB
}

// NewUserDAO 创建 UserDAO 实例
func NewUserDAO(db *gorm.DB) IUserDAO {
	return &UserDAO{
		db: db,
	}
}

// CreateUser 创建用户
func (d *UserDAO) CreateUser(ctx context.Context, user *model.User) error {
	err := d.db.WithContext(ctx).Create(user).Error

	if err != nil {
		global.Logger.Error("dao.CreateUser.failed",
			zap.String("username", user.Username),
			zap.Error(err),
		)
	}

	return err
}

// ExistsUsername 检查用户名是否存在（高性能查询，只查 ID）
func (d *UserDAO) ExistsUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := d.db.WithContext(ctx).
		Model(&model.User{}).
		Where("username = ?", username).
		Count(&count).Error

	if err != nil {
		global.Logger.Error("dao.ExistsUsername.db_error",
			zap.String("username", username),
			zap.Error(err),
		)
		return false, err
	}

	return count > 0, nil
}

// GetUserByUsername 根据用户名获取用户（完整信息）
func (d *UserDAO) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User

	err := d.db.WithContext(ctx).
		Where("username = ?", username).
		First(&user).Error

	// 只记录真正的数据库错误
	if err != nil && err != gorm.ErrRecordNotFound {
		global.Logger.Error("dao.GetUserByUsername.db_error",
			zap.String("username", username),
			zap.Error(err),
		)
	}

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据用户ID获取用户
func (d *UserDAO) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User

	err := d.db.WithContext(ctx).First(&user, id).Error

	// 只记录真正的数据库错误
	if err != nil && err != gorm.ErrRecordNotFound {
		global.Logger.Error("dao.GetUserByID.db_error",
			zap.Uint("user_id", id),
			zap.Error(err),
		)
	}

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUsersByIDs 批量查询用户信息
func (d *UserDAO) GetUsersByIDs(ctx context.Context, userIDs []uint) ([]*model.User, error) {
	if len(userIDs) == 0 {
		return []*model.User{}, nil
	}

	var users []*model.User
	err := d.db.WithContext(ctx).
		Where("id IN ?", userIDs).
		Find(&users).Error

	if err != nil {
		global.Logger.Error("dao.GetUsersByIDs.db_error",
			zap.Any("user_ids", userIDs),
			zap.Error(err),
		)
		return nil, err
	}

	global.Logger.Info("dao.GetUsersByIDs.success",
		zap.Int("count", len(users)),
	)

	return users, nil
}

// IncrementFollowCount 增加用户关注数
func (d *UserDAO) IncrementFollowCount(ctx context.Context, userID uint) error {
	err := d.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		UpdateColumn("follow_count", gorm.Expr("follow_count + ?", 1)).Error

	if err != nil {
		global.Logger.Error("dao.IncrementFollowCount.db_error",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return err
	}

	global.Logger.Info("dao.IncrementFollowCount.success",
		zap.Uint("user_id", userID),
	)

	return nil
}

// DecrementFollowCount 减少用户关注数
func (d *UserDAO) DecrementFollowCount(ctx context.Context, userID uint) error {
	err := d.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ? AND follow_count > 0", userID).
		UpdateColumn("follow_count", gorm.Expr("follow_count - ?", 1)).Error

	if err != nil {
		global.Logger.Error("dao.DecrementFollowCount.db_error",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return err
	}

	global.Logger.Info("dao.DecrementFollowCount.success",
		zap.Uint("user_id", userID),
	)

	return nil
}

// IncrementFollowerCount 增加用户粉丝数
func (d *UserDAO) IncrementFollowerCount(ctx context.Context, userID uint) error {
	err := d.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		UpdateColumn("follower_count", gorm.Expr("follower_count + ?", 1)).Error

	if err != nil {
		global.Logger.Error("dao.IncrementFollowerCount.db_error",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return err
	}

	global.Logger.Info("dao.IncrementFollowerCount.success",
		zap.Uint("user_id", userID),
	)

	return nil
}

// DecrementFollowerCount 减少用户粉丝数
func (d *UserDAO) DecrementFollowerCount(ctx context.Context, userID uint) error {
	err := d.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ? AND follower_count > 0", userID).
		UpdateColumn("follower_count", gorm.Expr("follower_count - ?", 1)).Error

	if err != nil {
		global.Logger.Error("dao.DecrementFollowerCount.db_error",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return err
	}

	global.Logger.Info("dao.DecrementFollowerCount.success",
		zap.Uint("user_id", userID),
	)

	return nil
}
