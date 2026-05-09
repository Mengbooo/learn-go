---
title: '路线图'
publishDate: '2026-05-08 12:00:00'
description: 'road to backend'
tags:
  - Backend
language: 'zh-CN'
draft: true
---


# 完全体后端通关秘籍：Go + AI 辅助学习路线图

> 版本：v1.0  
> 适用对象：有前端/Next.js API 基础，但缺少完整后端生产经验的人  
> 主语言：Go  
> 核心目标：不是“会写 Go”，而是“能用 Go 做出可上线、可维护、可排查、可扩展的后端系统”

---

## 0. 这份秘籍怎么用

你现在的情况很典型：

- 会写一些 Next.js API；
- 对后端整体生产链路不够熟；
- 想用 Go 补齐后端工程能力；
- 认可 AI 会降低“手写代码”的价值，但不会降低“理解模式、工具、边界和取舍”的价值。

所以这份路线不是传统“语法教程”，而是一份后端通关地图。你要沿着它完成三件事：

1. **掌握 Go 语言和工具链**：写出符合 Go 惯例的代码。
2. **掌握后端生产能力**：接口、数据库、认证、缓存、部署、监控、异步任务。
3. **掌握工程判断力**：知道某个需求该放哪一层、用什么工具、有哪些坑、上线后怎么排查。

你可以把学习目标理解成：

```txt
初级：能写 CRUD API
合格：能设计接口、表结构、认证、部署
生产：能处理日志、监控、限流、异步、故障排查
高阶：能设计系统边界、服务拆分、可用性、一致性和成本
```

---

## 1. 总路线图

完整后端能力可以分为 12 层：

```txt
第 01 层：Go 语言基础
第 02 层：Go 工具链与工程规范
第 03 层：HTTP / REST / Web 服务
第 04 层：数据库 / SQL / 事务 / 索引
第 05 层：项目结构 / 分层架构 / 业务建模
第 06 层：认证 / 鉴权 / 安全
第 07 层：缓存 / Redis / 性能优化
第 08 层：异步任务 / 消息队列 / 幂等
第 09 层：测试 / 文档 / 代码质量
第 10 层：部署 / Docker / Linux / CI/CD
第 11 层：日志 / 指标 / 链路追踪 / 可观测性
第 12 层：微服务 / 分布式 / 系统设计 / 架构能力
```

路线图可以这样看：

```txt
Go 语言基础
  ↓
标准库与工具链
  ↓
HTTP API
  ↓
数据库与事务
  ↓
认证鉴权与安全
  ↓
项目分层与工程化
  ↓
Docker 部署与 CI/CD
  ↓
日志、监控、链路追踪
  ↓
缓存、队列、异步任务
  ↓
微服务、分布式、系统设计
  ↓
真实项目沉淀 + 面试表达
```

---

## 2. 第一关：Go 语言基础

### 2.1 通关目标

这一关的目标不是“背语法”，而是能写出简洁、清晰、符合 Go 风格的代码。

你要掌握：

```txt
变量与常量
var 与 :=
零值
const 与 iota
作用域与变量遮蔽
基础数据类型
类型转换
数组、切片、map、结构体
if / switch / for / range
函数、多返回值、闭包
指针
方法与接口
错误处理
泛型基础
```

### 2.2 必学清单

#### 变量与常量

你需要理解：

- `var` 和 `:=` 的区别；
- 零值为什么重要；
- `const` 和 `iota` 怎么用于枚举；
- 变量遮蔽为什么容易制造 bug。

验收标准：

```txt
能解释 var、:=、const、iota 的使用场景
能识别变量遮蔽问题
能写出不依赖未初始化状态的代码
```

#### 数据类型

你需要掌握：

```txt
bool
int / int64 / uint
float32 / float64
complex
byte
rune
string
类型转换
```

重点不是类型名字，而是边界：

```txt
什么时候用 int64？
什么时候不能随便 float 比较？
string 和 []byte 的转换成本是什么？
rune 为什么和 Unicode 有关？
```

#### 复合类型

你需要掌握：

```txt
array
slice
map
struct
struct tag
JSON tag
嵌入结构体
comma-ok idiom
```

重点理解：

```txt
slice 的 len 和 cap
slice 扩容
map 读取不存在 key 的行为
struct tag 如何影响 JSON 编解码
map 并发读写为什么危险
```

#### 控制流

你需要掌握：

```txt
if
if-else
switch
for
for range
break
continue
goto 不推荐使用
```

Go 没有 while，你需要习惯用 `for` 表达循环。

#### 函数

你需要掌握：

```txt
普通函数
多返回值
可变参数
匿名函数
闭包
命名返回值
值传递
```

后端里最常见的是：

```go
result, err := service.DoSomething(ctx, input)
if err != nil {
    return err
}
```

多返回值和错误处理是 Go 后端的核心风格。

#### 指针

你需要理解：

```txt
指针是什么
什么时候需要指针
指针接收者 vs 值接收者
map / slice 本身的引用语义
结构体方法什么时候用 *T
```

不要把 Go 指针学成 C 指针。Go 的指针主要服务于：

```txt
避免大结构体拷贝
修改原对象
表达可选值
配合方法接收者
```

#### 方法与接口

你需要掌握：

```txt
method vs function
value receiver
pointer receiver
interface
empty interface / any
interface embedding
类型断言
type switch
```

最重要的意识：

> Go 的接口是隐式实现，不需要显式 `implements`。

