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

// IRelationService 关注服务接口
type IRelationService interface {
	// RelationAction 关注操作（关注/取消关注）
	RelationAction(ctx context.Context, userID uint, req *dto.RelationActionRequest) error
	// GetFollowList 获取关注列表
	GetFollowList(ctx context.Context, currentUserID, targetUserID uint) ([]*dto.UserInfo, error)
	// GetFollowerList 获取粉丝列表
	GetFollowerList(ctx context.Context, currentUserID, targetUserID uint) ([]*dto.UserInfo, error)
	// GetFriendList 获取好友列表
	GetFriendList(ctx context.Context, userID uint) ([]*dto.FriendInfo, error)
	// IsFriend 判断两个用户是否为好友（双向关注）
	IsFriend(ctx context.Context, userID1, userID2 uint) (bool, error)
}

// RelationService 关注服务实现
type RelationService struct {
	relationDAO dao.IRelationDAO
	userDAO     dao.IUserDAO
	db          *gorm.DB
}

// NewRelationService 创建 RelationService 实例
func NewRelationService(
	relationDAO dao.IRelationDAO,
	userDAO dao.IUserDAO,
	db *gorm.DB,
) IRelationService {
	return &RelationService{
		relationDAO: relationDAO,
		userDAO:     userDAO,
		db:          db,
	}
}

// RelationAction 关注操作（关注/取消关注）
func (s *RelationService) RelationAction(ctx context.Context, userID uint, req *dto.RelationActionRequest) error {
	// 验证不能关注自己
	if userID == req.ToUserID {
		global.Logger.Warn("service.RelationAction.cannot_follow_self",
			zap.Uint("user_id", userID),
		)
		return fmt.Errorf("不能关注自己")
	}

	// 验证目标用户是否存在
	_, err := s.userDAO.GetUserByID(ctx, req.ToUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.Warn("service.RelationAction.target_user_not_found",
				zap.Uint("to_user_id", req.ToUserID),
			)
			return fmt.Errorf("目标用户不存在")
		}
		global.Logger.Error("service.RelationAction.get_user_error",
			zap.Uint("to_user_id", req.ToUserID),
			zap.Error(err),
		)
		return fmt.Errorf("查询用户失败")
	}

	switch req.ActionType {
	case constant.RelationActionFollow:
		return s.followUser(ctx, userID, req.ToUserID)
	case constant.RelationActionUnfollow:
		return s.unfollowUser(ctx, userID, req.ToUserID)
	default:
		global.Logger.Warn("service.RelationAction.invalid_action_type",
			zap.Int32("action_type", req.ActionType),
		)
		return fmt.Errorf("无效的操作类型")
	}
}

// followUser 关注用户
func (s *RelationService) followUser(ctx context.Context, followerID, followeeID uint) error {
	// 检查是否已关注（幂等性）
	isFollowing, err := s.relationDAO.IsFollowing(ctx, followerID, followeeID)
	if err != nil {
		global.Logger.Error("service.followUser.check_following_error",
			zap.Uint("follower_id", followerID),
			zap.Uint("followee_id", followeeID),
			zap.Error(err),
		)
		return fmt.Errorf("检查关注状态失败")
	}

	if isFollowing {
		// 已关注，直接返回成功（幂等性）
		global.Logger.Info("service.followUser.already_followed",
			zap.Uint("follower_id", followerID),
			zap.Uint("followee_id", followeeID),
		)
		return nil
	}

	// 创建关注关系
	relation := &model.Relation{
		FollowerID: followerID,
		FolloweeID: followeeID,
	}

	// 使用事务：创建关注关系 + 更新双方统计数
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 创建关注记录
		if err := s.relationDAO.CreateRelation(ctx, relation); err != nil {
			return err
		}

		// 增加关注者的关注数
		if err := s.userDAO.IncrementFollowCount(ctx, followerID); err != nil {
			return err
		}

		// 增加被关注者的粉丝数
		if err := s.userDAO.IncrementFollowerCount(ctx, followeeID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		global.Logger.Error("service.followUser.transaction_error",
			zap.Uint("follower_id", followerID),
			zap.Uint("followee_id", followeeID),
			zap.Error(err),
		)
		return fmt.Errorf("关注失败")
	}

	global.Logger.Info("service.followUser.success",
		zap.Uint("follower_id", followerID),
		zap.Uint("followee_id", followeeID),
	)

	return nil
}

