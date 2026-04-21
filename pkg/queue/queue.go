package queue

import (
	"fmt"
	"strings"
	"sync"
	"time"

	queueprovider "go-mvc/pkg/queue/provider"

	"github.com/spf13/viper"
)

type Handler = queueprovider.Handler
type Option = queueprovider.Option

type registeredHandler struct {
	taskType string
	handler  Handler
}

func WithQueue(name string) Option {
	return queueprovider.WithQueue(name)
}

func WithMaxRetry(n int) Option {
	return queueprovider.WithMaxRetry(n)
}

func WithTimeout(timeout time.Duration) Option {
	return queueprovider.WithTimeout(timeout)
}

func WithDeadline(deadline time.Time) Option {
	return queueprovider.WithDeadline(deadline)
}

func WithUnique(ttl time.Duration) Option {
	return queueprovider.WithUnique(ttl)
}

func WithTaskID(id string) Option {
	return queueprovider.WithTaskID(id)
}

func WithRetention(retention time.Duration) Option {
	return queueprovider.WithRetention(retention)
}

// Start 启动任务队列
func Start() error {
	return defaultProvider.Start()
}

// Shutdown 关闭任务队列
func Shutdown() error {
	return defaultProvider.Shutdown()
}

// Close 关闭任务队列（统一组件协议入口）。
func Close() error {
	return Shutdown()
}

// Enqueue 立即执行任务
func Enqueue(taskType string, payload any, opts ...Option) error {
	return defaultProvider.Enqueue(taskType, payload, opts...)
}

// EnqueueIn 延迟执行任务
func EnqueueIn(taskType string, delay time.Duration, payload any, opts ...Option) error {
	return defaultProvider.EnqueueIn(taskType, delay, payload, opts...)
}

// EnqueueAt 指定时间执行
func EnqueueAt(taskType string, at time.Time, payload any, opts ...Option) error {
	return defaultProvider.EnqueueAt(taskType, at, payload, opts...)
}

func buildProvider(v *viper.Viper) (queueprovider.Provider, error) {
	providerName := "asynq"
	if v != nil {
		providerName = strings.TrimSpace(strings.ToLower(v.GetString("queue.provider")))
		if providerName == "" {
			providerName = "asynq"
		}
	}

	switch providerName {
	case "asynq":
		return queueprovider.NewAsynqProvider(), nil
	default:
		return nil, fmt.Errorf("不支持的队列 provider: %s", providerName)
	}
}

var defaultProvider queueprovider.Provider = queueprovider.NewAsynqProvider()
var (
	registrationMu sync.RWMutex
	registrations  []registeredHandler
)

// Init 初始化任务队列
func Init(v *viper.Viper) error {
	provider, err := buildProvider(v)
	if err != nil {
		return err
	}

	registrationMu.RLock()
	registered := append([]registeredHandler(nil), registrations...)
	registrationMu.RUnlock()

	for _, item := range registered {
		provider.Register(item.taskType, item.handler)
	}

	defaultProvider = provider
	if err := defaultProvider.Init(v); err != nil {
		return err
	}
	if v != nil && v.GetBool("queue.run_worker") {
		if err := defaultProvider.Start(); err != nil {
			return err
		}
	}
	return nil
}

// Register 注册处理器
func Register(taskType string, handler Handler) {
	registrationMu.Lock()
	registrations = append(registrations, registeredHandler{
		taskType: taskType,
		handler:  handler,
	})
	registrationMu.Unlock()
	defaultProvider.Register(taskType, handler)
}
