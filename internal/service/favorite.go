package service

import (
	"context"
	"errors"
	"time"

	"github.com/wangn-tech/tiny-douyin/internal/api/dto"
	"github.com/wangn-tech/tiny-douyin/internal/common/constant"
	"github.com/wangn-tech/tiny-douyin/internal/common/errc"
	"github.com/wangn-tech/tiny-douyin/internal/dao"
	"github.com/wangn-tech/tiny-douyin/internal/global"
	"github.com/wangn-tech/tiny-douyin/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// IFavoriteService 点赞服务接口
type IFavoriteService interface {
	// FavoriteAction 点赞操作（点赞或取消点赞）
	FavoriteAction(ctx context.Context, userID, videoID uint, actionType int32) error
	// GetFavoriteList 获取用户喜欢的视频列表
	GetFavoriteList(ctx context.Context, userID, currentUserID uint) (*dto.FavoriteListData, error)
}

// FavoriteService 点赞服务实现
type FavoriteService struct {
	favoriteDAO dao.IFavoriteDAO
	videoDAO    dao.IVideoDAO
	userDAO     dao.IUserDAO
	relationDAO dao.IRelationDAO
}

// NewFavoriteService 创建 FavoriteService 实例
func NewFavoriteService(
	favoriteDAO dao.IFavoriteDAO,
	videoDAO dao.IVideoDAO,
	userDAO dao.IUserDAO,
	relationDAO dao.IRelationDAO,
) IFavoriteService {
	return &FavoriteService{
		favoriteDAO: favoriteDAO,
		videoDAO:    videoDAO,
		userDAO:     userDAO,
		relationDAO: relationDAO,
	}
}

// FavoriteAction 点赞操作（点赞或取消点赞）
func (s *FavoriteService) FavoriteAction(ctx context.Context, userID, videoID uint, actionType int32) error {
	start := time.Now()

	global.Logger.Info("service.FavoriteAction.start",
		zap.Uint("user_id", userID),
		zap.Uint("video_id", videoID),
		zap.Int32("action_type", actionType),
	)

	// 1. 验证操作类型
	if actionType != constant.FavoriteActionLike && actionType != constant.FavoriteActionUnlike {
		global.Logger.Warn("service.FavoriteAction.invalid_action_type",
			zap.Int32("action_type", actionType),
		)
		return errors.New(errc.ErrMsg[errc.ErrInvalidActionType])
	}

	// 2. 检查视频是否存在
	video, err := s.videoDAO.GetVideoByID(ctx, videoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.Warn("service.FavoriteAction.video_not_found",
				zap.Uint("video_id", videoID),
			)
			return errors.New(errc.ErrMsg[errc.ErrVideoNotFound])
		}
		global.Logger.Error("service.FavoriteAction.get_video_error",
			zap.Uint("video_id", videoID),
			zap.Error(err),
		)
		return err
	}

	// 3. 检查当前点赞状态
	isFavorited, err := s.favoriteDAO.IsFavorite(ctx, userID, videoID)
	if err != nil {
		global.Logger.Error("service.FavoriteAction.check_favorite_error",
			zap.Uint("user_id", userID),
			zap.Uint("video_id", videoID),
			zap.Error(err),
		)
		return err
	}

	// 4. 根据操作类型和当前状态处理
	if actionType == constant.FavoriteActionLike {
		// 点赞操作
		if isFavorited {
			// 已经点赞过，幂等返回成功
			global.Logger.Info("service.FavoriteAction.already_favorited",
				zap.Uint("user_id", userID),
				zap.Uint("video_id", videoID),
			)
			return nil
		}

		// 执行点赞：创建点赞记录 + 增加视频点赞数
		if err := s.favoriteDAO.CreateFavorite(ctx, userID, videoID); err != nil {
			global.Logger.Error("service.FavoriteAction.create_favorite_error",
				zap.Uint("user_id", userID),
				zap.Uint("video_id", videoID),
				zap.Error(err),
			)
			return err
		}

		if err := s.videoDAO.IncrementFavoriteCount(ctx, videoID); err != nil {
			global.Logger.Error("service.FavoriteAction.increment_count_error",
				zap.Uint("video_id", videoID),
				zap.Error(err),
			)
			// 注意：这里可能导致数据不一致，实际应该用事务处理
			// 为了简化，这里先不回滚，后续可以优化
			return err
		}

		global.Logger.Info("service.FavoriteAction.like_success",
			zap.Uint("user_id", userID),
			zap.Uint("video_id", videoID),
			zap.Duration("duration", time.Since(start)),
		)

	} else {
		// 取消点赞操作
		if !isFavorited {
			// 未点赞过，幂等返回成功
			global.Logger.Info("service.FavoriteAction.not_favorited",
				zap.Uint("user_id", userID),
				zap.Uint("video_id", videoID),
			)
			return nil
		}

		// 执行取消点赞：删除点赞记录 + 减少视频点赞数
		if err := s.favoriteDAO.DeleteFavorite(ctx, userID, videoID); err != nil {
			global.Logger.Error("service.FavoriteAction.delete_favorite_error",
				zap.Uint("user_id", userID),
				zap.Uint("video_id", videoID),
				zap.Error(err),
			)
			return err
		}

		if err := s.videoDAO.DecrementFavoriteCount(ctx, videoID); err != nil {
			global.Logger.Error("service.FavoriteAction.decrement_count_error",
				zap.Uint("video_id", videoID),
				zap.Error(err),
			)
			// 注意：这里可能导致数据不一致，实际应该用事务处理
			return err
		}

		global.Logger.Info("service.FavoriteAction.unlike_success",
			zap.Uint("user_id", userID),
			zap.Uint("video_id", videoID),
			zap.Duration("duration", time.Since(start)),
		)
	}

	// 记录最终视频点赞数（用于调试）
	global.Logger.Info("service.FavoriteAction.final_count",
		zap.Uint("video_id", videoID),
		zap.Int64("favorite_count", video.FavoriteCount),
	)

	return nil
}

