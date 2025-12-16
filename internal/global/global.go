package global

import (
	"github.com/minio/minio-go/v7"
	amqp "github.com/rabbitmq/amqp091-go"
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
	MinIOClient *minio.Client     // MinIO 客户端
	RabbitConn  *amqp.Connection  // RabbitMQ 连接
	RabbitChan  *amqp.Channel     // RabbitMQ 频道
)
