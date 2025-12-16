package initialize

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/wangn-tech/tiny-douyin/internal/config"
	"github.com/wangn-tech/tiny-douyin/internal/global"
)

// InitMinIO 初始化 MinIO 客户端
func InitMinIO(cfg *config.MinIOConfig) *minio.Client {
	// 创建 MinIO 客户端
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		global.Logger.Fatal("Failed to create MinIO client: " + err.Error())
	}

	// 检查存储桶是否存在，不存在则创建
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		global.Logger.Fatal("Failed to check MinIO bucket: " + err.Error())
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{
			Region: cfg.Location,
		})
		if err != nil {
			global.Logger.Fatal("Failed to create MinIO bucket: " + err.Error())
		}
		global.Logger.Info("MinIO bucket created: " + cfg.BucketName)
	}

	// 设置存储桶策略为公开读
	policy := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::` + cfg.BucketName + `/*"]
			}
		]
	}`
	err = minioClient.SetBucketPolicy(ctx, cfg.BucketName, policy)
	if err != nil {
		global.Logger.Warn("Failed to set MinIO bucket policy: " + err.Error())
	}

	global.Logger.Info("MinIO client initialized successfully")
	return minioClient
}
