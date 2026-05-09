# 04 数据库与分层架构教程

## 1. 这章的重点

从这里开始，你不再只是处理内存里的数据，而是开始处理“要长期保存、要能查、要能保证一致”的数据。

同时，代码结构也会开始膨胀，所以数据库和分层架构最好一起学。

---

## 2. 先看表设计最小例子

假设你在做书签系统，最小表可以是：

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE bookmarks (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

这里你先看懂四件事：

- `PRIMARY KEY` 表示主键
- `UNIQUE` 表示唯一约束
- `REFERENCES` 表示外键关系
- `user_id` 表示书签归属于某个用户

---

## 3. 最小 SQL 查询

### 3.1 插入数据

```sql
INSERT INTO users (email, password_hash)
VALUES ('alice@example.com', 'hashed-password');
```

### 3.2 查询数据

```sql
SELECT id, email, created_at
FROM users;
```

### 3.3 条件查询

```sql
SELECT id, title, url
FROM bookmarks
WHERE user_id = 1;
```

### 3.4 排序

```sql
SELECT id, title
FROM bookmarks
WHERE user_id = 1
ORDER BY created_at DESC;
```

### 3.5 分页

```sql
SELECT id, title
FROM bookmarks
WHERE user_id = 1
ORDER BY created_at DESC
LIMIT 20 OFFSET 0;
```

这里已经出现了后端列表接口最常见的 SQL 形态。

---

## 4. 索引最小例子

如果你的常见查询是：

```sql
SELECT id, title
FROM bookmarks
WHERE user_id = 1
ORDER BY created_at DESC;
```

那你就应该开始想索引：

```sql
CREATE INDEX idx_bookmarks_user_id ON bookmarks(user_id);
```

如果你经常同时按 `user_id` 过滤并按 `created_at` 排序，可以继续思考联合索引：

```sql
CREATE INDEX idx_bookmarks_user_id_created_at
ON bookmarks(user_id, created_at DESC);
```

你要建立的不是“索引越多越好”，而是“索引要服务明确查询”。

---

## 5. 事务最小例子

假设你要：

1. 创建书签
2. 写标签关系

这两个动作应该一起成功或一起失败。

SQL 版事务：

```sql
BEGIN;

INSERT INTO bookmarks (user_id, title, url)
VALUES (1, 'Go 官网', 'https://go.dev');

INSERT INTO bookmark_tags (bookmark_id, tag_id)
VALUES (100, 2);

COMMIT;
```

如果中间某一步失败，就应该 `ROLLBACK`。

---

## 6. Go 里连接数据库的最小感觉

下面是一个最小 `pgx` 示例。你现在先看结构，不必急着一次记住所有细节。

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	var now string
	err = conn.QueryRow(ctx, "select now()::text").Scan(&now)
	if err != nil {
		panic(err)
	}

	fmt.Println(now)
}
```

运行前你需要先准备 `DATABASE_URL`。

你要先看懂：

- 数据库操作也要带 `context`
- `QueryRow().Scan()` 是很常见的查询模式

---

## 7. 插入和查询的最小代码

### 7.1 插入用户

```go
func createUser(ctx context.Context, conn *pgx.Conn, email, passwordHash string) error {
	_, err := conn.Exec(ctx,
		`INSERT INTO users (email, password_hash) VALUES ($1, $2)`,
		email, passwordHash,
	)
	return err
}
```

### 7.2 查询单个用户

```go
type User struct {
	ID    int64
	Email string
}

func getUserByEmail(ctx context.Context, conn *pgx.Conn, email string) (User, error) {
	var u User
	err := conn.QueryRow(ctx,
		`SELECT id, email FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email)
	return u, err
}
```

重点观察：

- 参数不是字符串拼接，而是占位符
- 这样能避免 SQL 注入

---

## 8. Repository 的最小形态

你不应该把所有 SQL 直接写在 handler 里。先感受一下 repository 的作用。

```go
package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type User struct {
	ID    int64
	Email string
}

type UserRepository struct {
	conn *pgx.Conn
}

func NewUserRepository(conn *pgx.Conn) *UserRepository {
	return &UserRepository{conn: conn}
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := r.conn.QueryRow(ctx,
		`SELECT id, email FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email)
	return u, err
}
```

repository 层负责的是：

- SQL
- 数据持久化
- 数据库细节

---

## 9. Service 的最小形态

```go
package service

import "context"

type User struct {
	ID    int64
	Email string
}

type UserRepo interface {
	GetByEmail(ctx context.Context, email string) (User, error)
}

type UserService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) FindUser(ctx context.Context, email string) (User, error) {
	return s.repo.GetByEmail(ctx, email)
}
```

这里你开始感受分层：

- handler 处理 HTTP
- service 处理业务
- repository 处理数据库

---

## 10. Handler 的最小形态

```go
package handler

import (
	"encoding/json"
	"net/http"
)

type User struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type UserService interface {
	FindUser(rctx context.Context, email string) (User, error)
}
```

上面只是形态示意。你需要建立的意识是：

- handler 不该直接写 SQL
- service 不该直接处理 HTTP 细节

---

## 11. N+1 查询怎么理解

假设你先查出 20 个书签，然后循环 20 次去查每个书签的标签：

```txt
1 次查列表
+ 20 次查标签
```

这就是典型的 N+1 问题。

你需要警惕这种写法：

```go
for _, b := range bookmarks {
	tags := queryTagsByBookmarkID(b.ID)
	...
}
```

看起来自然，实际上可能很慢。

---

## 12. `sqlc` 的生态位

如果你觉得自己手写扫描代码很多，可以开始了解 `sqlc`。

它的作用不是替你理解 SQL，而是：

- 保留 SQL 可见性
- 帮你生成类型安全代码
- 减少重复样板

所以正确顺序是：

```txt
先理解 SQL
  ↓
再用 sqlc 提效
```

---

## 13. 这一章建议你手敲的练习

### 练习 1：用户表

自己写出：

- `users` 表
- `email` 唯一约束

### 练习 2：书签表

自己写出：

- `bookmarks` 表
- `user_id` 外键

### 练习 3：查询语句

写三条 SQL：

- 按邮箱查用户
- 按用户查书签
- 书签列表按创建时间倒序

### 练习 4：分层重构

把一个“查用户”的功能拆成：

- handler
- service
- repository

---

## 14. 这章学完的最低标准

你应该至少能做到：

- 看懂简单表结构
- 写基本 CRUD SQL
- 理解事务为什么需要
- 知道为什么要把 SQL 从 handler 拿出去

做到这一步，数据库和工程化才算真正入门。
