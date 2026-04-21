package config

import (
	"fmt"

	adminrouter "go-mvc/internal/module/backend/admin/router"
	_ "go-mvc/internal/task"
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
	Name     string
	Register func(rg *gin.RouterGroup)
}

func registeredComponents() []runtimeComponent {
	return []runtimeComponent{
		{
			Name: "logger",
			Init: func(cfg *viper.Viper) error {
				if err := pkglogger.Init(cfg); err != nil {
					return fmt.Errorf("初始化日志组件失败: %w", err)
				}
				return nil
			},
			Close: pkglogger.Sync,
		},
		{
			Name:  "database",
			Init:  database.InitDB,
			Close: database.Close,
		},
		{
			Name: "i18n",
			Init: func(cfg *viper.Viper) error {
				if err := i18n.Init(cfg); err != nil {
					return fmt.Errorf("初始化多语言配置中心失败: %w", err)
				}
				return nil
			},
			Close: i18n.Close,
		},
		{
			Name: "casbin",
			Enabled: func(cfg *viper.Viper) bool {
				return cfg.GetBool("casbin.enabled")
			},
			Init: func(cfg *viper.Viper) error {
				if err := casbin.Init(cfg); err != nil {
					return fmt.Errorf("初始化 Casbin 失败: %w", err)
				}
				return nil
			},
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
			Init: func(cfg *viper.Viper) error {
				if err := upload.Init(cfg); err != nil {
					return fmt.Errorf("初始化上传组件失败: %w", err)
				}
				return nil
			},
			Close: upload.Close,
		},
		{
			Name: "queue",
			Enabled: func(cfg *viper.Viper) bool {
				return cfg.GetBool("queue.enabled")
			},
			Init: func(cfg *viper.Viper) error {
				if err := queue.Init(cfg); err != nil {
					return fmt.Errorf("初始化任务队列失败: %w", err)
				}
				return nil
			},
			Close: queue.Close,
		},
	}
}

func registeredModules() []runtimeModule {
	return []runtimeModule{
		{
			Name:     "backend.admin",
			Register: adminrouter.SetupAdminRoutes,
		},
	}
}
