package dto

import "github.com/wangn-tech/tiny-douyin/internal/common/response"

// FavoriteActionRequest 点赞操作请求
type FavoriteActionRequest struct {
	Token      string `form:"token" binding:"required"`       // 用户鉴权token
	VideoID    uint   `form:"video_id" binding:"required"`    // 视频ID
	ActionType int32  `form:"action_type" binding:"required"` // 操作类型：1-点赞，2-取消点赞
}

// FavoriteActionResponse 点赞操作响应
type FavoriteActionResponse struct {
	response.Response
}

// FavoriteListRequest 喜欢列表请求
type FavoriteListRequest struct {
	UserID uint   `form:"user_id" binding:"required"` // 用户ID
	Token  string `form:"token"`                      // 用户鉴权token（可选）
}

// FavoriteListResponse 喜欢列表响应
type FavoriteListResponse struct {
	response.Response
	Videos []Video `json:"video_list"` // 用户喜欢的视频列表
}

// FavoriteListData 喜欢列表数据（Service 层返回）
type FavoriteListData struct {
	Videos []Video
}
