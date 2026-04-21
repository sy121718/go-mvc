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

// runtimeComponent 描述一个由框架运行时统一编排的启动型组件。
//
// 注册准入原则：
// - 只有“系统启动必须感知”的组件才能放进这里
// - 典型对象是 database、logger、auth、i18n、cache、queue、upload、casbin
// - 纯工具包、DTO、model、helper、业务 service 不允许放进这里
// - 组件自身负责 Init/Ready/Close 的实现细节，这里只维护清单和顺序
type runtimeComponent struct {
	Name     string
	Critical bool
	Enabled  func(cfg *viper.Viper) bool
	Init     func(cfg *viper.Viper) error
	Ready    func() error
	Close    func() error
}

// runtimeModule 描述一个业务模块入口。
//
// 说明：
// - config 层只依赖模块入口，不直接依赖模块内部的 handle/service/model 细节
// - 模块内部如何拆分 router/handle/service，由模块自己决定
type runtimeModule struct {
	Register func(rg *gin.RouterGroup)
}

// runtimePreparers 是运行时编排前的显式预处理动作。
//
// 说明：
// - 这里只允许放“启动前必须显式执行一次”的预处理逻辑
// - 当前仅用于注册项目私有任务处理器
// - 不允许在这里做组件初始化，不允许替代 runtimeComponents
var runtimePreparers = []func(){
	internaltask.RegisterHandlers,
}

// runtimeComponents 是唯一的启动型组件注册清单。
//
// 使用规则：
// - Critical=true 表示关键启动阶段，优先初始化
// - Critical=false 表示扩展启动阶段，在关键组件之后初始化
// - Enabled 用于按配置启停组件
// - Init/Ready/Close 一律直接绑定到 pkg 组件入口，不在这里写业务胶水
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

// runtimeModules 是业务模块入口清单。
//
// 规则：
// - 这里只放模块入口函数
// - 不直接引用模块内部 router 子包以外的更细层级
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
