# tiny-douyin

简化版抖音后端服务，基于 Go + Gin + GORM + MySQL + Redis + MinIO + RabbitMQ 实现。

## ✨ 功能特性

### 已实现功能

- ✅ 用户模块
  - 用户注册
  - 用户登录
  - 获取用户信息
  - JWT 认证

- ✅ 视频模块
  - 视频发布（异步上传）
  - 视频流推荐
  - 用户视频列表
  - MinIO 对象存储
  - RabbitMQ 异步处理

### 技术栈

- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0
- **缓存**: Redis 7
- **对象存储**: MinIO
- **消息队列**: RabbitMQ
- **依赖注入**: Google Wire
- **日志**: Zap
- **认证**: JWT
- **密码加密**: Bcrypt

## 🚀 快速开始

### 前置要求

- Go 1.24+
- Docker & Docker Compose

### 1. 启动依赖服务

```bash
# 启动 MySQL, Redis, MinIO, RabbitMQ
./start-services.sh

# 或手动启动
docker-compose up -d
```

### 2. 配置

编辑 `config/config-dev.yaml` 文件，配置数据库、Redis、MinIO、RabbitMQ 连接信息。

### 3. 运行服务

```bash
# 编译
go build -o tiny-douyin

# 运行
./tiny-douyin --env=dev
```

服务将在 `http://localhost:8080` 启动。

### 4. 测试

```bash
# 测试视频上传功能
./test-video-upload.sh
```

## 📁 项目结构

```
tiny-douyin/
├── config/                  # 配置文件
│   └── config-dev.yaml     # 开发环境配置
├── internal/               # 内部代码
│   ├── api/               # API 层
│   │   ├── dto/          # 数据传输对象
│   │   └── handler/      # 请求处理器
│   ├── common/           # 公共模块
│   │   ├── constant/    # 常量定义
│   │   ├── errc/        # 错误码
│   │   └── response/    # 统一响应
│   ├── service/          # 业务逻辑层
│   ├── dao/              # 数据访问层
│   ├── model/            # 数据模型
│   ├── middleware/       # 中间件
│   ├── router/           # 路由配置
│   ├── pkg/              # 工具包
│   │   ├── hash/        # 密码加密
│   │   ├── jwt/         # JWT 认证
│   │   ├── upload/      # 文件上传服务 (接口化)
│   │   └── validator/   # 参数验证
│   ├── config/          # 配置管理
│   ├── global/          # 全局变量
│   ├── initialize/      # 初始化
│   └── wire/            # 依赖注入
├── tmp/                 # 临时文件目录 (git ignored)
│   ├── logs/           # 应用日志
│   └── uploads/        # 上传临时文件
├── docs/                # 文档
├── public/              # 静态文件
├── docker-compose.yml   # Docker Compose 配置
├── check-status.sh      # 状态检查脚本
└── test-video-upload.sh # 测试脚本
```

## 🎬 视频上传架构

### 异步上传流程

```
用户上传视频 → Handler 接收 → 保存临时文件 → 创建记录 
    → 发布到消息队列 → 立即返回成功
    
(后台异步)
Worker 消费任务 → 上传到 MinIO → 生成封面 → 更新数据库 → 清理临时文件
```

### 优势

- **快速响应**: 用户无需等待上传完成
- **高并发**: 消息队列缓冲请求
- **可靠性**: 消息持久化和重试机制
- **可扩展**: 可启动多个 Worker

详细文档: [MinIO 和 RabbitMQ 集成文档](./docs/minio-rabbitmq-integration.md)

## 📡 API 文档

### 用户接口

#### 注册
```bash
POST /douyin/user/register/
Content-Type: application/json

{
  "username": "testuser",
  "password": "123456"
}
```

#### 登录
```bash
POST /douyin/user/login/
Content-Type: application/json

{
  "username": "testuser",
  "password": "123456"
}
```

#### 获取用户信息
```bash
GET /douyin/user/?user_id=1
Authorization: Bearer {token}
```

### 视频接口

#### 发布视频
```bash
POST /douyin/publish/action/
Authorization: Bearer {token}
Content-Type: multipart/form-data

data: <video_file>
title: 视频标题（可选）
```

#### 获取视频流
```bash
GET /douyin/feed/?latest_time=1702742400
```

#### 获取用户视频列表
```bash
GET /douyin/publish/list/?user_id=1
Authorization: Bearer {token}
```

## 🛠 开发指南

### 添加新功能

1. 定义 DTO (internal/api/dto/)
2. 实现 DAO (internal/dao/)
3. 实现 Service (internal/service/)
4. 实现 Handler (internal/api/handler/)
5. 注册路由 (internal/router/)
6. 配置依赖注入 (internal/wire/)

### Wire 依赖注入

```bash
# 生成依赖注入代码
cd internal/wire
wire
```

### 数据库迁移

数据库表会在应用启动时自动创建（使用 GORM AutoMigrate）。

## 🔧 配置说明

### MinIO 配置

```yaml
minio:
  endpoint: localhost:9000
  access_key_id: minioadmin
  secret_access_key: minioadmin
  use_ssl: false
  bucket_name: tiny-douyin
  location: us-east-1
  url_prefix: http://localhost:9000/tiny-douyin
```

访问 MinIO 控制台: http://localhost:9090

### RabbitMQ 配置

```yaml
rabbitmq:
  host: localhost
  port: 5672
  user: guest
  password: guest
  vhost: /
  exchange: tiny-douyin.video.upload
  queue: video.upload.queue
```

访问 RabbitMQ 管理界面: http://localhost:15672

## 📊 监控

### MinIO
- 控制台: http://localhost:9090
- 用户名/密码: minioadmin/minioadmin

### RabbitMQ
- 管理界面: http://localhost:15672
- 用户名/密码: guest/guest

### 应用日志

日志输出到标准输出，包含详细的请求和处理信息。

## 🐛 问题排查

### MinIO 连接失败
- 检查 MinIO 服务是否启动: `docker ps | grep minio`
- 查看 MinIO 日志: `docker logs minio`

### RabbitMQ 连接失败
- 检查 RabbitMQ 服务是否启动: `docker ps | grep rabbitmq`
- 查看 RabbitMQ 日志: `docker logs rabbitmq`

### 视频上传后无法访问
- 检查 MinIO 存储桶策略是否设置为公开读
- 查看应用日志确认上传是否成功

## 📚 相关文档

- [项目结构最佳实践](./docs/project-structure-best-practice.md)
- [DTO 最佳实践](./docs/dto-best-practice.md)
- [依赖注入](./docs/dependency-injection.md)
- [Wire 快速参考](./docs/wire-quick-reference.md)
- [MinIO 和 RabbitMQ 集成](./docs/minio-rabbitmq-integration.md)

## 🤝 参考项目

- [TikGok](https://github.com/CyanAsterisk/TikGok) - 高性能抖音后端实现
- [x-tiktok](https://github.com/X-Engineer/x-tiktok) - 微服务架构抖音后端

## 📝 License

MIT License
