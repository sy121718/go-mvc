package task

import (
	"go-mvc/config"
	"go-mvc/pkg/queue"
)

// Init 初始化任务队列
func Init() error {
	cfg := config.GetRedis()

	// 一行搞定初始化
	queue.Init(
		cfg.Host+":6379",
		cfg.Password,
		cfg.DB,
		10, // 并发数
	)

	return nil
}

// StartQueue 启动任务队列
func StartQueue() {
	queue.Start()
}

// ShutdownQueue 关闭任务队列
func ShutdownQueue() {
	queue.Shutdown()
}
