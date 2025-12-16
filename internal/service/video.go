package service

import (
	"context"
	"errors"
	"time"

	"github.com/wangn-tech/tiny-douyin/internal/api/dto"
	"github.com/wangn-tech/tiny-douyin/internal/dao"
	"github.com/wangn-tech/tiny-douyin/internal/global"
	"github.com/wangn-tech/tiny-douyin/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// IVideoService 视频服务接口
type IVideoService interface {
	// PublishVideo 发布视频，返回视频 ID
	PublishVideo(ctx context.Context, req *dto.VideoPublishRequest, authorID uint, playURL, coverURL string) (uint, error)
	// GetVideoFeed 获取视频流
	GetVideoFeed(ctx context.Context, req *dto.VideoFeedRequest, currentUserID uint) (*dto.VideoFeedData, error)
	// GetVideoList 获取用户发布的视频列表
	GetVideoList(ctx context.Context, req *dto.VideoListRequest, currentUserID uint) (*dto.VideoListData, error)
}

// VideoService 视频服务实现
type VideoService struct {
	videoDAO    dao.IVideoDAO
	userDAO     dao.IUserDAO
	favoriteDAO dao.IFavoriteDAO
	relationDAO dao.IRelationDAO
}

// NewVideoService 创建 VideoService 实例
func NewVideoService(
	videoDAO dao.IVideoDAO,
	userDAO dao.IUserDAO,
	favoriteDAO dao.IFavoriteDAO,
	relationDAO dao.IRelationDAO,
) IVideoService {
	return &VideoService{
		videoDAO:    videoDAO,
		userDAO:     userDAO,
		favoriteDAO: favoriteDAO,
		relationDAO: relationDAO,
	}
}

// PublishVideo 发布视频，返回视频 ID
func (s *VideoService) PublishVideo(ctx context.Context, req *dto.VideoPublishRequest, authorID uint, playURL, coverURL string) (uint, error) {
	start := time.Now()

	global.Logger.Info("service.PublishVideo.start",
		zap.Uint("author_id", authorID),
		zap.String("title", req.Title),
	)

	// 创建视频记录
	video := &model.Video{
		AuthorID: authorID,
		PlayURL:  playURL,
		CoverURL: coverURL,
		Title:    req.Title,
	}

	if err := s.videoDAO.CreateVideo(ctx, video); err != nil {
		global.Logger.Error("service.PublishVideo.create_error",
			zap.Uint("author_id", authorID),
			zap.Error(err),
		)
		return 0, err
	}

	global.Logger.Info("service.PublishVideo.success",
		zap.Uint("author_id", authorID),
		zap.Uint("video_id", video.ID),
		zap.Duration("duration", time.Since(start)),
	)

	return video.ID, nil
}

// GetVideoFeed 获取视频流
func (s *VideoService) GetVideoFeed(ctx context.Context, req *dto.VideoFeedRequest, currentUserID uint) (*dto.VideoFeedData, error) {
	start := time.Now()

	latestTime := req.LatestTime
	// 注意：如果 latestTime 为 0，表示首次请求，不做时间过滤
	// 否则，返回比 latestTime 更早的视频

	global.Logger.Info("service.GetVideoFeed.start",
		zap.Int64("latest_time", latestTime),
		zap.Uint("current_user_id", currentUserID),
	)

	// 查询视频列表（限制30条）
	videos, err := s.videoDAO.GetVideoFeed(ctx, latestTime, 30)
	if err != nil {
		global.Logger.Error("service.GetVideoFeed.query_error",
			zap.Error(err),
		)
		return nil, err
	}

	global.Logger.Info("service.GetVideoFeed.dao_result",
		zap.Int("video_count", len(videos)),
	)

	// 组装响应数据
	videoList, nextTime, err := s.buildVideoDTOList(ctx, videos, currentUserID)
	if err != nil {
		global.Logger.Error("service.GetVideoFeed.build_error",
			zap.Error(err),
		)
		return nil, err
	}

	global.Logger.Info("service.GetVideoFeed.success",
		zap.Int("video_count", len(videoList)),
		zap.Int64("next_time", nextTime),
		zap.Duration("duration", time.Since(start)),
	)

	return &dto.VideoFeedData{
		NextTime: nextTime,
		Videos:   videoList,
	}, nil
}

// GetVideoList 获取用户发布的视频列表
func (s *VideoService) GetVideoList(ctx context.Context, req *dto.VideoListRequest, currentUserID uint) (*dto.VideoListData, error) {
	start := time.Now()

	global.Logger.Info("service.GetVideoList.start",
		zap.Uint("user_id", req.UserID),
		zap.Uint("current_user_id", currentUserID),
	)

	// 查询用户视频列表
	videos, err := s.videoDAO.GetVideosByUserID(ctx, req.UserID)
	if err != nil {
		global.Logger.Error("service.GetVideoList.query_error",
			zap.Uint("user_id", req.UserID),
			zap.Error(err),
		)
		return nil, err
	}

	// 组装响应数据
	videoList, _, err := s.buildVideoDTOList(ctx, videos, currentUserID)
	if err != nil {
		global.Logger.Error("service.GetVideoList.build_error",
			zap.Error(err),
		)
		return nil, err
	}

	global.Logger.Info("service.GetVideoList.success",
		zap.Uint("user_id", req.UserID),
		zap.Int("video_count", len(videoList)),
		zap.Duration("duration", time.Since(start)),
	)

	return &dto.VideoListData{
		Videos: videoList,
	}, nil
}

// buildVideoDTOList 构建视频DTO列表（包含作者信息、统计信息）
func (s *VideoService) buildVideoDTOList(ctx context.Context, videos []*model.Video, currentUserID uint) ([]dto.Video, int64, error) {
	if len(videos) == 0 {
		return []dto.Video{}, time.Now().Unix(), nil
	}

	var videoList []dto.Video
	var nextTime int64

	// 收集所有作者ID和视频ID
	authorIDs := make(map[uint]bool)
	videoIDs := make([]uint, 0, len(videos))
	for _, video := range videos {
		authorIDs[video.AuthorID] = true
		videoIDs = append(videoIDs, video.ID)
	}

	// 批量查询作者信息（优化：避免N+1查询）
	authorMap := make(map[uint]*model.User)
	for authorID := range authorIDs {
		user, err := s.userDAO.GetUserByID(ctx, authorID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, err
		}
		if user != nil {
			authorMap[authorID] = user
		}
	}

	// 批量查询当前用户的点赞状态（如果已登录）
	var favoriteMap map[uint]bool
	if currentUserID > 0 {
		var err error
		favoriteMap, err = s.favoriteDAO.BatchCheckFavorite(ctx, currentUserID, videoIDs)
		if err != nil {
			global.Logger.Error("service.buildVideoDTOList.batch_check_favorite_error",
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
	var followMap map[uint]bool
	if currentUserID > 0 {
		authorIDList := make([]uint, 0, len(authorIDs))
		for authorID := range authorIDs {
			authorIDList = append(authorIDList, authorID)
		}
		var err error
		followMap, err = s.relationDAO.BatchCheckFollowing(ctx, currentUserID, authorIDList)
		if err != nil {
			global.Logger.Error("service.buildVideoDTOList.batch_check_following_error",
				zap.Uint("current_user_id", currentUserID),
				zap.Error(err),
			)
			// 不阻断流程，继续处理
			followMap = make(map[uint]bool)
		}
	} else {
		followMap = make(map[uint]bool)
	}

	// 组装视频DTO
	for _, video := range videos {
		author := authorMap[video.AuthorID]
		if author == nil {
			continue // 跳过作者不存在的视频
		}

		videoDTO := dto.Video{
			ID:       video.ID,
			PlayURL:  video.PlayURL,
			CoverURL: video.CoverURL,
			Title:    video.Title,
			Author: dto.UserInfo{
				ID:            author.ID,
				Username:      author.Username,
				Avatar:        author.Avatar,
				Signature:     author.Signature,
				FollowCount:   author.FollowCount,
				FollowerCount: author.FollowerCount,
				IsFollow:      followMap[author.ID],
			},
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    favoriteMap[video.ID],
		}

		videoList = append(videoList, videoDTO)

		// 更新 nextTime（最早的视频时间）
		videoTime := video.CreatedAt.Unix()
		if nextTime == 0 || videoTime < nextTime {
			nextTime = videoTime
		}
	}

	return videoList, nextTime, nil
}