这会影响你理解 repository、service、mock、依赖注入。

#### 泛型

你需要掌握：

```txt
为什么需要泛型
泛型函数
泛型类型
类型约束
类型推断
```

但不要一开始滥用泛型。Go 的泛型适合：

```txt
通用数据结构
通用工具函数
减少重复但不牺牲可读性
```

不适合：

```txt
为了炫技写复杂抽象
把简单业务逻辑泛型化
模仿 TypeScript 类型体操
```

### 2.3 练习任务

```txt
任务 1：写一个 CLI Todo 工具
任务 2：写一个文件统计工具，统计目录下文件数量、大小、后缀分布
任务 3：写一个 JSON 配置读取器
任务 4：写一个并发 URL 检测器
任务 5：写一个简易日志库
```

### 2.4 AI 提示词

```txt
我正在学习 Go 语言基础。请你用“前端开发者能理解的方式”解释下面这个概念：[概念]。
要求：
1. 先说明它解决什么问题
2. 类比到 JavaScript / TypeScript
3. 给一个 Go 最小示例
4. 给一个后端真实场景
5. 列出常见误区
```

---

## 3. 第二关：Go 工具链与工程规范

### 3.1 通关目标

会使用 Go 的核心命令、模块系统、代码格式化、静态检查和安全扫描工具。

你需要掌握：

```txt
go run
go build
go install
go fmt
go test
go mod
go doc
go clean
go version
go generate
go vet
goimports
staticcheck
golangci-lint
govulncheck
pprof
trace
race detector
```

### 3.2 核心命令

```bash
go mod init example.com/myapp
go mod tidy
go run ./cmd/server
go build -o bin/server ./cmd/server
go test ./...
go test -race ./...
go test -cover ./...
go vet ./...
gofmt -w .
go install github.com/some/tool@latest
govulncheck ./...
```

### 3.3 目录结构建议

小型项目：

```txt
myapp/
  cmd/
    server/
      main.go
  internal/
    handler/
    service/
    repository/
    model/
    config/
  migrations/
  docs/
  go.mod
```

中型项目：

```txt
myapp/
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
      grpc/
    middleware/
    platform/
      logger/
      database/
      redis/
  migrations/
  deployments/
  docs/
  scripts/
  Makefile
```

### 3.4 工程规范

你要形成这些习惯：

```txt
每次提交前 gofmt
每次提交前 go test ./...
依赖变更后 go mod tidy
所有入口从 cmd 开始
业务代码尽量放 internal
不要把所有东西塞 main.go
不要 handler 里直接写 SQL
不要 panic 处理普通业务错误
```

### 3.5 AI 提示词

```txt
请你作为 Go 工程规范审查员，审查我的项目结构。
重点看：
1. 目录是否清晰
2. 包职责是否合理
3. 是否存在循环依赖风险
4. 是否符合 Go 项目惯例
5. 哪些地方过度设计
6. 哪些地方未来会难维护
```

---

## 4. 第三关：标准库与后端基础能力

### 4.1 通关目标

Go 的标准库很强，后端初学阶段不要过早依赖框架。你需要先掌握标准库能解决什么问题。

必学标准库：

```txt
fmt
os
io
bufio
net/http
encoding/json
time
context
sync
errors
log/slog
regexp
flag
embed
testing
httptest
```

### 4.2 标准库能力树

```txt
输入输出：io / os / bufio
网络服务：net/http
JSON：encoding/json
时间处理：time
上下文控制：context
并发同步：sync
错误处理：errors
日志：log/slog
测试：testing / httptest
命令行：flag
资源嵌入：embed
```

### 4.3 练习任务

```txt
任务 1：用 net/http 写一个 Hello API
任务 2：用 encoding/json 做请求和响应编解码
任务 3：用 context 给请求加超时
任务 4：用 slog 输出结构化日志
任务 5：用 httptest 给 handler 写测试
任务 6：用 embed 打包静态配置文件
```

### 4.4 验收标准

```txt
能不用框架写一个 HTTP API
能写 JSON 请求和响应
能理解 context 的取消和超时
能写基础测试
能使用 slog 输出结构化日志
```

---

## 5. 第四关：并发与 Context

### 5.1 通关目标

并发是 Go 的核心优势之一，但你不需要一开始写复杂并发模型。你要重点掌握后端场景下怎么安全使用并发。

你需要掌握：

```txt
goroutine
channel
buffered channel
unbuffered channel
select
sync.WaitGroup
sync.Mutex
sync.RWMutex
context.Context
deadline
timeout
cancellation
worker pool
fan-in
fan-out
pipeline
race detector
```

### 5.2 你必须理解的问题

```txt
什么时候需要 goroutine？
什么时候不该开 goroutine？
goroutine 泄露是什么？
channel 什么时候会阻塞？
context 为什么要从请求入口传下去？
数据库查询为什么也要带 context？
超时、取消、重试之间是什么关系？
```

### 5.3 后端真实场景

```txt
并发请求多个第三方服务
异步生成报告
批量处理任务
定时清理过期数据
worker 消费队列
请求超时取消数据库查询
```

### 5.4 练习任务

```txt
任务 1：并发检查 100 个 URL 的状态码
任务 2：实现一个 worker pool
任务 3：实现一个带 context 超时的任务
任务 4：写一个 fan-in 聚合多个 channel 结果
任务 5：故意制造 data race，然后用 go test -race 检测
```

### 5.5 AI 提示词

