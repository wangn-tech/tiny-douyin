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
	ErrTokenInvalid = 2001
	ErrTokenExpired = 2002
	ErrUnauthorized = 2003

	// Video 3xxx
	ErrVideoNotFound = 3001

	// Comment 4xxx
	ErrCommentNotFound = 4001

	// System 9xxx
	ErrInternalServer = 9001
	ErrDatabaseError  = 9002
)

var ErrMsg = map[int]string{
	Success:              "success",
	Failed:               "failed",
	ErrUserNotFound:      "用户不存在",
	ErrUserAlreadyExists: "用户已存在",
	ErrInvalidPassword:   "密码错误",
	ErrInvalidParams:     "参数错误",
	ErrTokenInvalid:      "token无效",
	ErrTokenExpired:      "token已过期",
	ErrUnauthorized:      "未授权",
	ErrVideoNotFound:     "视频不存在",
	ErrCommentNotFound:   "评论不存在",
	ErrInternalServer:    "服务器错误",
	ErrDatabaseError:     "数据库错误",
}

func GetMsg(code int) string {
	if msg, ok := ErrMsg[code]; ok {
		return msg
	}
	return "未知错误"
}
