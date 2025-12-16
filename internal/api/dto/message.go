package dto

// MessageActionRequest 发送消息请求
type MessageActionRequest struct {
	ToUserID   uint   `form:"to_user_id" binding:"required,gt=0"`     // 接收者用户ID
	ActionType int32  `form:"action_type" binding:"required,oneof=1"` // 操作类型：1-发送消息
	Content    string `form:"content" binding:"required"`             // 消息内容
}

// MessageActionResponse 发送消息响应
type MessageActionResponse struct {
	StatusCode int32  `json:"status_code"` // 状态码
	StatusMsg  string `json:"status_msg"`  // 状态信息
}

// MessageChatRequest 获取聊天记录请求
type MessageChatRequest struct {
	ToUserID   uint  `form:"to_user_id" binding:"required,gt=0"` // 对方用户ID
	PreMsgTime int64 `form:"pre_msg_time"`                       // 上次最新消息的时间戳（秒），可选，用于分页
}

// MessageChatResponse 获取聊天记录响应
type MessageChatResponse struct {
	StatusCode  int32     `json:"status_code"`  // 状态码
	StatusMsg   string    `json:"status_msg"`   // 状态信息
	MessageList []Message `json:"message_list"` // 消息列表
}

// Message 消息信息
type Message struct {
	ID         uint   `json:"id"`           // 消息ID
	ToUserID   uint   `json:"to_user_id"`   // 接收者用户ID
	FromUserID uint   `json:"from_user_id"` // 发送者用户ID
	Content    string `json:"content"`      // 消息内容
	CreateTime int64  `json:"create_time"`  // 消息发送时间（Unix时间戳，秒）
}
