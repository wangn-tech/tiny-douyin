package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wangn-tech/tiny-douyin/internal/common/errc"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"` // 移除 omitempty，确保始终输出
}

// Success 成功响应（使用默认成功消息）
func Success(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		StatusCode: errc.Success,
		StatusMsg:  errc.GetMsg(errc.Success),
	})
}

// SuccessWithMsg 成功响应（自定义消息）
func SuccessWithMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		StatusCode: errc.Success,
		StatusMsg:  msg,
	})
}

// SuccessWithData 成功响应（带数据，数据结构需包含 Response）
func SuccessWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// Error 错误响应（使用错误码 + 自定义消息）
func Error(c *gin.Context, code int32, msg string) {
	c.JSON(http.StatusOK, Response{
		StatusCode: code,
		StatusMsg:  msg,
	})
}

// ErrorWithCode 错误响应（使用错误码 + 标准消息）
func ErrorWithCode(c *gin.Context, code int32) {
	c.JSON(http.StatusOK, Response{
		StatusCode: code,
		StatusMsg:  errc.GetMsg(code),
	})
}
