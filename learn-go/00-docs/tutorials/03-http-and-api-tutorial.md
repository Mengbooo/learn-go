# 03 HTTP 与 API 设计教程

## 1. 这一章怎么学

这章不要只看概念。你至少要自己启动一个最小 HTTP 服务，然后用 `curl` 去打它。

后端接口不是“返回个 JSON”就结束了，你要开始理解：

- 请求怎么进来
- 路由怎么分发
- 参数怎么解析
- 状态码怎么表达
- 错误怎么统一

---

## 2. 第一个 HTTP 服务

新建 `main.go`：

```go
package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello api")
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	fmt.Println("server running at :8080")
	http.ListenAndServe(":8080", nil)
}
```

运行：

```bash
go run main.go
```

另开一个终端执行：

```bash
curl http://localhost:8080/hello
```

你会看到：

```txt
hello api
```

这就是从“写程序”进入“写服务”的第一步。

---

## 3. 返回 JSON 响应

把 handler 改成这样：

```go
package main

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := Response{
		Code:    "OK",
		Message: "hello api",
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":8080", nil)
}
```

再执行：

```bash
curl http://localhost:8080/hello
```

你会得到 JSON。

观察点：

- 设置 `Content-Type`
- 用 `json.NewEncoder` 直接写响应

---

## 4. 认识请求方法

加一个创建接口：

```go
package main

import (
	"encoding/json"
	"net/http"
)

type Todo struct {
	Title string `json:"title"`
}

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input Todo
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"code":    "OK",
		"message": "created",
		"data":    input,
	})
}

func main() {
	http.HandleFunc("/todos", createTodoHandler)
	http.ListenAndServe(":8080", nil)
}
```

调用：

```bash
curl -X POST http://localhost:8080/todos \
  -H 'Content-Type: application/json' \
  -d '{"title":"learn go"}'
```

再故意试错：

```bash
curl http://localhost:8080/todos
```

你应该看到 `405 method not allowed`。

这会帮你建立“HTTP method 有明确语义”的意识。

---

## 5. 查询参数

```go
package main

import (
	"encoding/json"
	"net/http"
)

func listHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	keyword := r.URL.Query().Get("keyword")

	_ = json.NewEncoder(w).Encode(map[string]any{
		"page":    page,
		"keyword": keyword,
	})
}

func main() {
	http.HandleFunc("/search", listHandler)
	http.ListenAndServe(":8080", nil)
}
```

调用：

```bash
curl 'http://localhost:8080/search?page=1&keyword=go'
```

观察点：

- 查询参数不在 body 里
- 列表接口通常要从 query 取分页和筛选条件

---

## 6. 统一成功响应和错误响应

真实项目不要今天返回字符串、明天返回对象。先建立统一格式。

```go
package main

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, resp APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

func handler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	if title == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Code:    "VALIDATION_ERROR",
			Message: "title is required",
		})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Code:    "OK",
		Message: "success",
		Data: map[string]string{
			"title": title,
		},
	})
}

func main() {
	http.HandleFunc("/demo", handler)
	http.ListenAndServe(":8080", nil)
}
```

分别试：

```bash
curl 'http://localhost:8080/demo'
curl 'http://localhost:8080/demo?title=learn-go'
```

---

## 7. 为什么状态码重要

你可以直接把上面的例子拿来观察：

- 参数错误时返回 `400`
- 创建成功时返回 `201`
- 方法不允许时返回 `405`

最常见的一组你先记住：

- `200` 成功
- `201` 创建成功
- `400` 参数错误
- `401` 未登录
- `403` 无权限
- `404` 资源不存在
- `500` 服务内部错误

---

## 8. 最小中间件示例

Go 里中间件本质上就是“包装 handler”。

```go
package main

import (
	"fmt"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		fmt.Println(r.Method, r.URL.Path, time.Since(start))
	})
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", hello)

	http.ListenAndServe(":8080", loggingMiddleware(mux))
}
```

每请求一次，你的终端都会打印一条简单日志。

这就是后面做：

- 日志
- 认证
- 限流
- recover
- metrics

的起点。

---

## 9. 最小 Todo API

下面给你一个非常小的内存版 Todo API，可以完整跟敲。

```go
package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var todos = []Todo{
	{ID: 1, Title: "learn go", Done: false},
}

func listTodos(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var input Todo
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	input.ID = len(todos) + 1
	todos = append(todos, input)

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(input)
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	idText := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idText)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	for _, todo := range todos {
		if todo.ID == id {
			_ = json.NewEncoder(w).Encode(todo)
			return
		}
	}

	http.Error(w, "not found", http.StatusNotFound)
}

func main() {
	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			listTodos(w, r)
			return
		}
		if r.Method == http.MethodPost {
			createTodo(w, r)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/todo", getTodo)
	http.ListenAndServe(":8080", nil)
}
```

你可以这样测：

```bash
curl http://localhost:8080/todos
curl -X POST http://localhost:8080/todos -H 'Content-Type: application/json' -d '{"title":"build api","done":false}'
curl 'http://localhost:8080/todo?id=2'
```

---

## 10. `chi` 和 `Gin` 的生态位

在你把上面的标准库例子跑顺以后，再理解第三方框架会更容易。

### `chi`

适合：

- 想贴近标准库
- 想自己掌控结构

### `Gin`

适合：

- 想快速上手
- 需要更多现成资料

建议顺序仍然是：

```txt
net/http
  ↓
chi 或 Gin
  ↓
OpenAPI
```

---

## 11. 这一章建议你手敲的练习

### 练习 1：Hello API

写一个 `/ping` 接口，返回：

```json
{"message":"pong"}
```

### 练习 2：列表查询

写一个 `/users?page=1&page_size=10` 接口，把 query 参数原样返回。

### 练习 3：统一响应

自己实现一个 `writeJSON` 辅助函数。

### 练习 4：Todo API

补充一个“标记完成”的接口。

---

## 12. 这章学完的最低标准

你应该至少能独立做出这些：

- 启动一个 HTTP 服务
- 返回 JSON
- 解析 query 和 body
- 使用合理状态码
- 写一个最小中间件

做到这一步，才算真正进入后端开发。
