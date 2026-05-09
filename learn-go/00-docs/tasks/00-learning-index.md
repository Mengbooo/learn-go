# Go 后端学习资料索引

## 1. 这份索引的作用

你现在已经有两类材料：

- 路线图：告诉你整体方向是什么。
- 任务文档：告诉你每一阶段要做什么。

这次新增的是第三类材料：

- 教程文档：告诉你知识本身是什么、为什么存在、应该怎么理解。

你可以把三类文档理解为：

```txt
路线图 = 全局地图
任务文档 = 每一关的执行清单
教程文档 = 每一关的知识讲解
工具生态文档 = 常用工具的定位说明
项目蓝图文档 = 把知识真正串成作品
```

## 2. 推荐阅读顺序

建议按下面顺序使用文档：

1. 先看 `roadmap/raodmap.md`，理解全局路线。
2. 进入某一阶段时，先看对应任务文档。
3. 做任务前或卡住时，再看对应教程文档。
4. 选技术栈或工具时，查工具生态文档。
5. 想做完整项目时，直接参考项目蓝图文档。

## 3. 任务文档与教程文档映射

### 第一组：语言、工具链、HTTP

- 任务：[01-go-language-basics.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/01-go-language-basics.md)
- 教程：[tutorials/01-go-language-tutorial.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/tutorials/01-go-language-tutorial.md)

- 任务：[02-go-toolchain-and-standards.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/02-go-toolchain-and-standards.md)
- 教程：[tutorials/02-go-toolchain-and-stdlib-tutorial.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/tutorials/02-go-toolchain-and-stdlib-tutorial.md)

- 任务：[03-http-rest-web-service.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/03-http-rest-web-service.md)
- 教程：[tutorials/03-http-and-api-tutorial.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/tutorials/03-http-and-api-tutorial.md)

### 第二组：数据库、架构、安全

- 任务：[04-database-sql-transaction-index.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/04-database-sql-transaction-index.md)
- 任务：[05-project-structure-and-modeling.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/05-project-structure-and-modeling.md)
- 教程：[tutorials/04-database-and-architecture-tutorial.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/tutorials/04-database-and-architecture-tutorial.md)

- 任务：[06-auth-authorization-security.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/06-auth-authorization-security.md)
- 教程：[tutorials/05-security-cache-async-tutorial.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/tutorials/05-security-cache-async-tutorial.md)

### 第三组：缓存、异步、质量、部署、可观测性、架构

- 任务：[07-cache-redis-performance.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/tasks/07-cache-redis-performance.md)
- 教程：[tutorials/07-cache-redis-performance-tutorial.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/tutorials/07-cache-redis-performance-tutorial.md)

- 任务：[08-async-queue-idempotency.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/tasks/08-async-queue-idempotency.md)
- 教程：[tutorials/08-async-queue-idempotency-tutorial.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/tutorials/08-async-queue-idempotency-tutorial.md)

- 任务：[09-testing-docs-code-quality.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/tasks/09-testing-docs-code-quality.md)
- 教程：[tutorials/09-testing-docs-code-quality-tutorial.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/tutorials/09-testing-docs-code-quality-tutorial.md)

- 任务：[10-deployment-docker-linux-cicd.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/tasks/10-deployment-docker-linux-cicd.md)
- 教程：[tutorials/10-deployment-docker-linux-cicd-tutorial.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/tutorials/10-deployment-docker-linux-cicd-tutorial.md)

- 任务：[11-observability-and-troubleshooting.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/tasks/11-observability-and-troubleshooting.md)
- 教程：[tutorials/11-observability-and-troubleshooting-tutorial.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/tutorials/11-observability-and-troubleshooting-tutorial.md)

- 任务：[12-microservices-distributed-system-design.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/tasks/12-microservices-distributed-system-design.md)
- 教程：[tutorials/12-microservices-distributed-system-design-tutorial.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/tutorials/12-microservices-distributed-system-design-tutorial.md)

## 4. 工具文档

- 工具生态说明：[tools/go-backend-tooling-ecosystem.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/tools/go-backend-tooling-ecosystem.md)

这份文档重点回答这些问题：

- `chi`、`Gin`、`Fiber` 分别适合什么阶段。
- `pgx`、`sqlc`、`GORM` 各自的生态位是什么。
- `Redis`、`Asynq`、`RabbitMQ`、`Kafka` 应该怎么选。
- `Prometheus`、`Grafana`、`OpenTelemetry`、`Sentry` 分别负责什么。

## 5. 项目蓝图文档

- RAG 项目蓝图：[projects/rag-backend-project-blueprint.md](/Users/qiumengbo.123/Desktop/basic-go/learn-go/00-docs/projects/rag-backend-project-blueprint.md)

这份文档不是一个玩具 demo，而是一份偏生产思维的项目设计稿，覆盖：

- 项目目标
- 技术栈建议
- API 设计
- 数据库与向量检索设计
- 异步任务与索引流程
- 权限、安全、可观测性
- MVP 到生产版的演进路径

## 6. 使用建议

如果你现在准备正式推进学习，建议采用下面节奏：

```txt
每周先选 1 个任务文档
  ↓
先看对应教程，理解知识背景
  ↓
再做任务文档中的练习和项目任务
  ↓
遇到“工具怎么选”的问题就查工具生态文档
  ↓
学到第 5-8 阶段时开始推进 RAG 项目蓝图
```