```txt
请你作为 Go 并发审查员，检查下面这段代码。
重点看：
1. 是否有 goroutine 泄露
2. channel 是否可能死锁
3. 是否需要 context
4. 是否有 data race
5. 是否应该用 worker pool
6. 是否有更简单的同步方式
```

---

## 6. 第五关：HTTP、REST API 与 Web 服务

### 6.1 通关目标

你要从“会写接口”升级为“会设计后端系统边界”。

你需要掌握：

```txt
HTTP 请求与响应
HTTP method
HTTP status code
RESTful API
JSON API
路由
中间件
参数校验
统一错误响应
分页
排序
筛选
CORS
文件上传
接口版本管理
OpenAPI / Swagger
```

### 6.2 框架选择

建议学习顺序：

```txt
net/http → chi 或 Gin → OpenAPI
```

框架选择：

```txt
想贴近标准库：chi
想资料多、生态广：Gin
想了解高性能非 net/http 生态：Fiber，后置学习
```

### 6.3 REST API 设计规范

#### 路径设计

```txt
GET    /api/v1/bookmarks
POST   /api/v1/bookmarks
GET    /api/v1/bookmarks/{id}
PATCH  /api/v1/bookmarks/{id}
DELETE /api/v1/bookmarks/{id}
```

#### 查询参数

```txt
GET /api/v1/bookmarks?page=1&page_size=20&tag=go&keyword=backend
```

#### 成功响应

```json
{
  "code": "OK",
  "message": "success",
  "data": {
    "id": "123",
    "title": "Go Backend Roadmap"
  }
}
```

#### 错误响应

```json
{
  "code": "VALIDATION_ERROR",
  "message": "title is required",
  "details": [
    {
      "field": "title",
      "reason": "required"
    }
  ]
}
```

### 6.4 练习项目：Todo API

功能：

```txt
创建 Todo
查询 Todo 列表
查询 Todo 详情
更新 Todo
删除 Todo
统一错误响应
基础参数校验
```

验收标准：

```txt
能用 curl 调用所有接口
能返回正确状态码
能处理参数错误
能写 handler 测试
能生成一份接口文档
```

### 6.5 AI 提示词

```txt
请你作为后端 API 设计审查员，审查下面这些接口。
重点看：
1. 路径是否符合 REST 风格
2. HTTP 方法是否合理
3. 状态码是否合理
4. 请求和响应字段是否稳定
5. 是否缺少分页、筛选、排序
6. 是否方便前端调用
7. 是否适合未来扩展
```

---

## 7. 第六关：数据库、SQL、事务与索引

### 7.1 通关目标

数据库是后端真正入门的分水岭。你要从“接口能跑”升级为“数据模型正确、查询稳定、事务可靠”。

你需要掌握：

```txt
PostgreSQL / MySQL 基础
表设计
主键
外键
唯一约束
索引
SQL 查询
JOIN
事务
隔离级别
连接池
migration
N+1 查询
慢查询
软删除
乐观锁
悲观锁
```

### 7.2 推荐技术栈

```txt
数据库：PostgreSQL
Go 访问：pgx / database/sql
类型安全 SQL：sqlc
迁移：golang-migrate
本地环境：docker-compose
```

### 7.3 不建议一开始依赖 ORM

你可以了解 GORM，但不要一开始就用 ORM 遮住 SQL。

推荐顺序：

```txt
手写 SQL → pgx/database/sql → sqlc → 了解 GORM/Bun
```

原因：

```txt
后端必须看得懂 SQL
必须知道索引怎么生效
必须知道事务边界
必须能排查慢查询
必须理解 ORM 背后生成了什么
```

### 7.4 表设计训练

以 Bookmark Manager 为例：

```txt
users
- id
- email
- password_hash
- role
- created_at
- updated_at

bookmarks
- id
- user_id
- title
- url
- description
- created_at
- updated_at
- deleted_at

tags
- id
- user_id
- name
- created_at

bookmark_tags
- bookmark_id
- tag_id
```

你需要思考：

```txt
email 是否需要唯一索引？
user_id 是否需要索引？
bookmark_tags 是否需要联合主键？
软删除后唯一约束怎么处理？
分页用 offset 还是 cursor？
搜索用 LIKE 还是全文索引？
```

### 7.5 事务场景

需要事务的情况：

```txt
创建订单 + 扣库存
创建文章 + 写标签关系
转账
创建检测任务 + 写入任务队列记录
更新用户信息 + 写审计日志
```

不一定需要事务的情况：

```txt
单表简单查询
纯缓存读取
无强一致要求的日志写入
异步通知
```

### 7.6 练习项目：Bookmark API v2

功能：

```txt
用户表
书签表
标签表
多对多关系
数据库 migration
分页查询
按标签筛选
按关键词搜索
事务创建书签和标签关系
```

验收标准：

```txt
能初始化数据库
能执行 migration up/down
能用 SQL 查询数据
能解释每个索引的作用
能说明哪些操作需要事务
能通过集成测试
```

### 7.7 AI 提示词

```txt
请你作为数据库架构师，审查下面这个表设计。
重点看：
1. 字段类型是否合理
2. 主键和外键是否合理
3. 是否缺少索引
4. 查询场景下是否会慢
5. 是否存在 N+1 风险
6. 是否需要事务
7. 数据量到 100 万后会出什么问题
```

---

## 8. 第七关：项目结构、分层架构与业务建模

### 8.1 通关目标

