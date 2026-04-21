package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AppRuntime 描述一次应用启动后的核心运行时对象。
//
// 当前阶段收纳：
// - Server: 服务配置
// - Router: Gin 路由引擎
// - HTTPServer: HTTP 服务实例
//
// 后续如果继续从全局状态向运行时容器演进，可以逐步把更多运行时依赖收进这里。
type AppRuntime struct {
	Server     ServerConfig
	Router     *gin.Engine
	HTTPServer *http.Server
	ready      func() error
	shutdown   func() error
}

// NewAppRuntime 创建应用运行时容器。
func NewAppRuntime(
	server ServerConfig,
	router *gin.Engine,
	httpServer *http.Server,
	ready func() error,
	shutdown func() error,
) *AppRuntime {
	return &AppRuntime{
		Server:     server,
		Router:     router,
		HTTPServer: httpServer,
		ready:      ready,
		shutdown:   shutdown,
	}
}

// Ready 执行运行时就绪检查。
func (a *AppRuntime) Ready() error {
	if a == nil || a.ready == nil {
		return nil
	}
	return a.ready()
}

// Shutdown 执行运行时关闭逻辑。
func (a *AppRuntime) Shutdown() error {
	if a == nil || a.shutdown == nil {
		return nil
	}
	return a.shutdown()
}