// unfollowUser 取消关注用户
func (s *RelationService) unfollowUser(ctx context.Context, followerID, followeeID uint) error {
	// 检查是否已关注（幂等性）
	isFollowing, err := s.relationDAO.IsFollowing(ctx, followerID, followeeID)
	if err != nil {
		global.Logger.Error("service.unfollowUser.check_following_error",
			zap.Uint("follower_id", followerID),
			zap.Uint("followee_id", followeeID),
			zap.Error(err),
		)
		return fmt.Errorf("检查关注状态失败")
	}

	if !isFollowing {
		// 未关注，直接返回成功（幂等性）
		global.Logger.Info("service.unfollowUser.not_followed",
			zap.Uint("follower_id", followerID),
			zap.Uint("followee_id", followeeID),
		)
		return nil
	}

	// 使用事务：删除关注关系 + 更新双方统计数
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 删除关注记录
		if err := s.relationDAO.DeleteRelation(ctx, followerID, followeeID); err != nil {
			return err
		}

		// 减少关注者的关注数
		if err := s.userDAO.DecrementFollowCount(ctx, followerID); err != nil {
			return err
		}

		// 减少被关注者的粉丝数
		if err := s.userDAO.DecrementFollowerCount(ctx, followeeID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		global.Logger.Error("service.unfollowUser.transaction_error",
			zap.Uint("follower_id", followerID),
			zap.Uint("followee_id", followeeID),
			zap.Error(err),
		)
		return fmt.Errorf("取消关注失败")
	}

	global.Logger.Info("service.unfollowUser.success",
		zap.Uint("follower_id", followerID),
		zap.Uint("followee_id", followeeID),
	)

	return nil
}

