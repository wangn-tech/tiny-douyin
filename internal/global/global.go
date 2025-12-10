package global

import (
	"github.com/redis/go-redis/v9"
	"github.com/wangn-tech/tiny-douyin/internal/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config      *config.AppConfig // 全局配置
	DB          *gorm.DB          // MySQL
	RedisClient *redis.Client     // Redis
	Logger      *zap.Logger       // Zap 日志
)
