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
}

// NewAppRuntime 创建应用运行时容器。
func NewAppRuntime(server ServerConfig, router *gin.Engine, httpServer *http.Server) *AppRuntime {
	return &AppRuntime{
		Server:     server,
		Router:     router,
		HTTPServer: httpServer,
	}
}
