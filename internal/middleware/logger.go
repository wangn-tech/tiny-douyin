package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wangn-tech/tiny-douyin/internal/global"
	"go.uber.org/zap"
)

// AccessLog 使用 zap 打印访问日志（不含请求链路追踪）
func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		dur := time.Since(start)
		global.Logger.Info("access",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Int64("duration_ms", dur.Milliseconds()),
		)
	}
}