你要从“handler 里写所有逻辑”升级为“分层清晰、边界稳定、可测试、可维护”。

经典分层：

```txt
handler / transport：处理 HTTP 请求和响应
service / usecase：处理业务逻辑
repository / dao：处理数据访问
model / entity：表达核心数据结构
middleware：处理横切逻辑
config：配置管理
```

### 8.2 推荐目录

```txt
backend/
  cmd/
    api/
      main.go
    worker/
      main.go
  internal/
    config/
    domain/
      user.go
      bookmark.go
      scan.go
    service/
      user_service.go
      bookmark_service.go
      scan_service.go
    repository/
      user_repo.go
      bookmark_repo.go
      scan_repo.go
    transport/
      http/
        handler/
        middleware/
        router.go
    platform/
      database/
      redis/
      logger/
      storage/
  migrations/
  deployments/
  docs/
  scripts/
  Makefile
```

### 8.3 每层放什么

#### Handler 层

负责：

```txt
解析请求
参数校验
调用 service
组装响应
处理 HTTP 状态码
```

不负责：

```txt
不写复杂业务逻辑
不直接写 SQL
不直接操作事务
不决定核心业务规则
```

#### Service 层

负责：

```txt
业务流程
权限判断
事务编排
调用多个 repository
调用外部服务
处理领域规则
```

#### Repository 层

负责：

```txt
SQL 查询
数据持久化
数据库错误转换
隐藏数据库细节
```

#### Middleware 层

负责：

```txt
认证
日志
请求 ID
CORS
限流
recover
metrics
```

### 8.4 常见模式

```txt
DTO
Entity
Repository
Service
Middleware
Dependency Injection
Config
Error Wrapping
Context Propagation
Transaction Script
State Machine
```

### 8.5 AI 提示词

```txt
我正在设计一个 Go 后端项目。请你判断下面这些逻辑应该放在哪一层：
- handler
- service
- repository
- middleware
- domain/model

请先说明分层原则，再逐条判断，并指出如果放错层会造成什么问题。
```

---

## 9. 第八关：认证、鉴权与安全

### 9.1 通关目标

你要理解：

```txt
认证 Authentication：你是谁？
鉴权 Authorization：你能做什么？
```

你需要掌握：

```txt
密码哈希
bcrypt / argon2
JWT
Session
Access Token
Refresh Token
RBAC
权限校验
CORS
CSRF
XSS 基础
SQL 注入
接口限流
敏感配置管理
审计日志
govulncheck
```

### 9.2 登录系统基础流程

```txt
用户注册
  ↓
密码哈希后入库
  ↓
用户登录
  ↓
校验密码
  ↓
签发 access token / refresh token
  ↓
请求携带 token
  ↓
中间件解析 token
  ↓
注入 user_id / role 到 context
  ↓
业务层做权限判断
```

### 9.3 RBAC 基础模型

```txt
user
role
permission
user_role
role_permission
```

小项目可以简化：

```txt
users.role = user / admin
```

### 9.4 常见安全问题

```txt
密码明文存储
JWT 永不过期
refresh token 不可撤销
用户 A 访问用户 B 的资源
接口没有限流
错误信息泄露内部细节
SQL 拼接导致注入
敏感配置提交到 Git
上传文件不校验类型
```

### 9.5 项目任务：Bookmark API v3

功能：

```txt
注册
登录
密码哈希
JWT 鉴权
Refresh Token
用户只能访问自己的资源
管理员角色
登录失败限流
审计日志
```

验收标准：

```txt
密码不明文入库
未登录不能访问受保护接口
用户不能访问别人的数据
管理员接口只允许 admin
token 过期后需要刷新
登录失败过多会被限制
```

### 9.6 AI 提示词

```txt
请你作为安全审查员，审查我的登录和鉴权方案。
重点看：
1. token 生命周期是否合理
2. refresh token 是否安全
3. 是否存在越权风险
4. 是否存在重放攻击风险
5. 是否需要限流
6. 是否有敏感信息泄露
7. 是否适合小型生产项目
```

---

## 10. 第九关：缓存、Redis 与性能优化

### 10.1 通关目标

你要理解缓存不是万能药。缓存是为了减少重复计算、降低数据库压力、改善响应速度，但会引入一致性问题。

你需要掌握：

```txt
Redis 基础数据结构
String
Hash
List
Set
Sorted Set
TTL
缓存穿透
缓存击穿
缓存雪崩
缓存一致性
分布式锁基础
限流
排行榜
会话存储
```

### 10.2 常见缓存场景

```txt
用户信息缓存
接口热点数据缓存
验证码
登录态
限流计数器
异步任务状态
排行榜
短链接映射
```

### 10.3 不应该缓存的场景

```txt
数据访问频率低
数据强一致要求高
缓存维护成本大于收益
查询本身已经很快
业务还没有性能瓶颈
```

### 10.4 缓存策略

```txt
Cache Aside：业务最常用
Read Through：读穿透
Write Through：写穿透
Write Back：异步写回
```

小项目优先学 Cache Aside：

```txt
读：先读缓存，未命中读数据库，再写缓存
写：先写数据库，再删除缓存
```

### 10.5 AI 提示词

```txt
请你作为后端性能优化顾问，判断下面这个场景是否需要缓存。
请分析：
1. 访问频率
2. 数据变化频率
3. 一致性要求
4. 缓存 key 设计
5. TTL 设计
6. 缓存穿透/击穿/雪崩风险
7. 不使用缓存是否也能接受
```

