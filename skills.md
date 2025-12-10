# tiny-douyin（极简版抖音后端）

这是一个参考并重写优化 x-tiktok 的极简版抖音后端，技术栈为 Go + Gin + GORM + MySQL + Redis。

## 表结构与模型设计

我们为核心领域设计以下数据表：用户（users）、视频（videos）、评论（comments）、点赞（favorites）、关系（relations，关注）、消息（messages）、标签（tags）与视频标签映射（video_tags）。设计目标：

- 高并发写入：避免使用强外键约束，改为通过索引与应用层保证一致性，降低锁与阻塞。
- 去重保证：在点赞与关注上使用复合唯一索引，避免重复数据。
- 查询效率：常用字段与时间戳建立索引；对需要审计的内容使用软删除。
- 字段长度合理：密码使用 bcrypt（char(60)），文本字段使用合适的 varchar 长度。

### 具体说明

- users
  - username 唯一，用于登录；password 存储 bcrypt 哈希（char(60)）。
  - avatar、signature、nickname 等为可选资料字段。
- videos
  - author_id 建索引，包含播放与封面地址；可选标题与描述。
  - visibility 预留隐私控制（公开/私密/好友可见）。
- comments
  - video_id、user_id 建索引；使用软删除，便于内容治理。
- favorites（点赞/喜欢）
  - （user_id, video_id）唯一，保证同一用户对同一视频只能点赞一次。
- relations（关注关系）
  - （follower_id, followee_id）唯一，避免重复关注记录。
- messages（私信）
  - sender_id、receiver_id 建索引；按时间与双方查询。
- tags 与 video_tags
  - tag.name 唯一；video_tags 使用复合主键（video_id, tag_id）。

### 为什么这样设计？

- 在高并发场景下，避免数据库外键可降低写入冲突与迁移复杂度；一致性由服务层保证，并可用定期校验作补充。
- 点赞与关注的复合唯一索引是保证数据正确性的关键（去重且性能好）。
- 对内容类表（评论/视频/标签）保留软删除能力，便于审核与恢复。
- 对 created_at 等时间字段建索引，有利于时间序排序的 Feed 查询。

## 自动迁移与索引

`internal/initialize/database.go` 在启动时自动执行模型迁移（GORM AutoMigrate）。
- 复合唯一索引通过模型结构体标签创建（例如 favorites 的 `uniqueIndex:uk_fav_user_video`）。
- 无需手写 SQL 建索引，避免版本差异带来的语法不兼容问题。

## 快速上手

1. 通过 docker-compose 启动 MySQL 与 Redis。

```bash
docker-compose up -d
```

2. 配置文件 `config/config-dev.yaml` 中的数据库与 Redis 地址需与 compose 映射一致（例如 127.0.0.1:13306）。

3. 启动应用：

```bash
go run main.go
```

## Git 工作流建议

- 初始化与首次推送：

```bash
git init
git add .
git commit -m "chore: 初始化项目结构与模型"
git branch -M main
git remote add origin https://github.com/<your-org>/tiny-douyin.git
git push -u origin main
```

- 使用特性分支迭代开发：

```bash
git checkout -b feature/model-schema
git add internal/model internal/initialize/database.go
git commit -m "feat(model): 定义核心实体并使用 AutoMigrate 与唯一索引"
git push -u origin feature/model-schema
```

- 推荐使用约定式提交（Conventional Commits）：
  - feat(model): 新增用户/视频/评论/点赞/关系/消息/标签模型
  - fix(init): 修正 MySQL/Redis 初始化选项
  - docs(readme): 解释表设计与快速上手

- 里程碑打标签（发布版本）：

```bash
git tag v0.1.0
git push origin v0.1.0
```

## 下一步计划

- 在 DAO 与 Service 层实现事务与计数更新逻辑。
- 根据接口文档完善各模块的 Handler。
- 集成 zap 日志（已完成基础访问日志），按需扩展业务日志。

## Swagger（OpenAPI）集成

为便于接口自描述与联调，项目将集成 Swagger（基于 swaggo）。不改变现有结构，步骤如下：

1. 安装工具与依赖：
   - 开发机安装 swag CLI：
     ```bash
     go install github.com/swaggo/swag/cmd/swag@latest
     ```
   - 项目依赖：
     ```bash
     go get -u github.com/swaggo/gin-swagger
     go get -u github.com/swaggo/files
     ```

2. 在 `main.go` 顶部添加用于生成文档的注释（示例）：
   ```go
   // @title Tiny Douyin API
   // @version 0.1.0
   // @description 极简版抖音后端 API
   // @host localhost:8080
   // @BasePath /
   ```

3. 在 `internal/router/router.go` 中开启 Swagger 路由：
   ```go
   // r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
   ```

4. 生成与查看文档：
   ```bash
   swag init --parseDependency --parseInternal
   go run main.go
   # 浏览器打开 http://localhost:8080/swagger/index.html
   ```

后续将为每个 Handler 添加注释（@Summary、@Description、@Param、@Success 等）以自动生成路由文档。
