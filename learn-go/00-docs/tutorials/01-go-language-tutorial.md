# 01 Go 语言基础教程

## 1. 怎么使用这篇教程

这篇文档不只讲概念，也给你可以直接跟着敲的最小例子。

建议你用下面节奏学习：

1. 先看概念解释，知道它是干什么的。
2. 把代码手敲一遍，不要直接复制。
3. 运行代码，观察输出。
4. 再改几行，故意试错。

如果你是从 JavaScript / TypeScript 转过来，重点不是背语法，而是把 Go 的思维方式建立起来：

- 显式
- 简单
- 少魔法
- 重视值语义和错误返回

---

## 2. 变量、常量和零值

### 2.1 `var` 和 `:=`

`var` 用于显式声明，`:=` 用于函数内部的快速定义。

先敲下面这段：

```go
package main

import "fmt"

func main() {
	var name string
	age := 18

	fmt.Println("name =", name)
	fmt.Println("age =", age)
}
```

运行：

```bash
go run main.go
```

你会看到类似输出：

```txt
name =
age = 18
```

这里最重要的不是 `age := 18`，而是 `var name string` 即使没赋值也能工作。这就是 Go 的零值设计。

### 2.2 零值是什么

再敲：

```go
package main

import "fmt"

func main() {
	var (
		n int
		ok bool
		s string
	)

	fmt.Println("int zero value:", n)
	fmt.Println("bool zero value:", ok)
	fmt.Println("string zero value:", s)
}
```

观察点：

- `int` 默认是 `0`
- `bool` 默认是 `false`
- `string` 默认是空字符串

这会直接影响你后面写结构体、接口响应、数据库模型时的默认行为。

### 2.3 `const` 和 `iota`

```go
package main

import "fmt"

const (
	StatusPending = iota
	StatusRunning
	StatusDone
)

func main() {
	fmt.Println(StatusPending, StatusRunning, StatusDone)
}
```

输出：

```txt
0 1 2
```

你可以把 `iota` 暂时理解成“帮你生成一组连续常量”。

---

## 3. 基础类型

你需要先熟悉这些：

- `bool`
- `int`
- `int64`
- `float64`
- `string`
- `byte`
- `rune`

### 3.1 一次把这些类型都看一遍

```go
package main

import "fmt"

func main() {
	var ok bool = true
	var count int = 10
	var userID int64 = 10000000001
	var price float64 = 19.99
	var text string = "Go"
	var b byte = 'A'
	var r rune = '中'

	fmt.Printf("ok=%v type=%T\n", ok, ok)
	fmt.Printf("count=%v type=%T\n", count, count)
	fmt.Printf("userID=%v type=%T\n", userID, userID)
	fmt.Printf("price=%v type=%T\n", price, price)
	fmt.Printf("text=%v type=%T\n", text, text)
	fmt.Printf("byte=%v type=%T\n", b, b)
	fmt.Printf("rune=%v type=%T\n", r, r)
}
```

观察点：

- `byte` 本质上是 `uint8`
- `rune` 本质上是 `int32`
- `rune` 更偏“字符”，`byte` 更偏“原始字节”

### 3.2 `byte` 和 `rune` 的区别

```go
package main

import "fmt"

func main() {
	s := "A中"

	fmt.Println("len(s) =", len(s))

	for i, b := range []byte(s) {
		fmt.Printf("byte[%d]=%v\n", i, b)
	}

	for i, r := range s {
		fmt.Printf("rune index=%d value=%c code=%d\n", i, r, r)
	}
}
```

观察点：

- `len(s)` 是字节长度，不是字符数
- 中文字符通常会占多个字节
- `range string` 取到的是 `rune`

### 3.3 显式类型转换

Go 不做隐式数值转换，这一点要尽快习惯。

```go
package main

import "fmt"

func main() {
	var a int = 10
	var b int64 = int64(a)
	var c float64 = float64(a)

	fmt.Printf("a=%v type=%T\n", a, a)
	fmt.Printf("b=%v type=%T\n", b, b)
	fmt.Printf("c=%v type=%T\n", c, c)
}
```

你可以故意删掉 `int64(...)` 试一次，编译器会直接报错。这是好事，因为它减少了很多模糊转换。

---

## 4. 数组、切片、map、结构体

### 4.1 数组和切片

```go
package main

import "fmt"

func main() {
	arr := [3]int{1, 2, 3}
	slice := []int{1, 2, 3}

	fmt.Printf("arr=%v len=%d\n", arr, len(arr))
	fmt.Printf("slice=%v len=%d cap=%d\n", slice, len(slice), cap(slice))
}
```

记住：

- 数组长度固定
- 切片是后端里真正高频使用的结构

### 4.2 `append` 和扩容

```go
package main

import "fmt"

func main() {
	nums := make([]int, 0, 2)
	fmt.Printf("init len=%d cap=%d %v\n", len(nums), cap(nums), nums)

	nums = append(nums, 1)
	nums = append(nums, 2)
	fmt.Printf("after 2 append len=%d cap=%d %v\n", len(nums), cap(nums), nums)

	nums = append(nums, 3)
	fmt.Printf("after 3 append len=%d cap=%d %v\n", len(nums), cap(nums), nums)
}
```

观察点：

- 前两次 `append` 可能不扩容
- 第三次开始通常会扩容
- `cap` 变化能帮助你理解切片底层行为

