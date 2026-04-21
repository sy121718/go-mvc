package config

import (
	"fmt"
	"time"

	adminrouter "go-mvc/internal/module/backend/admin/router"
	"go-mvc/internal/task"
	"go-mvc/pkg/auth"
	"go-mvc/pkg/cache"
	"go-mvc/pkg/casbin"
	"go-mvc/pkg/database"
	"go-mvc/pkg/i18n"
	pkglogger "go-mvc/pkg/logger"
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
				i18n.SetDefaultLang(cfg.GetString("i18n.default_lang"))
				if err := i18n.Init(); err != nil {
					return fmt.Errorf("初始化多语言配置中心失败: %w", err)
				}

				if cfg.GetBool("i18n.auto_refresh") {
					refreshInterval, err := time.ParseDuration(cfg.GetString("i18n.refresh_interval"))
					if err != nil {
						return fmt.Errorf("解析 i18n.refresh_interval 失败: %w", err)
					}
					i18n.StartAutoRefresh(refreshInterval)
				}
				return nil
			},
			Close: func() error {
				i18n.StopAutoRefresh()
				return nil
			},
		},
		{
			Name: "casbin",
			Enabled: func(cfg *viper.Viper) bool {
				return cfg.GetBool("casbin.enabled")
			},
			Init: func(cfg *viper.Viper) error {
				db, err := database.GetDB()
				if err != nil {
					return fmt.Errorf("获取数据库实例失败: %w", err)
				}
				if err := casbin.InitCasbin(db); err != nil {
					return fmt.Errorf("初始化 Casbin 失败: %w", err)
				}
				return nil
			},
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
				if err := task.Init(cfg); err != nil {
					return fmt.Errorf("初始化任务队列失败: %w", err)
				}
				if cfg.GetBool("queue.run_worker") {
					if err := task.StartQueue(); err != nil {
						return fmt.Errorf("启动任务队列失败: %w", err)
					}
				}
				return nil
			},
			Close: task.ShutdownQueue,
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
