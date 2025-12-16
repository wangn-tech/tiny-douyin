package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Init loads configuration from config-{env}.yaml (default env=dev)
func Init() *AppConfig {
	// 解析命令行参数 env, 默认为 dev
	env := pflag.String("env", "dev", "Specify the environment config to use: [dev, prod, test]")
	pflag.Parse()

	// 设置配置文件: ./config/config-{env}.yaml
	v := viper.New()
	v.AddConfigPath("./config")
	v.SetConfigName(fmt.Sprintf("config-%s", *env))
	v.SetConfigType("yaml")

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error reading config file: %w", err))
	}

	// 解析配置文件到结构体
	var c AppConfig
	if err := v.Unmarshal(&c); err != nil {
		panic(fmt.Errorf("unable to decode config: %w", err))
	}

	return &c
}

// AppConfig 整合所有配置
type AppConfig struct {
	Server   Server         `mapstructure:"server"`
	MySQL    MySQLConfig    `mapstructure:"mysql"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	MinIO    MinIOConfig    `mapstructure:"minio"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
}

type Server struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// MySQLConfig MySQL 配置
type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	TTL    int    `mapstructure:"ttl"`
}

// LogConfig Zap 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`        // debug, info, warn, error
	Format     string `mapstructure:"format"`       // json, console
	Output     string `mapstructure:"output"`       // stdout, stderr, file
	FilePath   string `mapstructure:"file_path"`    // 当 output=file 时使用
	MaxSizeMB  int    `mapstructure:"max_size_mb"`  // 滚动大小（可选，后续扩展）
	MaxBackups int    `mapstructure:"max_backups"`  // 保留文件数（可选）
	MaxAgeDays int    `mapstructure:"max_age_days"` // 保留天数（可选）
}

// MinIOConfig MinIO 对象存储配置
type MinIOConfig struct {
	Endpoint        string `mapstructure:"endpoint"`          // MinIO 服务地址
	AccessKeyID     string `mapstructure:"access_key_id"`     // 访问密钥 ID
	SecretAccessKey string `mapstructure:"secret_access_key"` // 访问密钥密码
	UseSSL          bool   `mapstructure:"use_ssl"`           // 是否使用 HTTPS
	BucketName      string `mapstructure:"bucket_name"`       // 存储桶名称
	Location        string `mapstructure:"location"`          // 存储桶位置（区域）
	URLPrefix       string `mapstructure:"url_prefix"`        // 文件访问 URL 前缀
}

// RabbitMQConfig RabbitMQ 消息队列配置
type RabbitMQConfig struct {
	Host     string `mapstructure:"host"`     // RabbitMQ 服务地址
	Port     int    `mapstructure:"port"`     // RabbitMQ 服务端口
	User     string `mapstructure:"user"`     // 用户名
	Password string `mapstructure:"password"` // 密码
	VHost    string `mapstructure:"vhost"`    // 虚拟主机
	Exchange string `mapstructure:"exchange"` // 交换机名称
	Queue    string `mapstructure:"queue"`    // 队列名称
}
