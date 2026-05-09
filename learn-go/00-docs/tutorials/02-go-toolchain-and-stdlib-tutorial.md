# 02 Go 工具链与标准库教程

## 1. 这一章学什么

语言基础只是“能写代码”，工具链决定“你能不能稳定做项目”。

这一章你要掌握三件事：

- 项目怎么初始化和运行
- 代码怎么格式化和检查
- 标准库怎么支撑后端基础能力

---

## 2. `go mod` 和项目初始化

### 2.1 初始化一个项目

先在一个空目录里执行：

```bash
go mod init example.com/hello-go
```

这会生成一个 `go.mod`。

最小项目：

```txt
hello-go/
  go.mod
  main.go
```

`main.go` 写成这样：

```go
package main

import "fmt"

func main() {
	fmt.Println("hello go module")
}
```

运行：

```bash
go run main.go
```

或者：

```bash
go run .
```

### 2.2 `go mod tidy`

当你增加或删除依赖后，执行：

```bash
go mod tidy
```

它负责：

- 补齐代码里实际用到的依赖
- 清理代码里已经不用的依赖

这一步以后会成为你的日常动作。

---

## 3. 运行、构建、安装

### 3.1 `go run`

适合开发阶段快速执行：

```bash
go run .
```

### 3.2 `go build`

适合生成产物：

```bash
go build -o bin/app .
```

执行完后你会看到：

```txt
bin/
  app
```

再运行它：

```bash
./bin/app
```

### 3.3 `go install`

适合安装命令行工具。

比如未来你会用：

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

先记住生态位：

- `run` 是临时执行
- `build` 是产出二进制
- `install` 是安装工具

---

## 4. 格式化和基础检查

### 4.1 `gofmt`

先写一个格式很乱的文件：

```go
package main
import "fmt"
func main( ){
fmt.Println("bad format")
}
```

执行：

```bash
gofmt -w main.go
```

你会发现它被整理成标准格式。

这是 Go 社区非常重要的习惯：少争论风格，多统一输出。

### 4.2 `go vet`

执行：

```bash
go vet ./...
```

它会帮你发现一些“编译能过，但写法很可疑”的问题。

### 4.3 `go test`

等会我们先写一个最小测试。

新建 `math.go`：

```go
package main

func add(a, b int) int {
	return a + b
}
```

新建 `math_test.go`：

```go
package main

import "testing"

func TestAdd(t *testing.T) {
	got := add(1, 2)
	want := 3

	if got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}
```

执行：

```bash
go test ./...
```

如果通过，你会看到 `ok`。

---

## 5. 表驱动测试最小例子

Go 很常见的测试写法是表驱动。

把 `math_test.go` 改成：

```go
package main

import "testing"

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{name: "1+2", a: 1, b: 2, want: 3},
		{name: "0+0", a: 0, b: 0, want: 0},
		{name: "-1+1", a: -1, b: 1, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := add(tt.a, tt.b)
			if got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}
```

这是后面写 service 测试时的主流模式。

---

## 6. 标准库 `encoding/json`

新建 `json_demo.go`：

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Todo struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func main() {
	todo := Todo{
		ID:    1,
		Title: "learn go",
		Done:  false,
	}

	data, err := json.Marshal(todo)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
```

运行后你会看到 JSON 字符串。

你可以再试一次反序列化：

```go
package main

import (
	"encoding/json"
	"fmt"
)

type Todo struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func main() {
	raw := `{"id":1,"title":"learn go","done":true}`

	var todo Todo
	if err := json.Unmarshal([]byte(raw), &todo); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", todo)
}
```

---

## 7. 标准库 `context`

很多人第一次看 `context` 会觉得抽象，先用超时例子理解。

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	select {
	case <-time.After(2 * time.Second):
		fmt.Println("finished work")
	case <-ctx.Done():
		fmt.Println("context cancelled:", ctx.Err())
	}
}
```

运行后你应该看到超时，而不是 `finished work`。

这就是后面请求超时、数据库超时、第三方调用超时控制的基础。

---

## 8. 标准库 `os` 和 `io`

### 8.1 读取文件

先新建一个 `demo.txt`：

```txt
hello go
```

再写：

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	data, err := os.ReadFile("demo.txt")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
```

### 8.2 写入文件

```go
package main

import "os"

func main() {
	err := os.WriteFile("output.txt", []byte("write by go"), 0644)
	if err != nil {
		panic(err)
	}
}
```

这类能力后面会用于配置、日志、导出文件等场景。

---

## 9. 标准库 `time`

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	fmt.Println("now:", now.Format(time.RFC3339))

	expireAt := now.Add(30 * time.Minute)
	fmt.Println("expireAt:", expireAt.Format(time.RFC3339))
}
```

后端里几乎所有系统都离不开时间：

- token 过期
- 任务超时
- 日志时间戳
- 缓存 TTL

---

## 10. 标准库 `log/slog`

```go
package main

import (
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("server started", "port", 8080, "env", "dev")
}
```

观察点：

- 这不是普通字符串日志
- 它带结构化字段

后面做生产项目时，这类日志会更容易被检索和分析。

---

## 11. 包和目录的最小感觉

你可以做一个这样的小结构：

```txt
demo/
  go.mod
  main.go
  internal/
    greet/
      greet.go
```

`internal/greet/greet.go`：

```go
package greet

func Hello(name string) string {
	return "hello, " + name
}
```

`main.go`：

```go
package main

import (
	"fmt"

	"example.com/demo/internal/greet"
)

func main() {
	fmt.Println(greet.Hello("alice"))
}
```

这个例子先让你感受：

- 包怎么拆
- 代码怎么组织
- 为什么不是所有内容都塞进 `main.go`

---

## 12. 这一章建议你手敲的练习

### 练习 1：项目初始化

自己新建一个目录，完成：

- `go mod init`
- `main.go`
- `go run .`
- `go build -o bin/app .`

### 练习 2：JSON 编解码

定义一个 `User` 结构体，完成：

- `Marshal`
- `Unmarshal`

### 练习 3：`context` 超时

写一个 2 秒任务，但超时只给 1 秒。

### 练习 4：最小测试

自己写一个 `multiply(a, b int)`，再补一个测试文件。

---

## 13. 这章学完的最低标准

你应该至少能独立做出这些事：

- 初始化一个 Go 项目
- 跑通 `go run`、`go build`、`go test`
- 理解 `go mod tidy` 在做什么
- 用 `encoding/json`、`context`、`time`、`os` 写小程序

如果这些还很生疏，后面学 Web 和数据库会一直卡在工程细节上。
