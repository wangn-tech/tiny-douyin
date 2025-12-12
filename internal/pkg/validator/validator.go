package validator

import (
	"github.com/gin-gonic/gin"
	"github.com/wangn-tech/tiny-douyin/internal/common/errc"
	"github.com/wangn-tech/tiny-douyin/internal/common/response"
)

func ValidateParams(c *gin.Context, params ...string) bool {
	for _, param := range params {
		if c.Query(param) == "" && c.PostForm(param) == "" {
			response.Error(c, errc.ErrInvalidParams, "缺少必需参数: "+param)
			return false
		}
	}
	return true
}

func GetQueryOrForm(c *gin.Context, key string) string {
	value := c.Query(key)
	if value == "" {
		value = c.PostForm(key)
	}
	return value
}