---

## 11. 第十关：异步任务、消息队列与幂等

### 11.1 通关目标

你要理解：不是所有事情都应该同步完成。

你需要掌握：

```txt
异步任务
消息队列
生产者
消费者
重试
死信队列
延迟任务
幂等
去重
最终一致性
事件驱动
Outbox Pattern
```

### 11.2 适合异步的场景

```txt
发送邮件
发送短信
生成报告
图片分析
视频转码
同步第三方系统
日志聚合
批量导入导出
AI 检测任务
```

### 11.3 不适合同步做的任务

```txt
耗时不可控
依赖第三方服务
失败需要重试
用户不需要立即拿到结果
任务可以排队处理
```

### 11.4 推荐学习顺序

```txt
Go goroutine 后台任务
  ↓
Redis + worker
  ↓
Asynq
  ↓
RabbitMQ
  ↓
Kafka
```

### 11.5 幂等设计

你必须理解：消息可能被重复消费。

幂等方案：

```txt
唯一请求 ID
任务状态机
数据库唯一约束
消费记录表
分布式锁
业务结果覆盖写
```

### 11.6 项目任务：BemoLens Scan Worker

流程：

```txt
用户上传图片
  ↓
创建 scan 记录，状态 pending
  ↓
写入任务队列
  ↓
worker 消费任务
  ↓
状态改为 processing
  ↓
执行图片检测
  ↓
成功：completed
  ↓
失败：failed / retry
```

任务状态：

```txt
pending
processing
completed
failed
expired
```

验收标准：

```txt
任务可以异步执行
用户可以查询任务状态
worker 挂了后任务不会永久丢失
重复消费不会产生错误结果
失败任务可以重试
超过次数进入失败状态
```

### 11.7 AI 提示词

```txt
请你作为后端架构师，审查这个异步任务设计。
重点看：
1. 是否真的需要消息队列
2. 消息格式是否合理
3. 消费失败如何处理
4. 是否支持重试
5. 是否幂等
6. 是否会重复通知
7. 是否需要死信队列
8. 状态机是否完整
```

---

## 12. 第十一关：测试、文档与代码质量

### 12.1 通关目标

你要能证明自己的代码是可靠的，而不是“我本地跑过了”。

你需要掌握：

```txt
testing
表驱动测试
mock
stub
httptest
集成测试
基准测试
覆盖率
race test
OpenAPI
README
架构文档
```

### 12.2 测试类型

```txt
单元测试：测试单个函数或模块
集成测试：测试数据库、Redis、HTTP 等组合行为
端到端测试：模拟真实用户流程
基准测试：测试性能
回归测试：防止旧 bug 复现
```

### 12.3 Go 测试命令

```bash
go test ./...
go test -v ./...
go test -race ./...
go test -cover ./...
go test -bench=. ./...
```

