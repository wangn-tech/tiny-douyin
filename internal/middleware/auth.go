package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wangn-tech/tiny-douyin/internal/common/errc"
	"github.com/wangn-tech/tiny-douyin/internal/common/response"
	"github.com/wangn-tech/tiny-douyin/internal/pkg/jwt"
)

// JWTAuth JWT 验证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查 token 是否存在
		token := getToken(c)
		if token == "" {
			response.Error(c, errc.ErrTokenMissing, errc.GetMsg(errc.ErrTokenMissing))
			c.Abort()
			return
		}

		// 解析和验证 token
		claims, err := jwt.ParseToken(token)
		if err != nil {
			response.Error(c, errc.ErrTokenInvalid, errc.GetMsg(errc.ErrTokenInvalid))
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

// JWTAuthOptional
func JWTAuthOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)
		if token != "" {
			claims, err := jwt.ParseToken(token)
			if err == nil {
				c.Set("user_id", claims.UserID)
			}
		}
		c.Next()
	}
}

// getToken 从请求中提取 JWT token
func getToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}

	token = c.PostForm("token")
	if token != "" {
		return token
	}

	return c.GetHeader("Authorization")
}
