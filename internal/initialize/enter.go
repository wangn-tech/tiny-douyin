package initialize

import (
	"github.com/wangn-tech/tiny-douyin/internal/config"
	"github.com/wangn-tech/tiny-douyin/internal/global"
)

func InitAll() {
	// 配置
	global.Config = config.Init()

	// 日志
	global.Logger = LoggerSetup(global.Config.Log)

	// 数据库
	global.DB = InitDB()

	// 迁移
	_ = AutoMigrate(global.DB)

	// Redis
	global.RedisClient = InitRedis()
}
