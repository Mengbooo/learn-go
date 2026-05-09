# RAG 后端项目蓝图

## 1. 项目定位

这个项目不是一个“调一下 LLM API 然后返回答案”的玩具 demo，而是一个结构完整、适合展示后端能力的 RAG 系统。

你要把它做成一个能体现下面能力的作品：

- 文件上传与解析
- 异步任务处理
- 向量索引与检索
- 多租户或多用户数据隔离
- 鉴权与权限控制
- 缓存与性能优化
- 可观测性与部署

如果做得完整，它会比普通 CRUD 项目更能体现你的系统设计能力。

## 2. 项目目标

### 核心业务目标

用户可以上传文档，系统会自动解析、切分、向量化、建立索引。之后用户可以在知识库上提问，系统返回带引用来源的回答。

### 项目价值

这个项目很好，因为它天然覆盖了多个真实后端能力点：

- 文档上传和对象存储
- 异步流水线
- 数据建模
- 检索链路
- AI 服务调用
- 成本与限流控制

## 3. 最小 MVP 范围

建议第一版只做下面这些：

- 用户注册登录
- 创建知识库
- 上传 PDF / Markdown / TXT
- 文档解析与切分
- 向量化并写入索引
- 对知识库提问
- 返回答案和引用片段
- 查询任务状态

不要第一天就做：

- 多模型切换
- 工作流编排
- 团队空间
- 权限矩阵
- 复杂 reranking
- 多模态

先做通一条完整链路，再谈扩展。

## 4. 推荐技术栈

### 后端主栈

```txt
语言：Go
Web：chi
数据库：PostgreSQL
向量：pgvector 或 Qdrant
缓存：Redis
异步任务：Asynq
对象存储：MinIO / S3
日志：slog
指标：Prometheus
追踪：OpenTelemetry
部署：Docker Compose
```

### 为什么这样选

- `chi` 足够轻，便于你自己掌控结构。
- `PostgreSQL` 可以承载用户、知识库、文档、任务等核心业务数据。
- `pgvector` 适合第一版统一技术栈；如果你后面更强调向量能力，可切到 Qdrant。
- `Asynq` 很适合处理文档解析、切分、embedding、索引构建这类后台任务。

## 5. 系统模块划分

建议拆成下面几个模块：

- `auth`：注册登录、token、权限
- `knowledge`：知识库管理
- `document`：文档上传、元数据、状态
- `ingestion`：解析、切分、embedding、索引
- `retrieval`：召回、组装上下文
- `chat`：提问、回答、引用
- `billing_or_quota`：可选，控制调用配额
- `platform`：数据库、Redis、存储、日志、配置

如果是单体项目，也要按模块化单体来组织。

## 6. 核心数据模型

### users

```txt
id
email
password_hash
role
created_at
updated_at
```

### knowledge_bases

```txt
id
user_id
name
description
visibility
created_at
updated_at
```

### documents

```txt
id
knowledge_base_id
user_id
file_name
file_type
storage_key
status
error_message
created_at
updated_at
```

### document_chunks

```txt
id
document_id
chunk_index
content
token_count
metadata_json
created_at
```

### embeddings

如果使用 `pgvector`：

```txt
id
chunk_id
embedding vector
created_at
```

### chat_sessions

```txt
id
user_id
knowledge_base_id
title
created_at
updated_at
```

### chat_messages

```txt
id
session_id
role
content
citations_json
created_at
```

### ingestion_tasks

```txt
id
document_id
task_type
status
retry_count
error_message
created_at
updated_at
```

## 7. 关键状态机

### 文档状态

```txt
uploaded
parsing
chunking
embedding
indexed
failed
```

### 任务状态

```txt
pending
processing
completed
failed
retrying
```

你必须把状态设计清楚，否则异步链路后面一定会难排查。

## 8. 核心业务流程

### 文档导入流程

```txt
用户上传文件
  ↓
写 documents 记录
  ↓
文件保存到 MinIO / S3
  ↓
创建 ingestion task
  ↓
worker 拉取任务
  ↓
解析文档文本
  ↓
按策略切分 chunk
  ↓
调用 embedding 服务
  ↓
写入向量索引
  ↓
更新状态为 indexed
```

### 提问流程

```txt
用户提交问题
  ↓
对问题做 embedding
  ↓
从向量库召回相关 chunk
  ↓
可选 rerank
  ↓
拼接 prompt
  ↓
调用 LLM
  ↓
返回答案 + 引用片段
  ↓
记录 chat message
```

## 9. API 设计建议

### 认证

```txt
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/refresh
```

### 知识库

```txt
GET    /api/v1/knowledge-bases
POST   /api/v1/knowledge-bases
GET    /api/v1/knowledge-bases/{id}
PATCH  /api/v1/knowledge-bases/{id}
DELETE /api/v1/knowledge-bases/{id}
```

### 文档

