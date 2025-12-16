package dto

// RelationActionRequest 关注操作请求
type RelationActionRequest struct {
	ToUserID   uint  `form:"to_user_id" binding:"required"`  // 目标用户ID
	ActionType int32 `form:"action_type" binding:"required"` // 操作类型：1-关注，2-取消关注
}

// RelationActionResponse 关注操作响应
type RelationActionResponse struct {
	StatusCode int32  `json:"status_code"` // 状态码
	StatusMsg  string `json:"status_msg"`  // 状态信息
}

// RelationFollowListRequest 关注列表请求
type RelationFollowListRequest struct {
	UserID uint `form:"user_id" binding:"required"` // 用户ID
}

// RelationFollowListResponse 关注列表响应
type RelationFollowListResponse struct {
	StatusCode int32       `json:"status_code"` // 状态码
	StatusMsg  string      `json:"status_msg"`  // 状态信息
	UserList   []*UserInfo `json:"user_list"`   // 用户列表
}

// RelationFollowerListRequest 粉丝列表请求
type RelationFollowerListRequest struct {
	UserID uint `form:"user_id" binding:"required"` // 用户ID
}

// RelationFollowerListResponse 粉丝列表响应
type RelationFollowerListResponse struct {
	StatusCode int32       `json:"status_code"` // 状态码
	StatusMsg  string      `json:"status_msg"`  // 状态信息
	UserList   []*UserInfo `json:"user_list"`   // 用户列表
}

// RelationFriendListRequest 好友列表请求
type RelationFriendListRequest struct {
	UserID uint `form:"user_id" binding:"required"` // 用户ID
}

// RelationFriendListResponse 好友列表响应
type RelationFriendListResponse struct {
	StatusCode int32         `json:"status_code"` // 状态码
	StatusMsg  string        `json:"status_msg"`  // 状态信息
	UserList   []*FriendInfo `json:"user_list"`   // 好友列表
}

// FriendInfo 好友信息（扩展 UserInfo，添加消息字段）
type FriendInfo struct {
	UserInfo
	Message string `json:"message,omitempty"` // 和该好友的最新聊天消息（暂时为空，等消息模块实现）
	MsgType int32  `json:"msg_type"`          // 消息类型：0-当前用户接收的消息，1-当前用户发送的消息
}
