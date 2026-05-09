# 07 缓存 Redis 与性能优化教程

## 1. 这章怎么学

这章最容易学偏的地方，是一上来就把 Redis 当成“性能万能药”。

你要先建立一个更现实的顺序：

1. 先确认瓶颈在哪里。
2. 再判断缓存是否合适。
3. 最后才决定 key、TTL 和失效策略。

缓存能提升性能，但也会带来一致性、热点和排查复杂度。

---

## 2. Redis 到底解决什么问题

你可以先把 Redis 理解成一个“非常快的内存型数据服务”。

它常见的用途不是只有缓存，还包括：

- 热点数据读取
- 计数器
- 短期状态保存
- 限流
- 排行榜
- 分布式协调

但在小型 Go 后端项目里，最常见的第一步仍然是：

- 读多写少数据缓存
- 登录或验证码限流
- 异步任务状态存储

---

## 3. 最小 Redis 连接例子

先安装依赖：

```bash
go get github.com/redis/go-redis/v9
```

再写：

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	if err := rdb.Set(ctx, "site:name", "learn-go", 30*time.Second).Err(); err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "site:name").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("value:", val)
}
```

你先看懂三件事：

- Redis 操作也要带 `context`
- `Set` 可以顺手设置 TTL
- `Get` 失败时不一定都是系统错误，也可能是 key 不存在

---

## 4. 认识常见数据结构

### 4.1 String

最常见，用来存：

- 用户摘要缓存
- token 黑名单标记
- 简单计数

```go
err := rdb.Set(ctx, "user:1:name", "alice", 5*time.Minute).Err()
```

### 4.2 Hash

适合一个对象多个字段。

```go
err := rdb.HSet(ctx, "user:1", map[string]any{
	"name":  "alice",
	"email": "alice@example.com",
}).Err()
```

### 4.3 List

适合简单队列、消息列表。

```go
err := rdb.LPush(ctx, "jobs", "task-1").Err()
```

### 4.4 Set

适合集合去重。

```go
err := rdb.SAdd(ctx, "online_users", 1, 2, 3).Err()
```

### 4.5 Sorted Set

适合排行榜和带分数排序场景。

```go
err := rdb.ZAdd(ctx, "rank", redis.Z{Score: 100, Member: "alice"}).Err()
```

一开始不需要全背，先建立“业务形态和数据结构要匹配”的意识。

---

## 5. Cache Aside 最小例子

这是一种很常见的缓存模式：

1. 先查缓存
2. 没命中再查数据库
3. 查到后写回缓存

例子：

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserProfile struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func loadUserFromDB(userID int64) (UserProfile, error) {
	fmt.Println("load from db")
	return UserProfile{ID: userID, Name: "Alice"}, nil
}

func getUserProfile(ctx context.Context, rdb *redis.Client, userID int64) (UserProfile, error) {
	key := fmt.Sprintf("user:profile:%d", userID)

	cached, err := rdb.Get(ctx, key).Result()
	if err == nil {
		var profile UserProfile
		if err := json.Unmarshal([]byte(cached), &profile); err == nil {
			return profile, nil
		}
	}

	profile, err := loadUserFromDB(userID)
	if err != nil {
		return UserProfile{}, err
	}

	data, _ := json.Marshal(profile)
	_ = rdb.Set(ctx, key, data, 5*time.Minute).Err()

	return profile, nil
}
```

这里你要理解：

- 缓存命中时，不打数据库
- 未命中时，回源数据库并回填
- 回填失败通常不该影响主流程成功

---

## 6. 为什么常见做法是“写库后删缓存”

先看更新用户昵称的伪代码：

```go
func updateUserName(ctx context.Context, userID int64, name string) error {
	if err := updateUserNameInDB(ctx, userID, name); err != nil {
		return err
	}

	key := fmt.Sprintf("user:profile:%d", userID)
	_ = rdb.Del(ctx, key).Err()
	return nil
}
```

你可能会问，为什么不直接更新缓存？

因为小项目里“写库后删缓存”通常更稳：

- 数据源仍然以数据库为准
- 下次读取自动回源并回填
- 少一层并发覆盖写风险

它也不是绝对完美，但比“写库同时更新一堆缓存对象”更容易控制。

---

## 7. TTL 应该怎么想

TTL 不是随便拍个数字。

你至少要看四件事：

- 数据变更频率
- 数据被访问频率
- 能容忍多久不一致
- 回源数据库的成本

例如：

- 配置类数据：可以更长
- 用户资料摘要：中等 TTL
- 订单支付状态：通常不能粗暴长缓存

一个很常见但不严谨的错误是：

“这个接口慢，那就缓存 24 小时。”

这样可能只是把一致性问题延后爆炸。

---

## 8. 缓存穿透、击穿、雪崩

### 8.1 穿透

查的是根本不存在的数据。

例如有人反复查不存在的 `user_id`，每次都打到数据库。

基础思路：

- 参数校验先拦一层
- 对“确实不存在”的结果做短 TTL 空值缓存

### 8.2 击穿

某个热点 key 突然失效，大量请求同时回源数据库。

基础思路：

- 热点数据不要和大量 key 同时过期
- 必要时增加单飞控制或互斥回源

### 8.3 雪崩

大量 key 在同一时间集中失效。

基础思路：

- TTL 加随机抖动
- 热点和普通数据分层设置 TTL

例如：

```go
ttl := 5*time.Minute + time.Duration(time.Now().UnixNano()%30)*time.Second
```

这只是最小思路，但能先帮你建立风险意识。

---

## 9. 限流最小例子

最简单的短期限流可以先用计数器做。

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func allowRequest(ctx context.Context, rdb *redis.Client, key string, limit int64, window time.Duration) (bool, error) {
	count, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		_ = rdb.Expire(ctx, key, window).Err()
	}

	return count <= limit, nil
}

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	ok, err := allowRequest(ctx, rdb, "login:ip:127.0.0.1", 5, time.Minute)
	if err != nil {
		panic(err)
	}

	fmt.Println("allow:", ok)
}
```

这个版本不追求最严谨，但很适合建立第一层直觉：

- 某个窗口期内计数
- 超过阈值就拒绝

很常见的使用场景：

- 登录
- 验证码发送
- 公开搜索接口

---

## 10. 性能优化别只盯缓存

很多性能问题其实先该看这些：

- SQL 有没有索引
- 有没有 N+1 查询
- JSON 编解码是不是过重
- 外部调用是不是太慢
- 是否做了不必要的大对象拷贝

一个很现实的原则是：

先做能长期降低复杂度的优化，再引入额外基础设施。

例如：

- 补索引
- 减少重复查询
- 限制返回字段
- 加分页

通常都比“先塞 Redis”更干净。

---

## 11. 这一章建议你手敲的练习

### 练习 1：最小缓存读取

为一个“获取用户资料”函数加上 Redis 缓存。

### 练习 2：缓存失效

写一个更新用户昵称的函数，采用“写库后删缓存”。

### 练习 3：登录限流

按 IP 做一分钟 5 次的限制。

### 练习 4：缓存设计说明

选一个接口，说明：

- 为什么值得缓存
- key 怎么命名
- TTL 为什么这样定

---

## 12. 这章学完的最低标准

你至少应该能做到：

- 用 Go 连接 Redis 并完成基本读写
- 解释 Cache Aside 的基本流程
- 说清楚为什么小项目常用“写库后删缓存”
- 识别缓存穿透、击穿、雪崩
- 为一个真实接口设计 key 和 TTL

做到这里，你才算真的开始理解缓存，而不是只会调用 Redis API。