### 12.4 表驱动测试模板

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name string
        input string
        want string
        wantErr bool
    }{
        {
            name: "valid input",
            input: "hello",
            want: "HELLO",
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Something(tt.input)
            if (err != nil) != tt.wantErr {
                t.Fatalf("unexpected error: %v", err)
            }
            if got != tt.want {
                t.Fatalf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 12.5 文档清单

一个可维护后端项目至少需要：

```txt
README.md
.env.example
API 文档
数据库设计文档
部署文档
错误码文档
架构说明
本地开发指南
```

### 12.6 AI 提示词

```txt
请你作为测试工程师，帮我为下面这个 Go service 设计测试方案。
要求：
1. 列出需要测试的正常路径
2. 列出异常路径
3. 列出边界条件
4. 判断哪些需要 mock
5. 判断哪些需要集成测试
6. 给出表驱动测试用例设计
```

---

## 13. 第十二关：部署、Docker、Linux 与 CI/CD

### 13.1 通关目标

你要能把服务从本地搬到服务器，并稳定运行。

你需要掌握：

```txt
Linux 基础命令
环境变量
Dockerfile
docker-compose
镜像构建
容器网络
数据卷
Nginx / Caddy
HTTPS
GitHub Actions
VPS 部署
健康检查
滚动重启基础
```

### 13.2 本地开发环境

```txt
app：Go API
postgres：数据库
redis：缓存/队列
minio：本地对象存储，可选
prometheus：指标采集，可选
grafana：监控面板，可选
```

### 13.3 Makefile 示例目标

```txt
make dev
make test
make lint
make migrate-up
make migrate-down
make docker-build
make docker-up
make docker-down
make deploy
```

### 13.4 生产部署最小清单

```txt
Dockerfile 使用多阶段构建
配置通过环境变量注入
数据库使用持久化卷或托管数据库
日志输出到 stdout
提供 /healthz 和 /readyz
反向代理开启 HTTPS
CI 中执行测试和 lint
敏感配置不提交到 Git
```

### 13.5 AI 提示词

```txt
请你作为 DevOps 审查员，审查我的 Go 后端部署方案。
重点看：
1. Dockerfile 是否适合生产
2. 镜像是否过大
3. 环境变量是否合理
4. 数据库是否持久化
5. 日志是否可收集
6. 是否有健康检查
7. CI/CD 是否有明显问题
8. 服务器重启后服务能否自动恢复
```

---

## 14. 第十三关：可观测性、稳定性与线上排查

### 14.1 通关目标

服务上线不是终点，能排查线上问题才是真正开始。

你需要掌握：

```txt
结构化日志
request_id
trace_id
metrics
tracing
OpenTelemetry
Prometheus
Grafana
Sentry
健康检查
慢查询日志
错误码
告警
限流
熔断
重试
降级
```

### 14.2 三件套

```txt
Logs：发生了什么
Metrics：系统整体状态如何
Traces：一次请求经过了哪些环节
```

### 14.3 必备指标

```txt
请求总量
请求耗时
错误率
接口状态码分布
数据库连接池使用率
Redis 请求耗时
队列长度
任务成功率
任务失败率
CPU / 内存 / 磁盘 / 网络
```

### 14.4 必备日志字段

```txt
timestamp
level
request_id
trace_id
method
path
status
latency
user_id
error_code
error_message
```

### 14.5 线上排查问题

你要能回答：

```txt
哪个接口慢？
慢在数据库还是外部服务？
错误率什么时候升高？
是哪一类用户受影响？
是代码发布引起的吗？
数据库连接池是否打满？
队列是否积压？
任务失败能否重试？
```

### 14.6 AI 提示词

```txt
请你作为 SRE，帮我设计这个 Go 后端服务的可观测性方案。
要求包括：
1. 应该记录哪些日志
2. 应该暴露哪些 metrics
3. 哪些地方需要 tracing
4. 哪些错误需要报警
5. 如何定位接口慢的问题
6. 如何区分应用问题、数据库问题和第三方服务问题
```

---

## 15. 第十四关：高级 Go 与性能调优

### 15.1 通关目标

这一关不是一开始就学，而是在你写过真实项目后补。

你需要掌握：

```txt
pprof
trace
escape analysis
内存分配
GC 基础
反射
unsafe
cgo
build tags
compiler flags
linker flags
cross compilation
插件机制
```

### 15.2 学习顺序

```txt
pprof 看 CPU / 内存
  ↓
race detector 查并发问题
  ↓
escape analysis 理解堆分配
  ↓
GC 基础
  ↓
反射
  ↓
unsafe / cgo / compiler flags
```

### 15.3 不要过早优化

先问：

```txt
是否真的慢？
慢在哪里？
有指标证明吗？
数据库是不是瓶颈？
网络是不是瓶颈？
算法是不是瓶颈？
是否可以通过缓存或异步解决？
```

### 15.4 AI 提示词

```txt
请你作为 Go 性能优化专家，分析下面这个性能问题。
要求：
1. 先判断是否需要优化
2. 可能瓶颈在哪里
3. 需要采集哪些数据
4. 如何用 pprof 验证
5. 优先级最高的优化是什么
6. 哪些优化可能是过早优化
```

---

## 16. 第十五关：微服务、分布式与系统设计

### 16.1 通关目标

这不是初学阶段的重点，但你最终要理解。

你需要掌握：

```txt
服务拆分
单体架构
模块化单体
微服务
RPC
gRPC
Protocol Buffers
服务发现
配置中心
API Gateway
分布式锁
分布式事务
最终一致性
CAP
BASE
限流
熔断
降级
幂等
重试
消息队列
事件驱动
Outbox Pattern
Saga
```

### 16.2 学习顺序

```txt
先写好单体
  ↓
模块化单体
  ↓
拆出 worker
  ↓
拆出独立服务
  ↓
引入 gRPC
  ↓
引入服务发现和网关
  ↓
再考虑 Kubernetes 和 Service Mesh
```

### 16.3 判断是否需要微服务

需要微服务的信号：

```txt
团队规模较大
模块发布节奏不同
不同模块资源需求不同
单体部署风险过高
边界非常清晰
有成熟的 DevOps 能力
```

不需要微服务的信号：

```txt
一个人或小团队
业务还没验证
接口数量很少
部署能力弱
监控能力弱
连单体都维护不好
只是为了“高级”
```

### 16.4 AI 提示词

```txt
请你作为系统架构师，判断下面这个系统是否需要拆成微服务。
请分析：
1. 当前业务规模
2. 团队规模
3. 模块边界
4. 部署复杂度
5. 数据一致性问题
6. 运维成本
7. 如果不拆，模块化单体怎么设计
8. 如果拆，第一步应该拆哪个服务
```

---

## 17. 四个主线项目

你不要只看教程。你需要做 4 个项目，把能力逐层补齐。

---

### 项目 1：Go Todo API

目标：入门 Go + HTTP。

功能：

```txt
Todo 增删改查
JSON API
统一错误响应
内存存储
基础测试
```

技术点：

```txt
net/http
encoding/json
handler
testing
httptest
```

验收：

```txt
能用 curl 完整调用
有 README
有测试
有统一错误格式
```

---

### 项目 2：Bookmark Manager

目标：进入真实后端。

功能：

```txt
用户注册登录
收藏链接
标签
搜索
分页
PostgreSQL
JWT
Docker
```

技术点：

```txt
Gin/chi
PostgreSQL
pgx/sqlc
migration
JWT
bcrypt
Docker Compose
```

验收：

```txt
有数据库设计文档
有 migration
有接口文档
有登录鉴权
有 Docker 一键启动
```

---

### 项目 3：Mini Blog / CMS API

目标：服务你的博客和内容管理需求。

功能：

```txt
文章管理
草稿
发布
标签
分类
Markdown 内容
图片上传
管理员权限
```

技术点：

```txt
RBAC
文件上传
对象存储
后台 API
内容模型
审计日志
```

验收：

```txt
能支撑一个 Astro/Next.js 博客后台
能上传图片
能管理草稿和发布状态
能区分管理员权限
```

---

### 项目 4：BemoLens Backend

目标：结合你的 AI 图片鉴别项目，做一个真正有业务含义的后端。

功能：

```txt
图片上传
创建检测任务
异步分析
状态轮询
结果报告
用户历史
限流
日志
监控
Docker 部署
```

技术栈：

```txt
Go
chi/Gin
PostgreSQL
Redis
MinIO / S3 / R2
Docker Compose
Prometheus
Grafana
OpenTelemetry
```

接口示例：

```txt
POST   /api/v1/scans
GET    /api/v1/scans/{id}
GET    /api/v1/scans/{id}/status
GET    /api/v1/scans/{id}/report
DELETE /api/v1/scans/{id}
```

状态机：

```txt
pending → processing → completed
pending → processing → failed
completed → expired
failed → retrying → processing
```

验收：

```txt
上传后不阻塞等待检测完成
任务可查询状态
worker 可独立运行
失败可重试
重复消费安全
有日志和 metrics
能部署到服务器
```

---

## 18. 12 周通关计划

### 第 1 周：Go 基础语法

任务：

```txt
变量、类型、函数、结构体、接口
完成 CLI Todo
完成文件统计工具
```

验收：

```txt
能独立写 Go 小程序
能解释 slice、map、struct、interface
```

---

### 第 2 周：错误处理、测试、标准库

任务：

```txt
错误处理
context 初步
testing
httptest
slog
encoding/json
```

验收：

```txt
能写表驱动测试
能写 JSON API demo
能处理错误并包装
```

---

### 第 3 周：HTTP API

任务：

```txt
net/http
路由
中间件
统一响应
Todo API
```

验收：

```txt
Todo API 可用
有 handler 测试
有 curl 文档
```

---

### 第 4 周：Web 框架与接口设计

任务：

```txt
chi 或 Gin
REST API 规范
OpenAPI
Bookmark API v1
```

验收：

```txt
完成 Bookmark CRUD
接口文档清晰
错误响应统一
```

---

### 第 5 周：PostgreSQL 与 SQL

任务：

```txt
表设计
SQL
索引
JOIN
migration
pgx/sqlc
```

验收：

```txt
Bookmark API 接入数据库
能解释索引设计
能跑 migration
```

---

### 第 6 周：事务与项目分层

任务：

```txt
handler/service/repository
事务
N+1
集成测试
```

验收：

```txt
项目结构清晰
复杂创建操作使用事务
有集成测试
```

---

### 第 7 周：认证鉴权

任务：

```txt
注册登录
密码哈希
JWT
Refresh Token
RBAC
```

验收：

```txt
用户只能访问自己的数据
管理员接口受保护
token 有过期和刷新机制
```

---

### 第 8 周：Redis 与缓存

任务：

```txt
Redis 基础
缓存策略
限流
任务状态缓存
```

验收：

```txt
能判断哪些数据适合缓存
实现一个简单限流中间件
```

---

### 第 9 周：异步任务

任务：

```txt
worker
Redis queue / Asynq
任务状态机
重试
幂等
```

验收：

```txt
实现异步任务系统
失败可重试
重复消费安全
```

---

### 第 10 周：Docker 与部署

任务：

```txt
Dockerfile
docker-compose
Makefile
GitHub Actions
VPS 部署
```

验收：

```txt
一键本地启动
CI 自动测试
服务部署到服务器
```

---

### 第 11 周：可观测性

任务：

```txt
结构化日志
request_id
metrics
Prometheus
Grafana
healthz/readyz
```

验收：

```txt
能看到请求量、耗时、错误率
日志能定位单次请求
```

---

### 第 12 周：BemoLens Backend MVP

任务：

```txt
图片上传
检测任务
worker
状态查询
结果报告
基础监控
```

验收：

```txt
形成一个可展示的完整后端项目
能写进简历
能讲清楚架构设计
```

---

## 19. AI 辅助学习工作流

### 19.1 正确用法

```txt
AI 负责：
- 解释概念
- 生成练习
- code review
- 找 bug
- 生成测试
- 比较方案
- 模拟面试
- 生成文档
- 解释报错

你负责：
- 判断方案
- 亲自实现核心逻辑
- 理解边界
- 做取舍
- debug
- 复盘
```

### 19.2 学习闭环

```txt
1. 让 AI 讲概念
2. 让 AI 给最小 demo
3. 你关掉 AI 自己写一遍
4. 跑测试
5. 把报错给 AI
6. 修完后让 AI code review
7. 让 AI 问你 10 个追问
8. 你写复盘
```

### 19.3 万能 Code Review 提示词

```txt
请你作为资深 Go 后端工程师 review 下面代码。
重点看：
1. Go 惯用写法
2. 错误处理
3. context 使用
4. 并发安全
5. 数据库访问
6. 可测试性
7. 是否适合生产
8. 是否存在过度设计
请先指出问题，再给修改建议，不要直接大段重写。
```

### 19.4 万能项目拆解提示词

```txt
我准备做一个 Go 后端项目：[项目描述]。
请你帮我拆成 5 个阶段：
1. MVP
2. 数据库版
3. 认证鉴权版
4. Docker 部署版
5. 可观测性/生产化版

每个阶段请给：
- 目标
- 功能列表
- 目录结构
- 需要掌握的知识点
- 验收标准
```

### 19.5 万能系统设计提示词

```txt
请你作为后端架构师，帮我设计 [系统]。
请按以下结构：
1. 需求澄清
2. 核心实体
3. API 设计
4. 数据库设计
5. 缓存设计
6. 异步任务设计
7. 权限设计
8. 错误处理
9. 可观测性
10. 未来扩展点
```

---

## 20. 后端能力树

```txt
后端能力
├── 语言能力
│   ├── Go 语法
│   ├── 标准库
│   ├── 并发
│   ├── 错误处理
│   └── 工具链
├── API 能力
│   ├── HTTP
│   ├── REST
│   ├── Web 框架
│   ├── 参数校验
│   └── OpenAPI
├── 数据能力
│   ├── SQL
│   ├── 表设计
│   ├── 索引
│   ├── 事务
│   └── 缓存
├── 工程能力
│   ├── 项目结构
│   ├── 分层架构
│   ├── 测试
│   ├── 文档
│   └── CI/CD
├── 安全能力
│   ├── 认证
│   ├── 鉴权
│   ├── 密码哈希
│   ├── 限流
│   └── 漏洞扫描
├── 运维能力
│   ├── Docker
│   ├── Linux
│   ├── 部署
│   ├── 日志
│   └── 监控
├── 稳定性能力
│   ├── 超时
│   ├── 重试
│   ├── 熔断
│   ├── 降级
│   └── 幂等
└── 架构能力
    ├── 异步任务
    ├── 消息队列
    ├── 微服务
    ├── 分布式一致性
    └── 系统设计
```

---

## 21. 秋招 / 面试衔接

### 21.1 简历项目表达模板

```txt
项目名称：BemoLens Backend - AI 图片检测任务系统

项目描述：
基于 Go 构建的图片检测后端服务，支持图片上传、异步检测任务、状态轮询、检测结果报告、用户历史记录和基础可观测性。

技术栈：
Go、Gin/chi、PostgreSQL、Redis、Docker、Prometheus、Grafana、OpenTelemetry

核心工作：
1. 设计 scan 任务状态机，支持 pending、processing、completed、failed、expired 等状态
2. 使用 Redis/Asynq 实现异步任务处理和失败重试
3. 使用 PostgreSQL 存储用户、任务和检测结果
4. 设计统一错误响应和鉴权中间件
5. 接入结构化日志、metrics 和健康检查，方便线上排查

项目亮点：
1. 将耗时图片检测从同步接口改为异步任务，提高接口响应稳定性
2. 通过任务状态机和幂等设计避免重复消费导致的数据错误
3. 使用 Docker Compose 编排本地开发环境，降低部署和调试成本
```

### 21.2 面试常问问题

Go：

```txt
slice 和 array 的区别
map 是否并发安全
interface 的底层理解
值接收者和指针接收者区别
goroutine 泄露怎么排查
channel 什么时候会阻塞
context 有什么作用
Go 的错误处理方式
defer 的执行顺序
GC 基础
```

数据库：

```txt
索引为什么能加速查询
什么情况下索引失效
事务 ACID
隔离级别
幻读、脏读、不可重复读
如何设计分页
如何排查慢查询
```

后端工程：

```txt
认证和鉴权区别
JWT 的优缺点
如何防止越权
什么是幂等
为什么要异步任务
如何设计重试
缓存一致性怎么处理
如何做限流
如何定位线上慢接口
```

系统设计：

```txt
设计短链接系统
设计登录系统
设计图片上传系统
设计消息通知系统
设计任务队列系统
设计博客后台系统
设计 AI 图片检测系统
```

---

## 22. 最终通关标准

当你完成这份路线后，应该能做到：

```txt
1. 能用 Go 独立写一个后端服务
2. 能设计 REST API
3. 能设计 PostgreSQL 表结构
4. 能处理登录、鉴权和权限隔离
5. 能使用 Redis 做缓存、限流或任务队列
6. 能写单元测试和接口测试
7. 能用 Docker 部署服务
8. 能用 GitHub Actions 做 CI
9. 能输出结构化日志
10. 能暴露 metrics 并做基础监控
11. 能设计异步任务和状态机
12. 能解释幂等、重试、超时、事务、索引、缓存一致性
13. 能把项目写进简历并讲清楚架构取舍
```

真正的后端通关不是背完所有技术名词，而是你面对一个需求时能说清楚：

```txt
这个需求的核心实体是什么？
接口怎么设计？
表怎么设计？
权限怎么控制？
是否需要事务？
是否需要缓存？
是否需要异步？
失败怎么处理？
上线后怎么观察？
出问题怎么排查？
未来怎么扩展？
```

---

## 23. 参考来源

- roadmap.sh Go Developer Roadmap
- roadmap.sh Go Roadmap PDF
- Go 官方文档
- Go 标准库与工具链实践
- 结合当前学习背景：前端/Next.js API 基础、目标转向 Go 后端生产能力、希望 AI 辅助学习

---

## 24. 最后一条心法

不要把路线图当成“收藏品”。

路线图真正的用法是：

```txt
每学一个点，就写一个小 demo；
每写一个 demo，就放进一个真实项目；
每做一个项目，就补一篇复盘；
每篇复盘，都变成你的面试表达和工程判断力。
```

AI 会让写代码更便宜，但不会让系统理解更便宜。

后端的核心竞争力不是“我会 Go”，而是：

> 我知道一个系统应该怎么被设计、实现、部署、观察、修复和演进。
