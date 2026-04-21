package task

import (
	"sync"

	"go-mvc/pkg/queue"
)

var registerOnce sync.Once

// RegisterHandlers 显式注册项目内的任务处理器。
func RegisterHandlers() {
	registerOnce.Do(func() {
		queue.Register(TypeEmailSend, HandleEmailSend)
		queue.Register(TypeOrderCancel, HandleOrderCancel)
	})
}
