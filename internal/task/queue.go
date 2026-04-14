package task

import (
	"fmt"
	"go-mvc/config"
	"go-mvc/pkg/queue"
)

// Init 初始化任务队列
func Init() error {
	v := config.GetViper()
	addr := fmt.Sprintf("%s:%d", v.GetString("redis.host"), v.GetInt("redis.port"))

	// 一行搞定初始化
	queue.Init(
		addr,
		v.GetString("redis.password"),
		v.GetInt("redis.db"),
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