// GetFollowList 获取关注列表
func (s *RelationService) GetFollowList(ctx context.Context, currentUserID, targetUserID uint) ([]*dto.UserInfo, error) {
	// 获取目标用户的关注列表（用户ID列表）
	followeeIDs, err := s.relationDAO.GetFollowList(ctx, targetUserID)
	if err != nil {
		global.Logger.Error("service.GetFollowList.get_follow_list_error",
			zap.Uint("target_user_id", targetUserID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("查询关注列表失败")
	}

	// 如果没有关注任何人，返回空列表
	if len(followeeIDs) == 0 {
		return []*dto.UserInfo{}, nil
	}

	// 批量查询用户信息
	users, err := s.userDAO.GetUsersByIDs(ctx, followeeIDs)
	if err != nil {
		global.Logger.Error("service.GetFollowList.get_users_error",
			zap.Uint("target_user_id", targetUserID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("查询用户信息失败")
	}

	// 批量查询当前用户对这些用户的关注状态
	var followMap map[uint]bool
	if currentUserID > 0 {
		followMap, err = s.relationDAO.BatchCheckFollowing(ctx, currentUserID, followeeIDs)
		if err != nil {
			global.Logger.Error("service.GetFollowList.batch_check_error",
				zap.Uint("current_user_id", currentUserID),
				zap.Error(err),
			)
			// 不阻断流程，继续处理
			followMap = make(map[uint]bool)
		}
	} else {
		followMap = make(map[uint]bool)
	}

	// 构建用户信息 DTO 列表
	userList := make([]*dto.UserInfo, 0, len(users))
	for _, user := range users {
		userInfo := &dto.UserInfo{
			ID:            user.ID,
			Username:      user.Username,
			Avatar:        user.Avatar,
			Signature:     user.Signature,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      followMap[user.ID],
		}
		userList = append(userList, userInfo)
	}

	global.Logger.Info("service.GetFollowList.success",
		zap.Uint("target_user_id", targetUserID),
		zap.Int("count", len(userList)),
	)

	return userList, nil
}

// GetFollowerList 获取粉丝列表
func (s *RelationService) GetFollowerList(ctx context.Context, currentUserID, targetUserID uint) ([]*dto.UserInfo, error) {
	// 获取目标用户的粉丝列表（用户ID列表）
	followerIDs, err := s.relationDAO.GetFollowerList(ctx, targetUserID)
	if err != nil {
		global.Logger.Error("service.GetFollowerList.get_follower_list_error",
			zap.Uint("target_user_id", targetUserID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("查询粉丝列表失败")
	}

	// 如果没有粉丝，返回空列表
	if len(followerIDs) == 0 {
		return []*dto.UserInfo{}, nil
	}

	// 批量查询用户信息
	users, err := s.userDAO.GetUsersByIDs(ctx, followerIDs)
	if err != nil {
		global.Logger.Error("service.GetFollowerList.get_users_error",
			zap.Uint("target_user_id", targetUserID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("查询用户信息失败")
	}

	// 批量查询当前用户对这些用户的关注状态
	var followMap map[uint]bool
	if currentUserID > 0 {
		followMap, err = s.relationDAO.BatchCheckFollowing(ctx, currentUserID, followerIDs)
		if err != nil {
			global.Logger.Error("service.GetFollowerList.batch_check_error",
				zap.Uint("current_user_id", currentUserID),
				zap.Error(err),
			)
			// 不阻断流程，继续处理
			followMap = make(map[uint]bool)
		}
	} else {
		followMap = make(map[uint]bool)
	}

	// 构建用户信息 DTO 列表
	userList := make([]*dto.UserInfo, 0, len(users))
	for _, user := range users {
		userInfo := &dto.UserInfo{
			ID:            user.ID,
			Username:      user.Username,
			Avatar:        user.Avatar,
			Signature:     user.Signature,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      followMap[user.ID],
		}
		userList = append(userList, userInfo)
	}

	global.Logger.Info("service.GetFollowerList.success",
		zap.Uint("target_user_id", targetUserID),
		zap.Int("count", len(userList)),
	)

	return userList, nil
}

// GetFriendList 获取好友列表（互相关注）
func (s *RelationService) GetFriendList(ctx context.Context, userID uint) ([]*dto.FriendInfo, error) {
	// 获取好友列表（互相关注的用户ID列表）
	friendIDs, err := s.relationDAO.GetFriendList(ctx, userID)
	if err != nil {
		global.Logger.Error("service.GetFriendList.get_friend_list_error",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("查询好友列表失败")
	}

	// 如果没有好友，返回空列表
	if len(friendIDs) == 0 {
		return []*dto.FriendInfo{}, nil
	}

	// 批量查询用户信息
	users, err := s.userDAO.GetUsersByIDs(ctx, friendIDs)
	if err != nil {
		global.Logger.Error("service.GetFriendList.get_users_error",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("查询用户信息失败")
	}

	// 构建好友信息 DTO 列表
	friendList := make([]*dto.FriendInfo, 0, len(users))
	for _, user := range users {
		friendInfo := &dto.FriendInfo{
			UserInfo: dto.UserInfo{
				ID:            user.ID,
				Username:      user.Username,
				Avatar:        user.Avatar,
				Signature:     user.Signature,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				IsFollow:      true, // 好友必定是互相关注的
			},
			Message: "", // 暂时为空，等消息模块实现后再填充
			MsgType: 0,  // 暂时为0
		}
		friendList = append(friendList, friendInfo)
	}

	global.Logger.Info("service.GetFriendList.success",
		zap.Uint("user_id", userID),
		zap.Int("count", len(friendList)),
	)

	return friendList, nil
}

// IsFriend 判断两个用户是否为好友（双向关注）
func (s *RelationService) IsFriend(ctx context.Context, userID1, userID2 uint) (bool, error) {
	// 检查 userID1 是否关注 userID2
	following, err := s.relationDAO.IsFollowing(ctx, userID1, userID2)
	if err != nil {
		global.Logger.Error("service.IsFriend.check_following_error",
			zap.Uint("user_id_1", userID1),
			zap.Uint("user_id_2", userID2),
			zap.Error(err),
		)
		return false, fmt.Errorf("检查关注关系失败")
	}

	if !following {
		return false, nil
	}

	// 检查 userID2 是否关注 userID1
	followed, err := s.relationDAO.IsFollowing(ctx, userID2, userID1)
	if err != nil {
		global.Logger.Error("service.IsFriend.check_followed_error",
			zap.Uint("user_id_1", userID1),
			zap.Uint("user_id_2", userID2),
			zap.Error(err),
		)
		return false, fmt.Errorf("检查关注关系失败")
	}

	return followed, nil
}
