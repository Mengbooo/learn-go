# 08 异步任务 消息队列与幂等教程

## 1. 这章的重点

这章不是让你马上上 Kafka 或 RabbitMQ，而是先把一个基础问题想清楚：

哪些事情不应该堵在用户请求里同步做完。

你要学会把“用户请求返回”和“后台任务完成”拆开。

---

## 2. 什么场景适合异步

这些都很典型：

- 发送邮件
- 生成报告
- 导入导出
- 调第三方 AI 或 OCR
- 上传后转码或分析

共同点通常是：

- 耗时不可控
- 失败后可能要重试
- 用户不需要立刻拿到最终结果

如果你把这些都塞进 HTTP 请求里，常见后果就是：

- 接口很慢
- 超时变多
- 用户体验不稳定
- 失败恢复困难

---

## 3. 先从最小后台任务模型开始

你不需要一上来就上完整队列，先用内存版理解模型。

```go
package main

import (
	"fmt"
	"time"
)

func createReportTask(reportID int64) {
	go func() {
		fmt.Println("start report task", reportID)
		time.Sleep(3 * time.Second)
		fmt.Println("finish report task", reportID)
	}()
}

func main() {
	createReportTask(1)
	fmt.Println("request returned quickly")
	time.Sleep(5 * time.Second)
}
```

这段代码帮你理解最核心的一步：

- 请求先返回
- 耗时任务在后台继续执行

但它也有明显问题：

- 进程挂了任务就没了
- 任务状态不可查
- 失败没法重试

这正是你下一步要补的能力。

---

## 4. 最小任务状态机

异步任务一旦进入真实项目，就必须可追踪。

常见状态：

- `pending`
- `processing`
- `completed`
- `failed`
- `retrying`

最小结构：

```go
package main

type Job struct {
	ID        string
	Type      string
	Status    string
	Payload   string
	Retry     int
	ErrorText string
}
```

你要建立一个强意识：

没有状态机的异步任务，后面排查问题会非常痛苦。

---

## 5. 任务创建和状态查询最小例子

下面用内存 map 先模拟：

```go
package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Job struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

var (
	jobs = map[string]*Job{}
	mu   sync.Mutex
	seq  int
)

func createJob(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	seq++
	id := strconv.Itoa(seq)
	job := &Job{ID: id, Status: "pending"}
	jobs[id] = job
	mu.Unlock()

	go func(jobID string) {
		mu.Lock()
		jobs[jobID].Status = "processing"
		mu.Unlock()

		time.Sleep(2 * time.Second)

		mu.Lock()
		jobs[jobID].Status = "completed"
		mu.Unlock()
	}(id)

	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(job)
}

func getJob(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	mu.Lock()
	job, ok := jobs[id]
	mu.Unlock()

	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(job)
}

func main() {
	http.HandleFunc("/jobs", createJob)
	http.HandleFunc("/job", getJob)
	http.ListenAndServe(":8080", nil)
}
```

测试：

```bash
curl -X POST http://localhost:8080/jobs
curl 'http://localhost:8080/job?id=1'
```

你现在先关注异步 API 的形态：

- 创建任务返回 `202 Accepted`
- 返回任务 ID
- 后续通过查询接口看状态

---

## 6. 队列和 Worker 的基本关系

你可以先这样理解：

- 生产者：创建任务
- 队列：临时保存任务
- Worker：从队列里取任务并执行

最小伪代码：

```txt
HTTP Handler
  -> push job
  -> return job id

Worker
  -> pop job
  -> execute
  -> update status
```

这套拆分的价值有三个：

- 请求延迟更稳定
- 后台处理可以独立扩容
- 失败恢复更容易设计

---

## 7. 失败重试最小思路

重试不是“失败了就无限再来”。

你至少要有：

- 最大重试次数
- 重试间隔
- 最终失败状态

例子：

```go
func process(job *Job) {
	job.Status = "processing"

	err := doWork(job)
	if err == nil {
		job.Status = "completed"
		return
	}

	job.Retry++
	if job.Retry >= 3 {
		job.Status = "failed"
		job.ErrorText = err.Error()
		return
	}

	job.Status = "retrying"
}
```

这里最重要的不是代码细节，而是边界：

- 重试几次
- 什么错误可重试
- 最终失败后怎么暴露给用户或运营

---

## 8. 幂等为什么是异步系统基础能力

很多人第一次做队列时会默认认为：

“一个任务只会被消费一次。”

这在真实系统里通常不可靠。

重复消费常见来源：

- worker 执行完但来不及确认
- 消费后进程崩溃
- 网络抖动导致重复投递

所以正确思路应该是：

即使同一个任务被执行多次，结果也不能错。

---

## 9. 幂等最小例子

假设业务是“给用户发欢迎券”，那你不能每次重试都再发一张。

你至少可以这样做：

```go
package main

import "fmt"

var sent = map[string]bool{}

func sendCoupon(taskID string, userID int64) error {
	if sent[taskID] {
		fmt.Println("already processed:", taskID)
		return nil
	}

	fmt.Println("issue coupon to user", userID)
	sent[taskID] = true
	return nil
}
```

真实项目里更常见的方案有：

- 数据库唯一约束
- 幂等请求 ID
- 消费记录表
- 结果覆盖写而不是追加写

你要逐步建立一个原则：

幂等通常要靠业务数据层保证，而不是只靠队列中间件。

---

## 10. Redis 队列和 Asynq 的生态位

### 10.1 Redis List

适合学习最小模型。

优点：

- 直观
- 够简单

缺点：

- 自己要补很多状态和重试逻辑

### 10.2 Asynq

很适合 Go 小到中型项目。

优点：

- 天然基于 Redis
- 已经提供任务、重试、延迟、调度等能力

### 10.3 RabbitMQ / Kafka

更适合复杂场景，但学习和运维成本也更高。

建议顺序仍然是：

```txt
goroutine 理解模型
  ↓
Redis + worker
  ↓
Asynq
  ↓
RabbitMQ / Kafka
```

---

## 11. 异步接口怎么设计

如果任务不是立即完成，就不要假装同步成功。

更合理的形态通常是：

### 创建任务

```http
POST /exports
```

返回：

```json
{
  "job_id": "123",
  "status": "pending"
}
```

状态码：

```txt
202 Accepted
```

### 查询任务

```http
GET /exports/123
```

返回：

```json
{
  "job_id": "123",
  "status": "processing"
}
```

这样用户和前端都能正确理解系统行为：任务已接收，但还没完成。

---

## 12. 这一章建议你手敲的练习

### 练习 1：最小后台任务

把一个耗时 3 秒的动作从 HTTP 请求里拆出去。

### 练习 2：任务状态查询

创建任务时返回任务 ID，再加一个查询状态接口。

### 练习 3：失败重试

模拟一个前两次失败、第三次成功的任务。

### 练习 4：幂等保护

为“发送通知”场景加一个防重复执行保护。

---

## 13. 这章学完的最低标准

你至少应该能做到：

- 判断什么业务适合同步，什么适合异步
- 设计最小任务状态机
- 用 `202 Accepted` 设计异步接口
- 理解失败重试的边界
- 解释为什么消费端必须考虑幂等

做到这里，你就开始真正接近可用的后台任务系统了。