```txt
POST   /api/v1/knowledge-bases/{id}/documents
GET    /api/v1/knowledge-bases/{id}/documents
GET    /api/v1/documents/{id}
GET    /api/v1/documents/{id}/status
DELETE /api/v1/documents/{id}
```

### 对话与问答

```txt
POST   /api/v1/knowledge-bases/{id}/ask
POST   /api/v1/chat-sessions
GET    /api/v1/chat-sessions/{id}
GET    /api/v1/chat-sessions/{id}/messages
```

### 任务

```txt
GET /api/v1/tasks/{id}
```

## 10. 向量存储怎么选

### 方案 A：PostgreSQL + pgvector

优点：

- 技术栈统一
- 上手简单
- 适合 MVP

缺点：

- 在大规模向量检索场景下不一定是最强方案

适合：

- 学习项目
- 小中型知识库

### 方案 B：Qdrant

优点：

- 专注向量检索
- 检索能力更强

缺点：

- 系统复杂度更高

适合：

- 你已经把第一版做通，想进一步强化检索层

建议：

- 第一版优先 `pgvector`
- 第二版再考虑 Qdrant

## 11. Chunking 和 Retrieval 的关键点

很多 RAG 项目做不好的根因，不在模型，而在前处理和检索策略。

### Chunking

你需要考虑：

- 按字符切，还是按语义段落切
- chunk 大小多大
- overlap 多大
- 是否保留标题层级和页码信息

建议第一版：

- 先用稳定、朴素的切分策略
- 元数据里记录文档名、页码、段落位置

### Retrieval

你需要考虑：

- top-k 召回多少
- 是否按知识库范围隔离
- 是否按文档类型过滤
- 是否做 rerank

建议第一版：

- 先做向量召回 + 元数据过滤
- rerank 放到第二阶段

## 12. 安全与权限

这个项目很容易忽视安全，但它其实很重要。

你至少要做：

- 用户只能访问自己的知识库和文档
- 上传文件大小和类型限制
- 文档下载需要权限控制
- API 鉴权
- 提问频率限制
- 第三方模型调用配额控制

如果不做这些，这个项目很快就会从“作品”退化成“demo”。

## 13. 异步任务设计建议

建议把导入链路拆成几个明确任务：

- `parse_document`
- `chunk_document`
- `embed_chunks`
- `index_embeddings`

也可以先做成一个复合任务，但内部仍然要记录阶段状态。

为什么要这样设计：

- 更容易排查失败点
- 更容易做重试
- 更容易扩展不同文档处理器

## 14. 缓存设计建议

这个项目适合缓存的地方：

- 知识库列表
- 文档状态摘要
- 高频问答结果短期缓存
- 用户配额信息

不建议一开始缓存：

- 所有检索结果
- 所有对话上下文
- 所有 chunk 数据

先保证正确，再做热点优化。

## 15. 可观测性设计

你至少应该记录：

### 日志

- request_id
- user_id
- knowledge_base_id
- document_id
- task_id
- model_name
- prompt_tokens
- completion_tokens
- latency
- error_code

### 指标

- 上传文件数量
- 文档处理成功率
- embedding 调用耗时
- 检索耗时
- LLM 响应耗时
- 问答错误率
- 队列长度
- 任务重试次数

### tracing

重点链路：

- 上传到索引完成
- 提问到回答返回
- 外部模型调用

## 16. 项目目录建议

```txt
rag-backend/
  cmd/
    api/
    worker/
  internal/
    app/
    config/
    domain/
    service/
    repository/
    transport/
      http/
    middleware/
    platform/
      database/
      redis/
      storage/
      logger/
      vector/
      llm/
  migrations/
  deployments/
  docs/
  scripts/
  Makefile
```

## 17. 开发阶段拆分建议

### 阶段 1：MVP

目标：

- 用户登录
- 知识库创建
- 文档上传
- 异步导入
- 简单提问

### 阶段 2：工程化

目标：

- 完整分层
- migration
- Docker Compose
- 基础测试
- 结构化日志

### 阶段 3：生产化增强

目标：

- 限流
- 配额
- 指标与 tracing
- 更完整错误码
- 文档和部署说明

### 阶段 4：效果增强

目标：

- chunk 优化
- rerank
- 多模型支持
- 检索策略优化

## 18. 面试或作品表达重点

这个项目完成后，你应该能清楚讲这些点：

- 为什么选 Go 做这个系统
- 为什么检索和导入要异步
- 为什么第一版选 `pgvector` 而不是更重的向量库
- 文档状态机如何设计
- 如何保证用户数据隔离
- 怎么做失败重试和幂等
- 怎么定位某次问答为什么很慢

如果你能把这些讲清楚，这个项目的价值会远高于一个普通 CRUD 项目。

## 19. 一句话建议

把这个项目当成“一个完整后端系统”，而不是“一个带 AI 调用的接口 demo”。只要你按这个思路做，它会非常适合作为你的代表项目。
