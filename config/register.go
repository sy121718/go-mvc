package config

import (
	adminrouter "go-mvc/internal/module/backend/admin/router"
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
	Name    string
	Enabled func(cfg *viper.Viper) bool
	Init    func(cfg *viper.Viper) error
	Close   func() error
}

type runtimeModule struct {
	Register func(rg *gin.RouterGroup)
}

var runtimePreparers = []func(){
	internaltask.RegisterHandlers,
}

var runtimeComponents = []runtimeComponent{
	{
		Name:  "logger",
		Init:  pkglogger.Init,
		Close: pkglogger.Sync,
	},
	{
		Name:  "database",
		Init:  database.InitDB,
		Close: database.Close,
	},
	{
		Name:  "i18n",
		Init:  i18n.Init,
		Close: i18n.Close,
	},
	{
		Name: "casbin",
		Enabled: func(cfg *viper.Viper) bool {
			return cfg.GetBool("casbin.enabled")
		},
		Init:  casbin.Init,
		Close: casbin.Close,
	},
	{
		Name: "cache",
		Enabled: func(cfg *viper.Viper) bool {
			return cfg.GetBool("redis.enabled")
		},
		Init:  cache.InitRedis,
		Close: cache.Close,
	},
	{
		Name: "auth",
		Init: auth.InitJWT,
	},
	{
		Name: "upload",
		Enabled: func(cfg *viper.Viper) bool {
			return cfg.GetBool("upload.enabled")
		},
		Init:  upload.Init,
		Close: upload.Close,
	},
	{
		Name: "queue",
		Enabled: func(cfg *viper.Viper) bool {
			return cfg.GetBool("queue.enabled")
		},
		Init:  queue.Init,
		Close: queue.Close,
	},
}

var runtimeModules = []runtimeModule{
	{
		Register: adminrouter.SetupAdminRoutes,
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
