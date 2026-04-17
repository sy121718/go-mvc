package config

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"go-mvc/internal/task"
	"go-mvc/pkg/auth"
	"go-mvc/pkg/cache"
	"go-mvc/pkg/casbin"
	"go-mvc/pkg/database"
	"go-mvc/pkg/i18n"
	pkglogger "go-mvc/pkg/logger"
	"go-mvc/pkg/upload"

	"github.com/spf13/viper"
)

var (
	v  *viper.Viper
	mu sync.Mutex
)

// ServerConfig 服务配置。
type ServerConfig struct {
	Port    int    `mapstructure:"port"`
	Mode    string `mapstructure:"mode"`
	AppName string `mapstructure:"app_name"`
}

// Init 初始化配置文件。
func Init(configPath string) error {
	mu.Lock()
	defer mu.Unlock()

	if v != nil {
		return nil
	}

	cfg := viper.New()
	setDefaults(cfg)
	cfg.SetConfigFile(configPath)

	if err := cfg.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	v = cfg
	log.Printf("配置加载成功: %s", configPath)
	return nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "debug")
	v.SetDefault("server.app_name", "go-mvc")

	v.SetDefault("database.driver", "mysql")
	v.SetDefault("database.host", "127.0.0.1")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.user", "root")
	v.SetDefault("database.password", "")
	v.SetDefault("database.dbname", "test")
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.max_open_conns", 100)
	v.SetDefault("database.log_level", "")

	v.SetDefault("redis.host", "127.0.0.1")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.enabled", true)
	v.SetDefault("redis.provider", "redis")
	v.SetDefault("redis.addrs", []string{})

	v.SetDefault("jwt.secret", "default-secret-key-please-change-in-production")
	v.SetDefault("jwt.expire_time", 24)
	v.SetDefault("jwt.issuer", "go-mvc")

	v.SetDefault("casbin.enabled", true)
	v.SetDefault("i18n.default_lang", "zh-CN")
	v.SetDefault("i18n.auto_refresh", true)
	v.SetDefault("i18n.refresh_interval", "20s")
	v.SetDefault("queue.enabled", false)
	v.SetDefault("queue.provider", "asynq")
	v.SetDefault("queue.run_worker", false)
	v.SetDefault("queue.concurrency", 10)
	v.SetDefault("queue.redis.host", "")
	v.SetDefault("queue.redis.port", 0)
	v.SetDefault("queue.redis.password", "")
	v.SetDefault("queue.redis.db", 0)

	v.SetDefault("upload.enabled", true)
	v.SetDefault("upload.default_provider", "local")

	v.SetDefault("log.level", "info")
	v.SetDefault("log.filename", "public/logs/app.log")
	v.SetDefault("log.base_dir", "public/logs")
	v.SetDefault("log.max_size", 100)
	v.SetDefault("log.max_backups", 10)
	v.SetDefault("log.max_age", 30)
	v.SetDefault("log.compress", false)
}

// GetViper 获取 viper 实例（供 pkg 使用）。
func GetViper() *viper.Viper {
	if v == nil {
		panic("配置未初始化，请先调用 config.Init()")
	}
	return v
}

// GetServer 获取服务配置。
func GetServer() (ServerConfig, error) {
	var cfg ServerConfig
	if err := GetViper().UnmarshalKey("server", &cfg); err != nil {
		return ServerConfig{}, fmt.Errorf("解析 Server 配置失败: %w", err)
	}
	return cfg, nil
}

// InitComponents 初始化所有组件。
func InitComponents() error {
	cfg := GetViper()
	log.Println("开始初始化组件...")

	if err := pkglogger.Init(cfg); err != nil {
		return fmt.Errorf("初始化日志组件失败: %w", err)
	}

	if err := database.InitDB(cfg); err != nil {
		return err
	}

	log.Println("初始化多语言配置中心...")
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

	if cfg.GetBool("casbin.enabled") {
		log.Println("初始化 Casbin...")
		db, err := database.GetDB()
		if err != nil {
			return fmt.Errorf("获取数据库实例失败: %w", err)
		}
		if err := casbin.InitCasbin(db); err != nil {
			return fmt.Errorf("初始化 Casbin 失败: %w", err)
		}
	}

	if cfg.GetBool("redis.enabled") {
		if err := cache.InitRedis(cfg); err != nil {
			return err
		}
	}

	if err := auth.InitJWT(cfg); err != nil {
		return err
	}

	if cfg.GetBool("upload.enabled") {
		if err := upload.Init(cfg); err != nil {
			return fmt.Errorf("初始化上传组件失败: %w", err)
		}
	}

	if cfg.GetBool("queue.enabled") {
		if err := task.Init(cfg); err != nil {
			return fmt.Errorf("初始化任务队列失败: %w", err)
		}
		if cfg.GetBool("queue.run_worker") {
			if err := task.StartQueue(); err != nil {
				return fmt.Errorf("启动任务队列失败: %w", err)
			}
		}
	}

	log.Println("组件初始化完成")
	return nil
}

// CloseComponents 关闭所有组件。
func CloseComponents() error {
	log.Println("开始关闭组件...")

	var closeErr error

	i18n.StopAutoRefresh()

	if err := task.ShutdownQueue(); err != nil {
		closeErr = errors.Join(closeErr, err)
	}

	if cache.IsInited() {
		if err := cache.Close(); err != nil {
			closeErr = errors.Join(closeErr, err)
		}
	}

	if upload.IsInited() {
		if err := upload.Close(); err != nil {
			closeErr = errors.Join(closeErr, err)
		}
	}

	if database.IsInited() {
		if err := database.Close(); err != nil {
			closeErr = errors.Join(closeErr, err)
		}
	}

	if err := pkglogger.Sync(); err != nil {
		closeErr = errors.Join(closeErr, err)
	}

	if closeErr != nil {
		return closeErr
	}

	log.Println("组件关闭完成")
	return nil
}
