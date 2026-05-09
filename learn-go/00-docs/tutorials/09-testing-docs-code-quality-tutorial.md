# 09 测试 文档与代码质量教程

## 1. 这章为什么重要

项目写到这个阶段，最大的风险已经不再是“不会写功能”，而是：

- 改一点就怕坏别的地方
- 接口行为越来越不稳定
- 过几天连自己都忘了怎么启动

所以这章的重点不是形式化，而是把项目变成“可验证、可维护、可交接”。

---

## 2. 单元测试到底测什么

单元测试最适合测这些：

- 纯业务规则
- 参数转换
- 状态流转
- 错误分支

例如状态规范化：

```go
package main

import "errors"

func normalizeRole(role string) (string, error) {
	switch role {
	case "admin", "user":
		return role, nil
	default:
		return "", errors.New("invalid role")
	}
}
```

对应测试：

```go
package main

import "testing"

func TestNormalizeRole(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "admin", input: "admin", want: "admin"},
		{name: "user", input: "user", want: "user"},
		{name: "invalid", input: "guest", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeRole(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}
```

这就是 Go 很典型的表驱动测试。

---

## 3. 表驱动测试为什么常用

因为很多业务规则都有“多输入、多边界”的特点。

例如：

- 角色校验
- 状态转换
- 分页参数处理
- 权限判断

表驱动的好处是：

- 用例集中
- 更容易补边界
- 结构一致

你会在 Go 项目里反复见到这种写法。

---

## 4. Handler 怎么测

HTTP 层很适合用 `httptest`。

先看 handler：

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

	if rec.Body.String() == "" {
		t.Fatal("empty body")
	}
}
```

这种测试非常适合稳定验证：

- 状态码
- 响应体
- 参数错误路径

---

## 5. 参数错误和鉴权失败也要测

很多人只测成功路径，这是不够的。

例如：

```go
func userHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
```

测试应至少覆盖：

- 缺少参数时返回 `400`
- 参数正确时返回 `200`

你要开始形成一个习惯：

风险最高的地方往往不是 happy path，而是边界和错误路径。

---

## 6. 集成测试和单元测试的边界

单元测试不是万能的。

这些更适合集成测试：

- 数据库事务
- Redis 缓存回源
- 鉴权中间件 + handler
- 异步任务状态流转

例如你要验证：

- 创建用户后是否真的写入数据库
- 更新后缓存是否真的失效

这种时候只 mock 掉一切，往往测不到真实问题。

所以更现实的理解是：

- 单元测试测规则
- 集成测试测协作

---

## 7. 文档不是可选装饰

一个能跑但没人知道怎么启动的项目，工程价值很低。

最小文档集合通常至少包括：

- `README.md`
- `.env.example`
- API 说明
- 数据库说明
- 部署说明
- 错误码说明

你可以先把 `README.md` 做到最小可用：

```md
# bookmark-service

## Run

```bash
go run ./cmd/server
```

## Env

- DATABASE_URL=
- REDIS_ADDR=

## Test

```bash
go test ./...
```
```

重点不是写长，而是让别人真的能照着跑起来。

---

## 8. `.env.example` 为什么重要

它的作用不是存秘密，而是说明项目需要什么配置。

例如：

```env
APP_PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/app?sslmode=disable
REDIS_ADDR=localhost:6379
JWT_SECRET=change-me
```

你要避免两种坏情况：

- 配置全靠口口相传
- `.env` 里藏了真实密钥还直接提交

---

## 9. 最小质量门槛应该是什么

小项目一开始不需要搞得非常复杂，但至少要有这些：

- 提交前跑 `gofmt`
- 提交前跑 `go test ./...`
- 关键业务模块有测试
- 接口变更时同步更新文档

一个很实用的日常顺序是：

```bash
gofmt -w .
go test ./...
```

如果项目再进一层，可以增加：

- `go vet ./...`
- CI 自动检查

---

## 10. 测试策略怎么想

不是所有代码都要同等力度测试。

你可以优先测：

- 容易出错的规则
- 历史上改动频繁的模块
- 一旦出错影响很大的流程

例如：

- 登录
- 权限判断
- 订单状态更新
- 异步任务去重

而不是把时间先花在低风险样板代码上。

---

## 11. 这一章建议你手敲的练习

### 练习 1：纯函数单元测试

写一个带错误分支的函数，再补表驱动测试。

### 练习 2：handler 测试

给一个带 query 参数校验的接口补 `httptest`。

### 练习 3：README 最小化整理

把启动、环境变量、测试命令写清楚。

### 练习 4：测试策略说明

为你当前项目写一段说明，解释：

- 哪些用单元测试
- 哪些用集成测试
- 哪些先不测

---

## 12. 这章学完的最低标准

你至少应该能做到：

- 为核心业务规则写表驱动测试
- 为 handler 写 `httptest`
- 区分单元测试和集成测试的职责
- 维护最基本的启动与配置文档
- 建立 `gofmt + go test` 的最小质量门槛

做到这里，你的项目才开始从“能跑”变成“可维护”。
