package constant

// Redis 键前缀
const (
	// RedisKeyUserPrefix 用户信息缓存键前缀
	RedisKeyUserPrefix = "user:"
	// RedisKeyVideoPrefix 视频信息缓存键前缀
	RedisKeyVideoPrefix = "video:"
	// RedisKeyUserVideosPrefix 用户视频列表缓存键前缀
	RedisKeyUserVideosPrefix = "user:videos:"
)

// RabbitMQ 常量
const (
	// RabbitMQRoutingKeyVideo 视频上传任务的 routing key
	RabbitMQRoutingKeyVideo = "video"
	// RabbitMQConsumerTag 消费者标签
	RabbitMQConsumerTag = "tiny-douyin-video-worker"
)

// MinIO 常量
const (
	// MinIOVideoPathFormat 视频对象路径格式: videos/{user_id}/{date}/{uuid}{ext}
	MinIOVideoPathFormat = "videos/%d/%s/%s"
	// MinioCoverPathFormat 封面对象路径格式: covers/{user_id}/{date}/{uuid}.jpg
	MinioCoverPathFormat = "covers/%d/%s/%s"
	// MinIODateFormat MinIO 存储路径中的日期格式
	MinIODateFormat = "2006-01-02"
	// MinioCoverExtension 封面文件扩展名
	MinioCoverExtension = ".jpg"
)

// 文件上传常量
const (
	// TempUploadDir 临时文件上传目录
	TempUploadDir = "./tmp/uploads"
	// MaxVideoSize 最大视频文件大小 (100MB)
	MaxVideoSize = 100 * 1024 * 1024
)

// 视频相关常量
const (
	// VideoStatusUploading 视频上传中
	VideoStatusUploading = "uploading"
	// VideoStatusReady 视频已就绪
	VideoStatusReady = "ready"
	// VideoStatusFailed 视频上传失败
	VideoStatusFailed = "failed"
)

// 日志相关常量
const (
	// LogDir 日志输出目录
	LogDir = "./tmp/logs"
	// LogFileName 日志文件名
	LogFileName = "tiny-douyin.log"
)

// 点赞操作类型
const (
	// FavoriteActionLike 点赞
	FavoriteActionLike = 1
	// FavoriteActionUnlike 取消点赞
	FavoriteActionUnlike = 2
)

// 评论操作类型
const (
	// CommentActionPublish 发布评论
	CommentActionPublish = 1
	// CommentActionDelete 删除评论
	CommentActionDelete = 2
)

// 评论内容限制
const (
	// CommentMaxLength 评论内容最大长度
	CommentMaxLength = 255
)
