package task

import (
	"go-mvc/pkg/queue"

	"github.com/spf13/viper"
)

// Init 初始化任务队列
func Init(v *viper.Viper) error {
	return queue.Init(v)
}

// StartQueue 启动任务队列
func StartQueue() error {
	return queue.Start()
}

// ShutdownQueue 关闭任务队列
func ShutdownQueue() error {
	return queue.Shutdown()
}
