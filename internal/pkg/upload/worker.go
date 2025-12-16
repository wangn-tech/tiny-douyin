package upload

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"github.com/wangn-tech/tiny-douyin/internal/common/constant"
	"github.com/wangn-tech/tiny-douyin/internal/dao"
	"github.com/wangn-tech/tiny-douyin/internal/global"
)

// Worker 上传任务工作器
type Worker struct {
	uploadService IUploadService
	videoDAO      dao.IVideoDAO
	queueName     string
}

// NewWorker 创建工作器（通过依赖注入）
func NewWorker(uploadService IUploadService, videoDAO dao.IVideoDAO) IUploadWorker {
	return &Worker{
		uploadService: uploadService,
		videoDAO:      videoDAO,
		queueName:     global.Config.RabbitMQ.Queue,
	}
}

// Start 启动工作器
func (w *Worker) Start(ctx context.Context) error {
	// 获取消息通道
	msgs, err := global.RabbitChan.Consume(
		w.queueName,                  // 队列名称
		constant.RabbitMQConsumerTag, // 消费者标签
		false,                        // 自动确认
		false,                        // 独占
		false,                        // 不等待
		false,                        // 额外参数
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	global.Logger.Info("Upload worker started, waiting for messages...")

	// 处理消息
	go func() {
		for {
			select {
			case <-ctx.Done():
				global.Logger.Info("Upload worker stopped")
				return
			case msg, ok := <-msgs:
				if !ok {
					global.Logger.Warn("Message channel closed")
					return
				}
				w.processMessage(msg)
			}
		}
	}()

	return nil
}

// processMessage 处理单条消息
func (w *Worker) processMessage(msg amqp.Delivery) {
	ctx := context.Background()

	// 解析任务
	var task VideoUploadTask
	if err := json.Unmarshal(msg.Body, &task); err != nil {
		global.Logger.Error("Failed to unmarshal task",
			zap.Error(err))
		_ = msg.Nack(false, false) // 拒绝消息，不重新入队
		return
	}

	global.Logger.Info("Processing upload task",
		zap.Uint("video_id", task.VideoID),
		zap.String("video_path", task.VideoPath))

	// 上传视频到 MinIO
	videoURL, err := w.uploadService.UploadToMinIO(ctx, task.VideoPath, task.VideoName, task.ContentType)
	if err != nil {
		global.Logger.Error("Failed to upload video to MinIO",
			zap.Uint("video_id", task.VideoID),
			zap.Error(err))
		_ = msg.Nack(false, true) // 拒绝消息，重新入队
		return
	}

	// 上传封面到 MinIO（如果存在）
	var coverURL string
	if task.CoverPath != "" {
		coverURL, err = w.uploadService.UploadToMinIO(ctx, task.CoverPath, task.CoverName, "image/jpeg")
		if err != nil {
			global.Logger.Error("Failed to upload cover to MinIO",
				zap.Uint("video_id", task.VideoID),
				zap.Error(err))
			// 封面上传失败不影响视频，继续处理
		}
	}

	// 如果没有封面，生成默认封面 URL
	if coverURL == "" {
		coverURL = videoURL + "?x-oss-process=video/snapshot,t_1000,f_jpg" // 1秒处截图
	}

	// 更新数据库中的视频 URL
	video, err := w.videoDAO.GetVideoByID(ctx, task.VideoID)
	if err != nil {
		global.Logger.Error("Failed to get video from database",
			zap.Uint("video_id", task.VideoID),
			zap.Error(err))
		_ = msg.Nack(false, true) // 拒绝消息，重新入队
		return
	}

	video.PlayURL = videoURL
	video.CoverURL = coverURL
	video.UpdatedAt = time.Now()

	if err := w.videoDAO.UpdateVideo(ctx, video); err != nil {
		global.Logger.Error("Failed to update video in database",
			zap.Uint("video_id", task.VideoID),
			zap.Error(err))
		_ = msg.Nack(false, true) // 拒绝消息，重新入队
		return
	}

	// 清理临时文件
	w.uploadService.CleanupTempFile(task.VideoPath)
	if task.CoverPath != "" {
		w.uploadService.CleanupTempFile(task.CoverPath)
	}

	// 确认消息
	if err := msg.Ack(false); err != nil {
		global.Logger.Error("Failed to acknowledge message",
			zap.Error(err))
	}

	global.Logger.Info("Upload task completed successfully",
		zap.Uint("video_id", task.VideoID),
		zap.String("video_url", videoURL),
		zap.String("cover_url", coverURL))
}
