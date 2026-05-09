# 11 日志 指标 链路追踪与可观测性教程

## 1. 这章为什么关键

系统上线后，最痛苦的问题通常不是“不会写功能”，而是：

- 用户说慢，但你不知道慢在哪
- 接口报错，但你不知道影响范围
- 发布后异常，你不知道是代码、数据库还是外部依赖的问题

所以可观测性的核心不是“堆工具”，而是让系统能回答问题。

---

## 2. Logs、Metrics、Traces 各管什么

你可以先这样记：

- Logs：发生了什么
- Metrics：整体趋势怎么样
- Traces：一次请求经过了什么路径

这三者不是互相替代，而是互相补位。

例如一个接口变慢时：

- Metrics 先告诉你整体延迟上升
- Traces 帮你找慢在数据库还是外部服务
- Logs 帮你看某次失败请求的具体上下文

---

## 3. 结构化日志最小例子

Go 现在很适合直接用 `slog`。

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
		"path", "/api/v1/users",
		"status", 200,
		"latency_ms", 12,
	)
}
```

这里真正重要的不是 JSON 形式本身，而是字段化。

你后面做检索和统计时，最依赖的是这些稳定字段。

---

## 4. 日志字段应该怎么设计

一个 HTTP 请求日志常见字段可以有：

- `timestamp`
- `level`
- `request_id`
- `method`
- `path`
- `status`
- `latency_ms`
- `user_id`
- `error_code`

不是每条日志都要全带，但关键链路日志最好有统一结构。

一个坏日志常见长这样：

```txt
request failed maybe db timeout user=12 path=/users
```

人眼也许勉强看得懂，但机器很难稳定检索和聚合。

---

## 5. 请求日志中间件最小例子

```go
package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

func loggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		logger.Info("request finished",
			"method", r.Method,
			"path", r.URL.Path,
			"latency_ms", time.Since(start).Milliseconds(),
		)
	})
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)

	http.ListenAndServe(":8080", loggingMiddleware(logger, mux))
}
```

这个版本还没记录状态码，但已经足够帮你建立中间件式接入的感觉。

---

## 6. Metrics 先看什么

你不用一开始就埋几十种指标，先关注最关键的几类：

- 请求总数
- 请求耗时
- 错误总数
- 队列积压
- 数据库连接池状态

最常见的指标名会类似：

```txt
http_requests_total
http_request_duration_seconds
http_request_errors_total
```

它们回答的是整体趋势，而不是某一次请求细节。

---

## 7. 为什么需要 request_id

当一次请求经过：

- HTTP handler
- service
- repository
- Redis
- PostgreSQL
- 第三方 API

如果没有统一请求标识，排查时你就很难把日志串起来。

最简单的做法是：

1. 入口生成 `request_id`
2. 放进 `context`
3. 所有日志都带上它

---

## 8. `request_id` 最小中间件例子

```go
package main

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

type contextKey string

const requestIDKey contextKey = "request_id"

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := strconv.FormatInt(time.Now().UnixNano(), 10)
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func requestIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(requestIDKey).(string)
	return v
}
```

这不是最终生产级方案，但足够让你先理解贯穿链路的基本思路。

---

## 9. Traces 解决什么问题

当系统越来越复杂时，仅靠日志会有两个问题：

- 单次请求上下游太多，很难人工拼接
- 慢调用到底卡在哪一段不直观

Trace 的价值在于：

- 把一次请求拆成多个 span
- 看到每段耗时
- 看到上下游调用关系

如果你以后接 OpenTelemetry，本质上就是把这件事标准化。

你现在至少要先把：

- `request_id`
- 下游调用耗时
- 错误字段

这些基础打好。

---

## 10. 线上排查的最小顺序

出问题时，不要直接跳进代码海里乱猜。

先按这个顺序想：

1. 影响范围多大
2. 是单接口还是全站问题
3. 是慢、错，还是资源耗尽
4. 变化发生在应用、数据库、缓存还是第三方
5. 最近有没有发布

例如某接口突然超时：

- 先看该接口耗时指标
- 再看数据库慢查询或连接池
- 再看 Redis 和第三方依赖
- 最后结合错误日志和发布时间点

这比“先 grep 所有日志”更有方向。

---

## 11. 报警不是越多越好

最开始可以先盯这几类：

- 错误率突增
- P95/P99 延迟升高
- 队列积压增长
- 数据库连接池耗尽
- 外部依赖超时增多

报警的目标不是让你手机一直响，而是异常发生时真的能触发行动。

---

## 12. 这一章建议你手敲的练习

### 练习 1：结构化日志

用 `slog` 打一条包含方法、路径、状态码的请求日志。

### 练习 2：请求 ID

写一个中间件，在请求入口生成 `request_id`。

### 练习 3：指标设计

为你的项目列一组最先该上的 5 个指标。

### 练习 4：排障清单

假设“用户反馈接口变慢”，写出你的排查顺序。

---

## 13. 这章学完的最低标准

你至少应该能做到：

- 说清 Logs、Metrics、Traces 的区别
- 为 HTTP 请求设计基础结构化日志字段
- 理解 `request_id` 的作用并在代码里贯穿
- 为服务列出最关键的一组指标
- 遇到线上性能问题时有基本排查顺序

做到这里，系统才不只是“能跑”，而是“出了问题也能查”。
