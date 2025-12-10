package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/wangn-tech/tiny-douyin/internal/global"
	"github.com/wangn-tech/tiny-douyin/internal/initialize"
	"github.com/wangn-tech/tiny-douyin/internal/router"
)

func main() {
	// 初始化
	initialize.InitAll()

	// 初始化路由
	gin.SetMode(global.Config.Server.Mode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	router.Init(r)

	// 启动服务
	port := fmt.Sprintf(":%d", global.Config.Server.Port)
	log.Println("Server starting on port " + port)
	if err := r.Run(port); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
