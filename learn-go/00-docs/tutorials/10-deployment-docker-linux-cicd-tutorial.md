# 10 部署 Docker Linux 与 CI CD 教程

## 1. 这章学什么

前面你一直主要在本地写代码，这章开始要解决一个更现实的问题：

代码怎么稳定地在别的环境跑起来。

这里你要掌握四件事：

- Linux 基本运行环境
- Docker 打包
- 本地多服务编排
- 最小 CI

---

## 2. Linux 基础别跳过

即使你主要用 Docker，也还是要能看懂 Linux 上的基本问题。

最常用的几类命令：

- 看文件：`ls`、`cat`
- 看进程：`ps`
- 看端口：`lsof -i :8080`
- 看环境变量：`env`
- 看日志：`tail`

例如排查端口占用：

```bash
lsof -i :8080
```

你不一定天天手写系统脚本，但至少要能定位：

- 服务有没有启动
- 端口是不是被占了
- 环境变量有没有注入

---

## 3. 为什么 Go 很适合 Docker

因为 Go 通常可以直接编译成单个二进制。

这意味着：

- 部署简单
- 运行依赖少
- 多阶段构建很顺手

最小 Dockerfile：

```dockerfile
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/server /app/server

EXPOSE 8080

CMD ["/app/server"]
```

这就是最常见的两阶段：

- 第一阶段构建
- 第二阶段运行

---

## 4. Dockerfile 每一层在干什么

### `FROM golang:1.24 AS builder`

用带 Go 编译器的镜像做构建环境。

### `WORKDIR /app`

设置工作目录。

### `COPY go.mod go.sum ./`

先拷依赖描述文件，有利于利用缓存层。

### `RUN go mod download`

先下载依赖。

### `COPY . .`

再拷完整项目代码。

### `RUN ... go build`

编译出最终二进制。

### 第二个 `FROM`

换成更轻的运行时镜像，只带可执行文件。

这能显著减小镜像体积，也减少不必要内容进入生产环境。

---

## 5. 本地构建和运行镜像

构建：

```bash
docker build -t bookmark-service .
```

运行：

```bash
docker run --rm -p 8080:8080 bookmark-service
```

你现在先看懂这件事：

不是“在我的电脑能跑”，而是“这个镜像在标准环境里也能跑”。

---

## 6. Docker Compose 的价值

很多后端项目不是单机程序，而是至少还依赖：

- PostgreSQL
- Redis

这时你就不想每次手动一个个起服务。

最小 `docker-compose.yml`：

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      APP_PORT: 8080
      DATABASE_URL: postgres://postgres:postgres@db:5432/app?sslmode=disable
      REDIS_ADDR: redis:6379
    depends_on:
      - db
      - redis

  db:
    image: postgres:16
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: app
    ports:
      - "5432:5432"

  redis:
    image: redis:7
    ports:
      - "6379:6379"
```

启动：

```bash
docker compose up --build
```

它解决的是环境复现问题，而不只是“偷懒少打命令”。

---

## 7. 健康检查接口最小例子

部署后你要让系统知道服务是不是还活着。

最简单的接口：

```go
package main

import (
	"encoding/json"
	"net/http"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func main() {
	http.HandleFunc("/healthz", healthz)
	http.ListenAndServe(":8080", nil)
}
```

一般先区分两个概念：

- `/healthz`：服务活着没
- `/readyz`：服务是否已经准备好接流量

例如数据库都没连上时，服务也许活着，但不代表 ready。

---

## 8. 日志为什么要打到 stdout

容器环境里，最简单稳定的做法通常是：

- 程序日志直接输出到标准输出

这样：

- Docker 能收
- Kubernetes 能收
- 日志平台也更容易接

所以不要一开始就把日志写死到本地文件路径里。

---

## 9. 最小 CI 例子

GitHub Actions：

`.github/workflows/ci.yml`

```yaml
name: ci

on:
  push:
  pull_request:

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.24"
      - run: go test ./...
```

这是最小持续集成。

它先不解决自动发布，只解决一个更基础的问题：

每次提交至少要自动验证没有明显破坏。

---

## 10. CI 里常见还会加什么

在 `go test` 之外，常见还会加：

- `gofmt` 检查
- `go vet ./...`
- 镜像构建

例如：

```yaml
- run: test -z "$(gofmt -l .)"
- run: go vet ./...
- run: go build ./...
```

你现在不一定一次全上，但要知道 CI 是逐步加门槛的。

---

## 11. 部署思路先求稳定，不求炫

在学习阶段，更合理的顺序通常是：

1. 本地 `go run` 跑通
2. 本地 `docker build` 跑通
3. 本地 `docker compose up` 跑通
4. CI 能自动测

不要一开始就把注意力全放在复杂云原生名词上。

真正高价值的是：

- 环境可复现
- 构建可重复
- 失败可排查

---

## 12. 这一章建议你手敲的练习

### 练习 1：Dockerfile

为自己的 Go 服务写一个多阶段 Dockerfile。

### 练习 2：Compose 编排

把 API、PostgreSQL、Redis 一起编排启动。

### 练习 3：健康检查

增加 `/healthz` 和 `/readyz` 两个接口。

### 练习 4：最小 CI

增加一条在提交时自动跑 `go test ./...` 的工作流。

---

## 13. 这章学完的最低标准

你至少应该能做到：

- 看懂并编写基础 Go Dockerfile
- 用 `docker compose` 拉起本地依赖环境
- 为服务提供健康检查接口
- 理解日志输出到 stdout 的意义
- 配一个最小 CI 自动跑测试

做到这里，你的项目才开始具备“可交付”的样子。
