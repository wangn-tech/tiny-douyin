# Tiny-Douyin 项目设计文档

## 1. 项目概述

Tiny-Douyin 是一个简化版的短视频社交平台，基于字节跳动青训营后端大作业需求开发。

### 技术栈
- **语言**: Go 1.24+
- **框架**: Gin
- **数据库**: MySQL 8.0
- **缓存**: Redis
- **消息队列**: RabbitMQ
- **对象存储**: MinIO
- **ORM**: GORM
- **认证**: JWT
- **容器化**: Docker

## 2. 数据模型设计

### 设计原则
- 使用 GORM 内置的 `gorm.Model` 管理时间戳和软删除
- 保持模型简洁，避免冗余字段
- 统计数据通过实时查询或缓存获取
- 使用唯一索引防止重复数据

### 核心模型

#### User (用户)
```go
type User struct {
    gorm.Model
    Username  string  // 用户名，唯一索引
    Password  string  // bcrypt 加密密码
    Nickname  string  // 昵称
    Avatar    string  // 头像URL
    Signature string  // 个性签名
}
```

#### Video (视频)
```go
type Video struct {
    gorm.Model
    AuthorID    uint    // 作者ID，索引
    PlayURL     string  // 播放地址
    CoverURL    string  // 封面地址
    Title       string  // 标题
    Description string  // 描述
}
```

#### Comment (评论)
```go
type Comment struct {
    gorm.Model
    VideoID uint    // 视频ID，索引
    UserID  uint    // 用户ID，索引
    Content string  // 评论内容
}
```

#### Favorite (点赞)
```go
type Favorite struct {
    ID      uint
    VideoID uint  // 唯一索引(UserID, VideoID)
    UserID  uint
}
```

#### Relation (关注)
```go
type Relation struct {
    ID         uint
    FollowerID uint  // 唯一索引(FollowerID, FolloweeID)
    FolloweeID uint
}
```

#### Message (私信)
```go
type Message struct {
    ID         uint
    SenderID   uint      // 索引
    ReceiverID uint      // 索引
    Content    string
    CreatedAt  time.Time
}
```

## 3. 项目结构

```
tiny-douyin/
├── cmd/                # 程序入口
├── config/             # 配置文件
├── docs/               # 文档
├── internal/           # 内部包
│   ├── common/         # 公共组件
│   │   ├── errc/       # 错误码
│   │   └── response/   # 统一响应
│   ├── config/         # 配置解析
│   ├── dao/            # 数据访问层
│   ├── handler/        # HTTP 处理器
│   ├── initialize/     # 初始化
│   ├── middleware/     # 中间件
│   ├── model/          # 数据模型
│   ├── pkg/            # 工具包
│   ├── router/         # 路由
│   └── service/        # 业务逻辑层
├── scripts/            # 脚本
└── docker-compose.yml  # Docker 配置
```

## 4. 响应格式设计

### 基础响应
```json
{
  "status_code": 0,
  "status_msg": "success"
}
```

### 带数据响应
```json
{
  "status_code": 0,
  "status_msg": "success",
  "data": {}
}
```

### 错误响应
```json
{
  "status_code": 1001,
  "status_msg": "用户不存在"
}
```

## 5. 开发计划

### Phase 1: 基础设施
- [x] 项目初始化
- [x] 数据模型设计
- [x] 响应格式设计
- [ ] 工具包开发
- [ ] 中间件开发

### Phase 2: 用户模块
- [ ] 用户注册
- [ ] 用户登录
- [ ] 获取用户信息

### Phase 3: 视频模块
- [ ] 视频上传
- [ ] 视频流
- [ ] 视频列表

### Phase 4: 互动模块
- [ ] 点赞
- [ ] 评论

### Phase 5: 社交模块
- [ ] 关注
- [ ] 粉丝列表
- [ ] 好友列表
- [ ] 消息

### Phase 6: 优化
- [ ] Redis 缓存
- [ ] RabbitMQ 异步处理
- [ ] 性能优化

## 6. Git 分支策略

- `main`: 生产分支
- `develop/basic`: 基础功能开发
- `develop/advance`: 高级功能开发
- `feature/*`: 功能分支

## 7. 接口文档

参考：
- https://s.apifox.cn/15216c90-4b29-46c3-b0e9-0a91a366bf68/api-63913539
- https://documenter.getpostman.com/view/20584759/2s93CHuuiQ
