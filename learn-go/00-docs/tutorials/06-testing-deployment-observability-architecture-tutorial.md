# 06 测试 部署 可观测性与架构教程

## 1. 这一章讲什么

前面几章帮你把“业务功能”搭起来了，这一章要补的是“怎么把它做成一个像样的项目”。

重点有四个：

- 测试
- 部署
- 可观测性
- 架构判断

---

## 2. 单元测试最小例子

先写一个 service 风格的小函数：

`status.go`

```go
package main

import "errors"

func normalizeStatus(s string) (string, error) {
	switch s {
	case "pending", "processing", "completed":
		return s, nil
	default:
		return "", errors.New("invalid status")
	}
}
```

`status_test.go`

```go
package main

import "testing"

func TestNormalizeStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "valid pending", input: "pending", want: "pending", wantErr: false},
		{name: "valid completed", input: "completed", want: "completed", wantErr: false},
		{name: "invalid", input: "done", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeStatus(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error state: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
```

运行：

```bash
go test ./...
```

---

## 3. `httptest` 最小例子

如果你已经有一个 handler，可以这样测它。

```go
package main

import (
	"encoding/json"
	"net/http"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "pong",
	})
}
```

测试：

```go
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	pingHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	if body == "" {
		t.Fatal("empty body")
	}
}
```

这就是 handler 层最常见的测试入口。

---

## 4. 最小 Dockerfile

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

你先记住它的两层：

- 第一层负责构建
- 第二层负责运行

这就是多阶段构建。

---

## 5. 最小 `docker-compose.yml`

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
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

作用很直接：

- 一次拉起 app
- 一次拉起 PostgreSQL
- 一次拉起 Redis

---

## 6. GitHub Actions 最小 CI 例子

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

这是最小 CI，不复杂，但很有价值。

---

## 7. 结构化日志最小例子

```go
package main

import (
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("request finished",
		"method", "GET",
		"path", "/api/v1/bookmarks",
		"status", 200,
		"latency_ms", 12,
	)
}
```

你要关注的是字段化，而不是“日志看起来漂亮”。

---

## 8. Metrics 最小感觉

如果你后面接入 Prometheus，常见思路是：

- 请求总数
- 请求耗时
- 错误总数

最简单的观念例子：

```txt
http_requests_total
http_request_duration_seconds
http_request_errors_total
```

你现在先不用把整套 Prometheus 接完，但要知道 metrics 是在看“整体趋势”，不是某一次请求。

---

## 9. `request_id` 最小概念

为什么要有 `request_id`：

- 一次请求经过多个模块时，你需要把日志串起来
- 线上排查时，不能只靠时间猜

最简单做法就是：

- 入口生成一个 ID
- 放进 context
- 日志统一带上它

---

## 10. 架构判断最小例子

假设你现在有这些功能：

- 用户
- 书签
- 上传
- 异步检测

你先别急着拆微服务。更合理的第一步通常是模块化单体：

```txt
cmd/
  api/
  worker/
internal/
  user/
  bookmark/
  upload/
  scan/
  platform/
```

这已经能表达边界，而且部署还简单。

---

## 11. 什么情况下别急着拆微服务

如果你现在是这些情况：

- 一个人开发
- 功能还在快速变化
- 监控还没搭好
- 单体都还没整理清楚

那就不要急着拆服务。

先把单体做到：

- 模块清晰
- 日志清晰
- 测试可跑
- 部署可复现

这比“形式上拆了几个服务”更有价值。

---

## 12. 这一章建议你手敲的练习

### 练习 1：单元测试

自己写一个函数，再补一个表驱动测试。

### 练习 2：handler 测试

给 `/ping` 接口写一个 `httptest`。

### 练习 3：Dockerfile

为自己的 Go 服务写一个最小多阶段 Dockerfile。

### 练习 4：日志

用 `slog` 打一条带字段的日志。

---

## 13. 这章学完的最低标准

你至少应该能做到：

- 给函数和 handler 写基础测试
- 看懂 Dockerfile 和 docker-compose
- 理解结构化日志和 metrics 的区别
- 判断当前项目更适合模块化单体还是微服务

做到这里，你的项目才开始像一个真正可维护的后端项目。
