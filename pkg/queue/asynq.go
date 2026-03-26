// Package queue /*
package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

// Client 客户端
type Client struct {
	client *asynq.Client
}

// Server 服务器
type Server struct {
	server *asynq.Server
	mux    *asynq.ServeMux
}

var (
	defaultClient *Client
	defaultServer *Server
	registerFuncs []func()
)

// ========== 初始化 ==========

// Init 初始化客户端和服务器
func Init(redisAddr, redisPassword string, redisDB int, concurrency int) {
	// 初始化客户端
	defaultClient = &Client{
		client: asynq.NewClient(asynq.RedisClientOpt{
			Addr:     redisAddr,
			Password: redisPassword,
			DB:       redisDB,
		}),
	}

	// 初始化服务器
	defaultServer = &Server{
		mux: asynq.NewServeMux(),
	}
	defaultServer.server = asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     redisAddr,
			Password: redisPassword,
			DB:       redisDB,
		},
		asynq.Config{
			Concurrency: concurrency,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	// 执行自动注册
	for _, fn := range registerFuncs {
		fn()
	}
}

// Start 启动服务器
func Start() {
	defaultServer.server.Start(defaultServer.mux)
}

// Shutdown 关闭
func Shutdown() {
	if defaultServer != nil {
		defaultServer.server.Shutdown()
	}
	if defaultClient != nil {
		defaultClient.client.Close()
	}
}

// ========== 客户端方法 ==========

// Enqueue 立即执行任务
func Enqueue(taskType string, payload interface{}, opts ...asynq.Option) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(taskType, data)
	_, err = defaultClient.client.Enqueue(task, opts...)
	return err
}

// EnqueueIn 延迟执行任务
func EnqueueIn(taskType string, delay time.Duration, payload interface{}, opts ...asynq.Option) error {
	opts = append(opts, asynq.ProcessIn(delay))
	return Enqueue(taskType, payload, opts...)
}

// EnqueueAt 指定时间执行
func EnqueueAt(taskType string, at time.Time, payload interface{}, opts ...asynq.Option) error {
	opts = append(opts, asynq.ProcessAt(at))
	return Enqueue(taskType, payload, opts...)
}

// ========== 服务器方法 ==========

// RegisterHandler 注册处理器（自动注册用）
func RegisterHandler(taskType string, handler func(ctx context.Context, t *asynq.Task) error) {
	defaultServer.mux.HandleFunc(taskType, handler)
}

// ========== 自动注册机制 ==========

// Register 自动注册函数（内部使用）
func Register(fn func()) {
	registerFuncs = append(registerFuncs, fn)
}
