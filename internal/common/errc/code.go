package errc

const (
	Success = 0
	Failed  = -1

	// User 1xxx
	ErrUserNotFound      = 1001
	ErrUserAlreadyExists = 1002
	ErrInvalidPassword   = 1003
	ErrInvalidParams     = 1004

	// Auth 2xxx
	ErrTokenMissing = 2001
	ErrTokenInvalid = 2002
	ErrTokenExpired = 2003
	ErrUnauthorized = 2004

	// Video 3xxx
	ErrVideoNotFound      = 3001
	ErrVideoUploadFailed  = 3002
	ErrVideoFileInvalid   = 3003
	ErrVideoProcessFailed = 3004

	// Comment 4xxx
	ErrCommentNotFound          = 4001
	ErrCommentPermissionDenied  = 4002
	ErrCommentTooLong           = 4003

	// Favorite 5xxx
	ErrInvalidActionType = 5001
	ErrAlreadyFavorited  = 5002
	ErrNotFavorited      = 5003

	// System 9xxx
	ErrInternalServer = 9001
	ErrDatabaseError  = 9002
)

var ErrMsg = map[int32]string{
	Success:               "success",
	Failed:                "failed",
	ErrUserNotFound:       "用户不存在",
	ErrUserAlreadyExists:  "用户已存在",
	ErrInvalidPassword:    "密码错误",
	ErrInvalidParams:      "参数错误",
	ErrTokenMissing:       "token缺失",
	ErrTokenInvalid:       "token无效",
	ErrTokenExpired:       "token已过期",
	ErrUnauthorized:       "未授权",
	ErrVideoNotFound:      "视频不存在",
	ErrVideoUploadFailed:  "视频上传失败",
	ErrVideoFileInvalid:   "视频文件无效",
	ErrVideoProcessFailed: "视频处理失败",
	ErrCommentNotFound:    "评论不存在",
	ErrCommentPermissionDenied: "无权限删除评论",
	ErrCommentTooLong:     "评论内容过长",
	ErrInvalidActionType:  "无效的操作类型",
	ErrAlreadyFavorited:   "已经点赞过",
	ErrNotFavorited:       "未点赞",
	ErrInternalServer:     "服务器错误",
	ErrDatabaseError:      "数据库错误",
}

func GetMsg(code int32) string {
	if msg, ok := ErrMsg[code]; ok {
		return msg
	}
	return "未知错误"
}
