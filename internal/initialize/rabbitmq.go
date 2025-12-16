package initialize

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/wangn-tech/tiny-douyin/internal/common/constant"
	"github.com/wangn-tech/tiny-douyin/internal/config"
	"github.com/wangn-tech/tiny-douyin/internal/global"
)

// InitRabbitMQ 初始化 RabbitMQ 连接和频道
func InitRabbitMQ(cfg *config.RabbitMQConfig) (*amqp.Connection, *amqp.Channel) {
	// 构建连接字符串
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.VHost)

	// 创建连接
	conn, err := amqp.Dial(connStr)
	if err != nil {
		global.Logger.Fatal("Failed to connect to RabbitMQ: " + err.Error())
	}

	// 创建频道
	ch, err := conn.Channel()
	if err != nil {
		global.Logger.Fatal("Failed to open RabbitMQ channel: " + err.Error())
	}

	// 声明交换机
	err = ch.ExchangeDeclare(
		cfg.Exchange, // 交换机名称
		"direct",     // 交换机类型
		true,         // 持久化
		false,        // 自动删除
		false,        // 内部使用
		false,        // 不等待
		nil,          // 额外参数
	)
	if err != nil {
		global.Logger.Fatal("Failed to declare RabbitMQ exchange: " + err.Error())
	}

	// 声明队列
	_, err = ch.QueueDeclare(
		cfg.Queue, // 队列名称
		true,      // 持久化
		false,     // 自动删除
		false,     // 独占
		false,     // 不等待
		nil,       // 额外参数
	)
	if err != nil {
		global.Logger.Fatal("Failed to declare RabbitMQ queue: " + err.Error())
	}

	// 绑定队列到交换机
	err = ch.QueueBind(
		cfg.Queue,                        // 队列名称
		constant.RabbitMQRoutingKeyVideo, // routing key
		cfg.Exchange,                     // 交换机名称
		false,                            // 不等待
		nil,                              // 额外参数
	)
	if err != nil {
		global.Logger.Fatal("Failed to bind RabbitMQ queue: " + err.Error())
	}

	global.Logger.Info("RabbitMQ connection and channel initialized successfully")
	return conn, ch
}