当你第三次 append 时，len 会从 2 变成 3，已经超过旧容量 2，Go 就会重新分配底层数组。Go 官方源码里 append 触发扩容时会走 growslice，而 nextslicecap 规定：当旧容量小于 256 时，新容量直接翻倍 doublecap。所以这里 2 -> 4。官方源码里能看到这个规则。[runtime/slice.go](https://go.dev/src/runtime/slice.go) 里有 threshold = 256，并且 oldCap < threshold 时返回 doublecap。

大切片增长会放缓，源码里是大约 1.25x

### 4.3 `map` 的读取行为

```go
package main

import "fmt"

func main() {
	scores := map[string]int{
		"go": 100,
	}

	fmt.Println(scores["go"])
	fmt.Println(scores["java"])

	v, ok := scores["java"]
	fmt.Println("value:", v, "exists:", ok)
}
```

观察点：

- 不存在的键会返回值类型的零值
- 所以很多时候必须用 `v, ok := m[key]`

### 4.4 结构体

```go
package main

import "fmt"

type User struct {
	ID    int64
	Name  string
	Email string
}

func main() {
	u := User{
		ID:    1,
		Name:  "Alice",
		Email: "alice@example.com",
	}

	fmt.Printf("%+v\n", u)
}
```

### 4.5 JSON tag

```go
package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func main() {
	u := User{ID: 1, Name: "Alice"}
	data, _ := json.Marshal(u)
	fmt.Println(string(data))
}
```

输出类似：

```txt
{"id":1,"name":"Alice"}
```

这和后面写 API 响应直接相关。

---

## 5. 控制流和函数

### 5.1 `for` 和 `range`

```go
package main

import "fmt"

func main() {
	nums := []string{"go", "redis", "sql"}

	for i := 0; i < len(nums); i++ {
		fmt.Println("for:", i, nums[i])
	}

	for i, v := range nums {
		fmt.Println("range:", i, v)
	}
}
```

### 5.2 多返回值

```go
package main

import (
	"errors"
	"fmt"
)

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func main() {
	result, err := divide(10, 2)
	fmt.Println(result, err)

	result, err = divide(10, 0)
	fmt.Println(result, err)
}
```

这段代码你要看到的是 Go 最典型的控制流：

```go
result, err := ...
if err != nil {
    ...
}
```

### 5.3 匿名函数和闭包

```go
package main

import "fmt"

func main() {
	base := 10
	add := func(v int) int {
		return base + v
	}

	fmt.Println(add(5))
}
```

这就是闭包最小例子：内部函数拿到了外部变量 `base`。

---

## 6. 指针、方法、接口

### 6.1 为什么需要指针

```go
package main

import "fmt"

type Counter struct {
	Value int
}

func incByValue(c Counter) {
	c.Value++
}

func incByPointer(c *Counter) {
	c.Value++
}

func main() {
	c := Counter{Value: 1}

	incByValue(c)
	fmt.Println("after incByValue:", c.Value)

	incByPointer(&c)
	fmt.Println("after incByPointer:", c.Value)
}
```

观察点：

- 值传递改不到原对象
- 指针传递可以改原对象

### 6.2 方法接收者

```go
package main

import "fmt"

type Counter struct {
	Value int
}

func (c *Counter) Inc() {
	c.Value++
}

func main() {
	c := Counter{Value: 10}
	c.Inc()
	fmt.Println(c.Value)
}
```

### 6.3 接口的最小例子

```go
package main

import "fmt"

type Speaker interface {
	Speak() string
}

type User struct {
	Name string
}

func (u User) Speak() string {
	return "hello, " + u.Name
}

func say(s Speaker) {
	fmt.Println(s.Speak())
}

func main() {
	u := User{Name: "Alice"}
	say(u)
}
```

重点不是接口语法本身，而是：

- `User` 没有写 `implements`
- 只要方法满足接口，就算实现了

---

## 7. 错误处理

### 7.1 最小错误处理链

```go
package main

import (
	"errors"
	"fmt"
)

func findUser(id int64) (string, error) {
	if id <= 0 {
		return "", errors.New("invalid user id")
	}
	return "alice", nil
}

func main() {
	name, err := findUser(0)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("name:", name)
}
```

这段代码的重点是：Go 倾向于把错误当普通返回值处理，而不是依赖异常。

---

## 8. 泛型最小例子

如果你已经学到这里，可以先看一个最小泛型例子，但不要急着深挖。

```go
package main

import "fmt"

func first[T any](items []T) T {
	return items[0]
}

func main() {
	fmt.Println(first([]int{1, 2, 3}))
	fmt.Println(first([]string{"go", "redis"}))
}
```

你只需要先理解：

- `T` 是类型参数
- 同一套逻辑可以作用于不同类型

---

## 9. 这一章建议你手敲的练习

### 练习 1：类型实验

自己写一个文件，把下面内容都打印出来：

- `bool`
- `int`
- `int64`
- `float64`
- `string`
- `byte`
- `rune`

并用 `%T` 打印类型。

### 练习 2：切片实验

做一个 `[]int`，连续 `append` 10 次，每次打印：

- 当前值
- `len`
- `cap`

### 练习 3：结构体 + JSON

定义一个 `Todo` 结构体：

- `id`
- `title`
- `done`

再把它编码成 JSON 输出。

### 练习 4：接口实验

定义一个 `Animal` 接口，再定义 `Dog` 和 `Cat` 两个结构体实现它。

---

## 10. 这章学完的最低标准

如果你已经能独立写出下面这些，说明基础够用了：

- 带结构体的 CLI 小程序
- 读写 JSON 的小程序
- 带错误返回的函数链
- 一个简单的接口示例

如果这些还不顺，先别急着上数据库和 Web 框架，继续把这章敲熟。
