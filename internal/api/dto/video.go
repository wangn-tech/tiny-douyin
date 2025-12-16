package dto

import "github.com/wangn-tech/tiny-douyin/internal/common/response"

// VideoPublishRequest 视频发布请求
type VideoPublishRequest struct {
	Token string `form:"token"`                   // 用户鉴权token（由中间件验证，这里不做 required 校验）
	Title string `form:"title" binding:"max=128"` // 视频标题（可选，不传则使用文件名）
}

// VideoPublishResponse 视频发布响应
type VideoPublishResponse struct {
	response.Response
}

// VideoFeedRequest 视频流请求
type VideoFeedRequest struct {
	Token      string `form:"token"`       // 可选参数，用户登录状态下传递 token
	LatestTime int64  `form:"latest_time"` // 可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
}

// VideoFeedResponse 视频流响应
type VideoFeedResponse struct {
	response.Response
	NextTime int64   `json:"next_time"` // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	Videos   []Video `json:"video_list"`
}

// VideoFeedData 视频流数据（Service 层返回）
type VideoFeedData struct {
	NextTime int64
	Videos   []Video
}

// VideoListRequest 获取用户发布列表请求
type VideoListRequest struct {
	UserID uint `form:"user_id" binding:"required"` // 用户ID
}

// VideoListResponse 获取用户发布列表响应
type VideoListResponse struct {
	response.Response
	Videos []Video `json:"video_list"`
}

// VideoListData 用户发布列表数据（Service 层返回）
type VideoListData struct {
	Videos []Video
}

// Video 视频信息
type Video struct {
	ID            uint     `json:"id"`             // 视频ID
	Author        UserInfo `json:"author"`         // 作者信息
	PlayURL       string   `json:"play_url"`       // 视频播放地址
	CoverURL      string   `json:"cover_url"`      // 视频封面地址
	FavoriteCount int64    `json:"favorite_count"` // 点赞数
	CommentCount  int64    `json:"comment_count"`  // 评论数
	IsFavorite    bool     `json:"is_favorite"`    // 是否点赞
	Title         string   `json:"title"`          // 视频标题
}
