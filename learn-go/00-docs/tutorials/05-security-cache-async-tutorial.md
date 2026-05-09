# 05 安全 缓存与异步任务教程

## 1. 这章讲什么

这章的内容开始明显接近真实线上系统：

- 用户怎么登录
- 资源怎么做权限控制
- 数据什么时候该缓存
- 耗时任务怎么放后台

你要开始从“功能能跑”过渡到“系统在复杂场景下还能稳住”。

---

## 2. 密码哈希最小例子

先安装依赖：

```bash
go get golang.org/x/crypto/bcrypt
```

再写：

```go
package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "123456"

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	fmt.Println("hash:", string(hash))

	err = bcrypt.CompareHashAndPassword(hash, []byte("123456"))
	fmt.Println("compare correct password:", err == nil)

	err = bcrypt.CompareHashAndPassword(hash, []byte("wrong"))
	fmt.Println("compare wrong password:", err == nil)
}
```

你应该看到：

- 哈希值不是明文
- 正确密码校验通过
- 错误密码校验失败

这就是“密码不能明文存储”的最小实践。

---

## 3. JWT 最小例子

先安装依赖：

```bash
go get github.com/golang-jwt/jwt/v5
```

代码：

```go
package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	secret := []byte("my-secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 1,
		"role":    "user",
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		panic(err)
	}

	fmt.Println("token:", tokenString)

	parsed, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return secret, nil
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("parsed valid:", parsed.Valid)
}
```

这段代码先帮你建立 token 的最小感觉：

- 签发
- 携带声明
- 校验

---

## 4. 资源权限最小例子

假设用户只能看自己的资源：

```go
package main

import "fmt"

type User struct {
	ID   int64
	Role string
}

type Bookmark struct {
	ID     int64
	UserID int64
	Title  string
}

func canAccess(user User, bookmark Bookmark) bool {
	if user.Role == "admin" {
		return true
	}
	return user.ID == bookmark.UserID
}

func main() {
	user := User{ID: 1, Role: "user"}
	bookmark := Bookmark{ID: 10, UserID: 2, Title: "Go"}

	fmt.Println(canAccess(user, bookmark))
}
```

你要理解的是：登录成功不代表什么都能访问。认证和鉴权不是一回事。

---

## 5. Redis 最小使用例子

先安装依赖：

```bash
go get github.com/redis/go-redis/v9
```

还要先本地有 Redis。

代码：

```go
package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	err := rdb.Set(ctx, "hello", "go-redis", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "hello").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println(val)
}
```

这只是最小读写，但已经足够帮你建立 Redis 的第一层直觉。

---

## 6. 缓存最小例子

下面是一个很简化的 Cache Aside 思路：

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func getUserProfile(ctx context.Context, rdb *redis.Client, userID string) (string, error) {
	cacheKey := "user_profile:" + userID

	val, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		fmt.Println("from cache")
		return val, nil
	}

	fmt.Println("from database")
	dbValue := "alice profile"

	if err := rdb.Set(ctx, cacheKey, dbValue, 5*time.Minute).Err(); err != nil {
		return "", err
	}

	return dbValue, nil
}
```

重点观察：

- 先查缓存
- 没命中再查数据库
- 再回填缓存

这就是小项目最常见的缓存策略。

---

## 7. 限流最小思路

伪代码先理解：

```txt
key = "login_limit:{ip}"
count = INCR key
if count == 1:
  EXPIRE key 60
if count > 5:
  reject
```

这类逻辑非常适合放在 Redis，因为它：

- 快
- 自带原子计数
- 适合短期状态

---

## 8. goroutine 后台任务最小例子

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("background job done")
	}()

	fmt.Println("request returned quickly")
	time.Sleep(3 * time.Second)
}
```

这能帮助你理解“后台执行”最基本的形态。

但你也要马上知道：这不等于可靠队列。进程挂了，任务就可能丢。

---

## 9. Asynq 的生态位和最小感觉

如果你继续做严肃一点的异步任务，可以看 Asynq。

安装：

```bash
go get github.com/hibiken/asynq
```

生产者最小示意：

```go
task := asynq.NewTask("email:send", []byte(`{"user_id":1}`))
info, err := client.Enqueue(task)
```

消费者最小示意：

```go
srv := asynq.NewServer(...)
mux := asynq.NewServeMux()
mux.HandleFunc("email:send", func(ctx context.Context, t *asynq.Task) error {
	return nil
})
```

这里你先建立概念就够了：

- 生产者入队
- worker 消费
- 失败可以重试

---

## 10. 幂等最小例子

假设一个任务可能被重复消费，你不能重复发奖励。

最简单思路是数据库唯一约束：

```sql
CREATE TABLE reward_records (
    id BIGSERIAL PRIMARY KEY,
    request_id TEXT NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    amount INT NOT NULL
);
```

业务含义：

- 第一次插入成功
- 第二次同一个 `request_id` 再插就失败

这就是一个非常常见的幂等手段。

---

## 11. 文档处理任务最小状态机

你可以把一个任务先设计成：

```txt
pending
processing
completed
failed
```

更新流程：

```txt
创建任务 -> pending
worker 开始处理 -> processing
成功 -> completed
失败 -> failed
```

为什么要有状态机：

- 用户能查状态
- 你能排查流程卡在哪
- 重试逻辑有落点

---

## 12. 这一章建议你手敲的练习

### 练习 1：密码哈希

自己输入一个密码，生成 hash，再校验两次：

- 一次正确密码
- 一次错误密码

### 练习 2：JWT

签发一个带 `user_id` 的 token，再解析它。

### 练习 3：Redis

本地启动 Redis，手敲：

- `SET`
- `GET`
- 带 TTL 的 `SET`

### 练习 4：后台任务

写一个 goroutine，模拟“生成报告”任务。

---

## 13. 这章学完的最低标准

你应该至少能理解：

- 为什么密码不能明文存储
- JWT 的基本作用是什么
- 缓存不是主数据库
- goroutine 不等于可靠任务系统
- 幂等是“重复也不出错”

做到这一步，后面接异步任务和真实登录系统就不会只停留在概念层。
