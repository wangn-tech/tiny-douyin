package upload

import "context"

// IUploadService 上传服务接口
type IUploadService interface {
	// SaveTempFile 保存上传的文件到临时目录
	SaveTempFile(fileData []byte, ext string) (string, error)
	// PublishUploadTask 发布上传任务到消息队列
	PublishUploadTask(ctx context.Context, task *VideoUploadTask) error
	// UploadToMinIO 上传文件到 MinIO
	UploadToMinIO(ctx context.Context, filePath, objectName, contentType string) (string, error)
	// CleanupTempFile 清理临时文件
	CleanupTempFile(filePath string)
	// GenerateObjectName 生成 MinIO 中的对象名称
	GenerateObjectName(userID uint, ext string) string
	// GenerateCoverObjectName 生成封面对象名称
	GenerateCoverObjectName(userID uint) string
}

// IUploadWorker 上传任务工作器接口
type IUploadWorker interface {
	// Start 启动工作器
	Start(ctx context.Context) error
}
