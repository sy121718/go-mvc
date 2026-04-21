package config

import (
	adminmodule "go-mvc/internal/module/backend/admin"
	internaltask "go-mvc/internal/task"
	"go-mvc/pkg/auth"
	"go-mvc/pkg/cache"
	"go-mvc/pkg/casbin"
	"go-mvc/pkg/database"
	"go-mvc/pkg/i18n"
	pkglogger "go-mvc/pkg/logger"
	"go-mvc/pkg/queue"
	"go-mvc/pkg/upload"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type runtimeComponent struct {
	Name     string
	Critical bool
	Enabled  func(cfg *viper.Viper) bool
	Init     func(cfg *viper.Viper) error
	Ready    func() error
	Close    func() error
}

type runtimeModule struct {
	Register func(rg *gin.RouterGroup)
}

var runtimePreparers = []func(){
	internaltask.RegisterHandlers,
}

var runtimeComponents = []runtimeComponent{
	{
		Name:     "logger",
		Critical: true,
		Init:     pkglogger.Init,
		Ready:    pkglogger.Ready,
		Close:    pkglogger.Close,
	},
	{
		Name:     "database",
		Critical: true,
		Init:     database.Init,
		Ready:    database.Ready,
		Close:    database.Close,
	},
	{
		Name:     "i18n",
		Critical: true,
		Init:     i18n.Init,
		Ready:    i18n.Ready,
		Close:    i18n.Close,
	},
	{
		Name: "casbin",
		Enabled: func(cfg *viper.Viper) bool {
			return cfg.GetBool("casbin.enabled")
		},
		Init:  casbin.Init,
		Ready: casbin.Ready,
		Close: casbin.Close,
	},
	{
		Name: "cache",
		Enabled: func(cfg *viper.Viper) bool {
			return cfg.GetBool("redis.enabled")
		},
		Init:  cache.Init,
		Ready: cache.Ready,
		Close: cache.Close,
	},
	{
		Name:     "auth",
		Critical: true,
		Init:     auth.Init,
		Ready:    auth.Ready,
		Close:    auth.Close,
	},
	{
		Name: "upload",
		Enabled: func(cfg *viper.Viper) bool {
			return cfg.GetBool("upload.enabled")
		},
		Init:  upload.Init,
		Ready: upload.Ready,
		Close: upload.Close,
	},
	{
		Name: "queue",
		Enabled: func(cfg *viper.Viper) bool {
			return cfg.GetBool("queue.enabled")
		},
		Init:  queue.Init,
		Ready: queue.Ready,
		Close: queue.Close,
	},
}

var runtimeModules = []runtimeModule{
	{
		Register: adminmodule.RegisterRoutes,
	},
}

// ModuleRegistrars 返回当前已注册的业务模块路由注册函数列表。
func ModuleRegistrars() []func(*gin.RouterGroup) {
	registrars := make([]func(*gin.RouterGroup), 0, len(runtimeModules))
	for _, module := range runtimeModules {
		registrars = append(registrars, module.Register)
	}
	return registrars
}
