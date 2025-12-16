package dto

import "github.com/wangn-tech/tiny-douyin/internal/common/response"

// ============== 请求 DTO ==============

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Username string `json:"username" form:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" form:"password" binding:"required,min=6,max=32"`
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// UserInfoRequest 获取用户信息请求
type UserInfoRequest struct {
	UserID uint `json:"user_id" form:"user_id" binding:"required,gt=0"`
}

// ============== 响应 DTO ==============

// UserRegisterResponse 用户注册响应
type UserRegisterResponse struct {
	response.Response
	UserID uint   `json:"user_id"`
	Token  string `json:"token"`
}

// UserLoginResponse 用户登录响应
type UserLoginResponse struct {
	response.Response
	UserID uint   `json:"user_id"`
	Token  string `json:"token"`
}

// UserInfoResponse 获取用户信息响应
type UserInfoResponse struct {
	response.Response
	User UserInfo `json:"user"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID            uint   `json:"id"`             // 用户ID
	Username      string `json:"name"`           // 用户名
	Avatar        string `json:"avatar"`         // 头像
	Signature     string `json:"signature"`      // 个性签名
	FollowCount   int64  `json:"follow_count"`   // 关注数
	FollowerCount int64  `json:"follower_count"` // 粉丝数
	IsFollow      bool   `json:"is_follow"`      // 是否关注（当前用户是否关注该用户）
}
