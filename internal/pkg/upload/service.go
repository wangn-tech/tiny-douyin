package upload

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"github.com/wangn-tech/tiny-douyin/internal/common/constant"
	"github.com/wangn-tech/tiny-douyin/internal/global"
)

// VideoUploadTask 视频上传任务
type VideoUploadTask struct {
	VideoID     uint   `json:"video_id"`     // 视频 ID
	VideoPath   string `json:"video_path"`   // 本地视频文件路径
	CoverPath   string `json:"cover_path"`   // 本地封面文件路径（可选）
	VideoName   string `json:"video_name"`   // MinIO 中的视频文件名
	CoverName   string `json:"cover_name"`   // MinIO 中的封面文件名
	ContentType string `json:"content_type"` // 视频内容类型
	UserID      uint   `json:"user_id"`      // 上传用户 ID
	Title       string `json:"title"`        // 视频标题
	Description string `json:"description"`  // 视频描述
}

// UploadService 上传服务
type UploadService struct {
	minioClient *minio.Client
	rabbitChan  *amqp.Channel
	bucketName  string
	urlPrefix   string
	exchange    string
	routingKey  string
}

// NewUploadService 创建上传服务
func NewUploadService() IUploadService {
	return &UploadService{
		minioClient: global.MinIOClient,
		rabbitChan:  global.RabbitChan,
		bucketName:  global.Config.MinIO.BucketName,
		urlPrefix:   global.Config.MinIO.URLPrefix,
		exchange:    global.Config.RabbitMQ.Exchange,
		routingKey:  constant.RabbitMQRoutingKeyVideo,
	}
}

// SaveTempFile 保存上传的文件到临时目录
func (s *UploadService) SaveTempFile(fileData []byte, ext string) (string, error) {
	// 创建临时目录
	tempDir := constant.TempUploadDir
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// 生成唯一文件名
	filename := uuid.New().String() + ext
	filepath := filepath.Join(tempDir, filename)

	// 写入文件
	if err := os.WriteFile(filepath, fileData, 0644); err != nil {
		return "", fmt.Errorf("failed to write temp file: %w", err)
	}

	return filepath, nil
}

// PublishUploadTask 发布上传任务到消息队列
func (s *UploadService) PublishUploadTask(ctx context.Context, task *VideoUploadTask) error {
	// 序列化任务
	body, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	// 发布消息
	err = s.rabbitChan.PublishWithContext(
		ctx,
		s.exchange,   // 交换机
		s.routingKey, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // 持久化消息
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	global.Logger.Info("Video upload task published to queue",
		zap.Uint("video_id", task.VideoID),
		zap.Uint("user_id", task.UserID))

	return nil
}

// UploadToMinIO 上传文件到 MinIO
func (s *UploadService) UploadToMinIO(ctx context.Context, filePath, objectName, contentType string) (string, error) {
	// 上传文件
	_, err := s.minioClient.FPutObject(ctx, s.bucketName, objectName, filePath, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to MinIO: %w", err)
	}

	// 返回文件 URL
	url := fmt.Sprintf("%s/%s", s.urlPrefix, objectName)
	global.Logger.Info("File uploaded to MinIO",
		zap.String("object_name", objectName),
		zap.String("url", url))

	return url, nil
}

// CleanupTempFile 清理临时文件
func (s *UploadService) CleanupTempFile(filePath string) {
	if err := os.Remove(filePath); err != nil {
		global.Logger.Warn("Failed to cleanup temp file",
			zap.String("file_path", filePath),
			zap.String("error", err.Error()))
	}
}

// GenerateObjectName 生成 MinIO 中的对象名称
func (s *UploadService) GenerateObjectName(userID uint, ext string) string {
	// 格式: videos/{user_id}/{date}/{uuid}.{ext}
	date := time.Now().Format(constant.MinIODateFormat)
	filename := uuid.New().String() + ext
	return fmt.Sprintf(constant.MinIOVideoPathFormat, userID, date, filename)
}

// GenerateCoverObjectName 生成封面对象名称
func (s *UploadService) GenerateCoverObjectName(userID uint) string {
	// 格式: covers/{user_id}/{date}/{uuid}.jpg
	date := time.Now().Format(constant.MinIODateFormat)
	filename := uuid.New().String() + constant.MinioCoverExtension
	return fmt.Sprintf(constant.MinioCoverPathFormat, userID, date, filename)
}