// GetFavoriteList 获取用户喜欢的视频列表
func (s *FavoriteService) GetFavoriteList(ctx context.Context, userID, currentUserID uint) (*dto.FavoriteListData, error) {
	start := time.Now()

	global.Logger.Info("service.GetFavoriteList.start",
		zap.Uint("user_id", userID),
		zap.Uint("current_user_id", currentUserID),
	)

	// 1. 获取用户点赞的视频ID列表
	videoIDs, err := s.favoriteDAO.GetUserFavoriteVideoIDs(ctx, userID)
	if err != nil {
		global.Logger.Error("service.GetFavoriteList.get_video_ids_error",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return nil, err
	}

	if len(videoIDs) == 0 {
		global.Logger.Info("service.GetFavoriteList.empty_list",
			zap.Uint("user_id", userID),
		)
		return &dto.FavoriteListData{
			Videos: []dto.Video{},
		}, nil
	}

	// 2. 批量查询视频详情
	videos, err := s.videoDAO.GetVideosByIDs(ctx, videoIDs)
	if err != nil {
		global.Logger.Error("service.GetFavoriteList.get_videos_error",
			zap.Int("video_count", len(videoIDs)),
			zap.Error(err),
		)
		return nil, err
	}

	// 3. 构建视频DTO列表（包含作者信息和统计信息）
	videoList, err := s.buildVideoDTOListFromModels(ctx, videos, currentUserID)
	if err != nil {
		global.Logger.Error("service.GetFavoriteList.build_dto_error",
			zap.Error(err),
		)
		return nil, err
	}

	global.Logger.Info("service.GetFavoriteList.success",
		zap.Uint("user_id", userID),
		zap.Int("video_count", len(videoList)),
		zap.Duration("duration", time.Since(start)),
	)

	return &dto.FavoriteListData{
		Videos: videoList,
	}, nil
}

// buildVideoDTOList 构建视频DTO列表（包含作者信息、统计信息）
func (s *FavoriteService) buildVideoDTOListFromModels(ctx context.Context, videos []*model.Video, currentUserID uint) ([]dto.Video, error) {
	if len(videos) == 0 {
		return []dto.Video{}, nil
	}

	var videoList []dto.Video

	// 收集所有作者ID和视频ID
	authorIDs := make(map[uint]bool)
	videoIDs := make([]uint, 0, len(videos))
	for _, video := range videos {
		authorIDs[video.AuthorID] = true
		videoIDs = append(videoIDs, video.ID)
	}

	// 批量查询作者信息（优化：避免N+1查询）
	authorMap := make(map[uint]*dto.UserInfo)
	for authorID := range authorIDs {
		user, err := s.userDAO.GetUserByID(ctx, authorID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if user != nil {
			authorMap[authorID] = &dto.UserInfo{
				ID:            user.ID,
				Username:      user.Username,
				Avatar:        user.Avatar,
				Signature:     user.Signature,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				IsFollow:      false, // 稍后批量更新
			}
		}
	}

	// 批量检查当前用户是否点赞了这些视频
	var favoriteMap map[uint]bool
	if currentUserID > 0 {
		var err error
		favoriteMap, err = s.favoriteDAO.BatchCheckFavorite(ctx, currentUserID, videoIDs)
		if err != nil {
			global.Logger.Error("service.buildVideoDTOListFromModels.batch_check_error",
				zap.Uint("current_user_id", currentUserID),
				zap.Error(err),
			)
			// 不阻断流程，继续处理
			favoriteMap = make(map[uint]bool)
		}
	} else {
		favoriteMap = make(map[uint]bool)
	}

	// 批量查询当前用户对作者的关注状态（如果已登录）
	if currentUserID > 0 {
		authorIDList := make([]uint, 0, len(authorIDs))
		for authorID := range authorIDs {
			authorIDList = append(authorIDList, authorID)
		}
		followMap, err := s.relationDAO.BatchCheckFollowing(ctx, currentUserID, authorIDList)
		if err != nil {
			global.Logger.Error("service.buildVideoDTOListFromModels.batch_check_following_error",
				zap.Uint("current_user_id", currentUserID),
				zap.Error(err),
			)
			// 不阻断流程，继续处理
		} else {
			// 更新作者的关注状态
			for authorID, author := range authorMap {
				author.IsFollow = followMap[authorID]
			}
		}
	}

	// 组装视频DTO
	for _, video := range videos {
		author := authorMap[video.AuthorID]
		if author == nil {
			continue // 跳过作者不存在的视频
		}

		videoDTO := dto.Video{
			ID:            video.ID,
			PlayURL:       video.PlayURL,
			CoverURL:      video.CoverURL,
			Title:         video.Title,
			Author:        *author,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    favoriteMap[video.ID],
		}

		videoList = append(videoList, videoDTO)
	}

	return videoList, nil
}
