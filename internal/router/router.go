package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wangn-tech/tiny-douyin/internal/middleware"
)

func Init(r *gin.Engine) {
	// 全局中间件
	// 访问日志
	r.Use(middleware.AccessLog())
	// 跨域请求中间件
	r.Use(middleware.CORS())

	// 健康检查 (ping, pong)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Swagger 路由占位（后续集成）
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 业务路由
	// apiRouter := r.Group("/douyin")
	// 基础接口

}
