package dto

// CommentActionRequest 评论操作请求
type CommentActionRequest struct {
	VideoID     uint   `form:"video_id" binding:"required"`    // 视频ID
	ActionType  int32  `form:"action_type" binding:"required"` // 操作类型：1-发布评论，2-删除评论
	CommentText string `form:"comment_text"`                   // 评论内容（发布评论时使用）
	CommentID   uint   `form:"comment_id"`                     // 评论ID（删除评论时使用）
}

// CommentActionResponse 评论操作响应
type CommentActionResponse struct {
	StatusCode int32    `json:"status_code"`       // 状态码
	StatusMsg  string   `json:"status_msg"`        // 状态信息
	Comment    *Comment `json:"comment,omitempty"` // 发布评论时返回评论信息
}

// CommentListRequest 评论列表请求
type CommentListRequest struct {
	VideoID uint `form:"video_id" binding:"required"` // 视频ID
}

// CommentListResponse 评论列表响应
type CommentListResponse struct {
	StatusCode  int32      `json:"status_code"`  // 状态码
	StatusMsg   string     `json:"status_msg"`   // 状态信息
	CommentList []*Comment `json:"comment_list"` // 评论列表
}

// Comment 评论信息
type Comment struct {
	ID         uint      `json:"id"`          // 评论ID
	User       *UserInfo `json:"user"`        // 评论用户信息
	Content    string    `json:"content"`     // 评论内容
	CreateDate string    `json:"create_date"` // 评论发布日期（MM-DD）
}
