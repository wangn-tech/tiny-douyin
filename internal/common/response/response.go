package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

func Success(c *gin.Context) {
	c.JSON(http.StatusOK, Response{StatusCode: 0})
}

func SuccessWithMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: msg})
}

func SuccessWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func Error(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{StatusCode: code, StatusMsg: msg})
}
