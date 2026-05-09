# Go 后端工具与生态位说明

## 1. 这份文档解决什么问题

你现在缺的不是“工具名单”，而是“为什么要用这个工具，它和相邻工具的边界是什么”。这份文档就是用来解决选型混乱问题的。

这里的核心词是“生态位”：

- 它在整个系统里负责什么
- 它和同类工具的差异是什么
- 它适合你现在这个阶段吗

## 2. 先给一个推荐组合

如果你现在是一个人或小团队做 Go 后端，推荐从下面这套开始：

```txt
Web：chi
DB：PostgreSQL + pgx + sqlc
Migration：golang-migrate
Cache / Queue：Redis + Asynq
Auth：JWT + bcrypt
Config：env + config package
Logging：slog
Testing：testing + httptest
Deploy：Docker + docker-compose + GitHub Actions
Observability：Prometheus + Grafana + OpenTelemetry
Storage：MinIO / S3
```

这套不是唯一答案，但很适合认真学习和做作品。

## 3. Web 框架与路由层

### `net/http`

生态位：

- HTTP 服务底层标准库

适合：

- 学基础
- 简单项目
- 理解请求处理模型

### `chi`

生态位：

- 轻量路由和中间件框架

适合：

- 希望贴近标准库
- 想控制结构
- 中小型项目

### `Gin`

生态位：

- 高使用率通用 Web 框架

适合：

- 快速起步
- 想参考大量现成资料

### `Fiber`

生态位：

- 风格更自成体系的 Web 框架

适合：

- 后期横向了解生态

建议：

- 学习阶段优先 `net/http` + `chi`
- 如果你更看重资料密度，也可以选 `Gin`

## 4. 数据库访问层

### PostgreSQL

生态位：

- 主数据库

适合：

- 绝大多数内容型、业务型、管理型后端
- 后续扩展全文检索、JSON 字段、向量检索

### MySQL

生态位：

- 另一类主流关系数据库

适合：

- 团队已有经验
- 兼容已有业务环境

学习阶段更推荐先用 PostgreSQL。

### `pgx`

生态位：

- Go 访问 PostgreSQL 的主流库

适合：

- 要认真做 Go + PostgreSQL 项目

### `sqlc`

生态位：

- SQL 代码生成器

适合：

- 想保留 SQL 可见性
- 想减少重复样板代码
- 想要类型安全

### `GORM`

生态位：

- ORM

适合：

- 快速原型
- 简单 CRUD 管理后台

不适合：

- 你完全不会 SQL 却想跳过数据库认知

### `Bun`

生态位：

- 较现代的 Go ORM / 查询构建方案

适合：

- 你已经知道 ORM 与 SQL 的取舍，想找介于原生 SQL 和重 ORM 之间的方案

## 5. Migration 工具

### `golang-migrate`

生态位：

- 数据库结构变更管理工具

适合：

- 需要稳定执行 up/down
- 想把结构演进纳入版本控制

没有 migration 的项目，数据库结构通常会很快失控。

## 6. 缓存与队列

### Redis

生态位：

- 内存缓存
- 轻量任务中转
- 限流状态存储
- 会话和计数器

不是主数据库替代品。

### Asynq

生态位：

- 基于 Redis 的 Go 后台任务框架

适合：

- 邮件
- 图像处理
- AI 调用任务
- 报告生成

### RabbitMQ

生态位：

- 通用消息队列

适合：

- 明确消息投递和消费模型
- 传统业务消息系统

### Kafka

生态位：

- 高吞吐事件流平台

适合：

- 埋点流
- 日志流
- 事件平台

不建议小项目为了显得高级而上 Kafka。

## 7. 认证与安全

### `bcrypt`

生态位：

- 密码哈希

### `argon2`

生态位：

- 更现代的密码哈希选择之一

### JWT

生态位：

- 无状态认证令牌

适合：

- API 鉴权
- 前后端分离项目

### Session

生态位：

- 服务端状态化登录态

适合：

- 传统 Web 应用
- 强服务端控制场景

## 8. 配置与环境管理

### 环境变量

生态位：

- 最基本配置注入方式

适合：

- 数据库地址
- Redis 地址
- token 密钥
- 第三方 API Key

### `.env`

生态位：

- 本地开发环境变量管理辅助

注意：

- `.env.example` 可以提交
- 真正的 `.env` 和密钥文件不要进 Git

## 9. 日志与可观测性

### `slog`

生态位：

- Go 官方结构化日志

适合：

- 中小型项目默认日志方案

### Prometheus

生态位：

- 指标采集和存储

### Grafana

生态位：

- 面板可视化

### OpenTelemetry

生态位：

- tracing 和 telemetry 标准层

### Sentry

生态位：

- 错误聚合和报警

## 10. 部署工具

### Docker

生态位：

- 打包和运行环境标准化

### docker-compose

生态位：

- 本地多服务编排

### GitHub Actions

生态位：

- CI/CD 入口

### Nginx

生态位：

- 传统反向代理

### Caddy

生态位：

- 更轻量、更省配置的反向代理

对个人项目和小团队来说，Caddy 常常更顺手。

## 11. AI / RAG 相关工具生态位

### `pgvector`

生态位：

- PostgreSQL 内的向量检索扩展

适合：

- 希望把结构化数据和向量数据放在一套数据库体系里

### Qdrant

生态位：

- 专注向量检索的数据库

适合：

- 更明确的向量检索场景
- 需要更丰富的向量检索能力

### Milvus

生态位：

- 更偏大型向量检索系统

适合：

- 数据规模较大
- 检索能力要求更高

### MinIO / S3

生态位：

- 对象存储

适合：

- 原始文件
- 文档上传
- 切分前原文存储

## 12. 学习阶段的选型原则

如果你现在的目标是“学会做可上线的 Go 后端”，而不是“堆一堆热门词”，建议遵守下面几条：

- 同一层只选一套主线工具，不要同时学三种框架。
- 优先选概念清晰、资料稳定、社区成熟的工具。
- 先掌握标准库和底层原理，再上更高层封装。
- 能用 PostgreSQL + Redis + Docker 解决的问题，不要急着引入更重系统。
- 工具的价值在于降低问题成本，不在于让项目看起来更复杂。

## 13. 一套很稳的学习型生产栈

如果你问“我现在最应该用什么栈”，给你的建议是：

```txt
Go
chi
PostgreSQL
pgx
sqlc
golang-migrate
Redis
Asynq
slog
Prometheus
Grafana
Docker
GitHub Actions
MinIO
```

这套栈既适合学习，也足够支撑一个不错的作品项目。
