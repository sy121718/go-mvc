package config

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"
)

var (
	v  *viper.Viper
	mu sync.Mutex
)

// ServerConfig 服务配置。
type ServerConfig struct {
	Port              int
	Mode              string
	AppName           string
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
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
	v.SetDefault("server.read_header_timeout", "3s")
	v.SetDefault("server.read_timeout", "15s")
	v.SetDefault("server.write_timeout", "30s")
	v.SetDefault("server.idle_timeout", "60s")

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
	type serverConfigRaw struct {
		Port              int    `mapstructure:"port"`
		Mode              string `mapstructure:"mode"`
		AppName           string `mapstructure:"app_name"`
		ReadHeaderTimeout string `mapstructure:"read_header_timeout"`
		ReadTimeout       string `mapstructure:"read_timeout"`
		WriteTimeout      string `mapstructure:"write_timeout"`
		IdleTimeout       string `mapstructure:"idle_timeout"`
	}

	var raw serverConfigRaw
	if err := GetViper().UnmarshalKey("server", &raw); err != nil {
		return ServerConfig{}, fmt.Errorf("解析 Server 配置失败: %w", err)
	}

	readHeaderTimeout, err := parseServerDuration("read_header_timeout", raw.ReadHeaderTimeout)
	if err != nil {
		return ServerConfig{}, err
	}
	readTimeout, err := parseServerDuration("read_timeout", raw.ReadTimeout)
	if err != nil {
		return ServerConfig{}, err
	}
	writeTimeout, err := parseServerDuration("write_timeout", raw.WriteTimeout)
	if err != nil {
		return ServerConfig{}, err
	}
	idleTimeout, err := parseServerDuration("idle_timeout", raw.IdleTimeout)
	if err != nil {
		return ServerConfig{}, err
	}

	return ServerConfig{
		Port:              raw.Port,
		Mode:              raw.Mode,
		AppName:           raw.AppName,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}, nil
}

func parseServerDuration(field string, raw string) (time.Duration, error) {
	duration, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("解析 server.%s 失败: %w", field, err)
	}
	return duration, nil
}
