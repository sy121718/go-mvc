package config

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-mvc/pkg/auth"
	"go-mvc/pkg/defaults"

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
	RequestBodyLimit  int64
	UploadBodyLimit   int64
	RateLimitEnabled  bool
	RateLimitLimit    int
	RateLimitWindow   time.Duration
	PortStrategy      string
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
	v.SetDefault("server.request_body_limit", "2MB")
	v.SetDefault("server.upload_body_limit", "32MB")
	v.SetDefault("server.rate_limit_enabled", true)
	v.SetDefault("server.rate_limit_limit", 120)
	v.SetDefault("server.rate_limit_window", "1m")
	v.SetDefault("server.port_strategy", "")

	v.SetDefault("database.driver", defaults.DefaultDatabaseDriver)
	v.SetDefault("database.host", defaults.DefaultDatabaseHost)
	v.SetDefault("database.port", defaults.DefaultDatabasePort)
	v.SetDefault("database.user", defaults.DefaultDatabaseUser)
	v.SetDefault("database.password", defaults.DefaultDatabasePassword)
	v.SetDefault("database.dbname", defaults.DefaultDatabaseName)
	v.SetDefault("database.max_idle_conns", defaults.DefaultDatabaseMaxIdleConns)
	v.SetDefault("database.max_open_conns", defaults.DefaultDatabaseMaxOpenConns)
	v.SetDefault("database.log_level", "")
	v.SetDefault("database.prepare_stmt", false)
	v.SetDefault("database.skip_default_transaction", false)
	v.SetDefault("database.slow_threshold", "200ms")

	v.SetDefault("redis.host", defaults.DefaultRedisHost)
	v.SetDefault("redis.port", defaults.DefaultRedisPort)
	v.SetDefault("redis.password", defaults.DefaultRedisPassword)
	v.SetDefault("redis.db", defaults.DefaultRedisDB)
	v.SetDefault("redis.enabled", true)
	v.SetDefault("redis.provider", defaults.DefaultRedisProvider)
	v.SetDefault("redis.addrs", []string{})

	v.SetDefault("jwt.secret", defaults.DefaultJWTSecret)
	v.SetDefault("jwt.expire_time", defaults.DefaultJWTExpireTime)
	v.SetDefault("jwt.issuer", defaults.DefaultJWTIssuer)

	v.SetDefault("casbin.enabled", true)
	v.SetDefault("i18n.default_lang", "zh-CN")
	v.SetDefault("i18n.auto_refresh", true)
	v.SetDefault("i18n.refresh_interval", "20s")
	v.SetDefault("queue.enabled", false)
	v.SetDefault("queue.provider", defaults.DefaultQueueProvider)
	v.SetDefault("queue.run_worker", false)
	v.SetDefault("queue.concurrency", defaults.DefaultQueueConcurrency)
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
	v.SetDefault("log.sample.enabled", false)
	v.SetDefault("log.sample.initial", 100)
	v.SetDefault("log.sample.thereafter", 100)
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
		RequestBodyLimit  string `mapstructure:"request_body_limit"`
		UploadBodyLimit   string `mapstructure:"upload_body_limit"`
		RateLimitEnabled  bool   `mapstructure:"rate_limit_enabled"`
		RateLimitLimit    int    `mapstructure:"rate_limit_limit"`
		RateLimitWindow   string `mapstructure:"rate_limit_window"`
		PortStrategy      string `mapstructure:"port_strategy"`
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
	requestBodyLimit, err := parseByteSize("request_body_limit", raw.RequestBodyLimit)
	if err != nil {
		return ServerConfig{}, err
	}
	uploadBodyLimit, err := parseByteSize("upload_body_limit", raw.UploadBodyLimit)
	if err != nil {
		return ServerConfig{}, err
	}
	rateLimitWindow, err := parseServerDuration("rate_limit_window", raw.RateLimitWindow)
	if err != nil {
		return ServerConfig{}, err
	}
	if raw.RateLimitLimit <= 0 {
		return ServerConfig{}, fmt.Errorf("解析 server.rate_limit_limit 失败: 值必须大于 0")
	}

	return ServerConfig{
		Port:              raw.Port,
		Mode:              raw.Mode,
		AppName:           raw.AppName,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		RequestBodyLimit:  requestBodyLimit,
		UploadBodyLimit:   uploadBodyLimit,
		RateLimitEnabled:  raw.RateLimitEnabled,
		RateLimitLimit:    raw.RateLimitLimit,
		RateLimitWindow:   rateLimitWindow,
		PortStrategy:      raw.PortStrategy,
	}, nil
}

func parseServerDuration(field string, raw string) (time.Duration, error) {
	duration, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("解析 server.%s 失败: %w", field, err)
	}
	return duration, nil
}

func parseByteSize(field string, raw string) (int64, error) {
	if raw == "" {
		return 0, fmt.Errorf("解析 server.%s 失败: 值不能为空", field)
	}

	normalized := strings.ToUpper(strings.TrimSpace(raw))
	units := []struct {
		Suffix string
		Scale  int64
	}{
		{Suffix: "KB", Scale: 1024},
		{Suffix: "MB", Scale: 1024 * 1024},
		{Suffix: "GB", Scale: 1024 * 1024 * 1024},
		{Suffix: "B", Scale: 1},
	}

	for _, unit := range units {
		if strings.HasSuffix(normalized, unit.Suffix) {
			number := strings.TrimSpace(strings.TrimSuffix(normalized, unit.Suffix))
			value, err := strconv.ParseInt(number, 10, 64)
			if err != nil {
				return 0, fmt.Errorf("解析 server.%s 失败: %w", field, err)
			}
			if value <= 0 {
				return 0, fmt.Errorf("解析 server.%s 失败: 值必须大于 0", field)
			}
			return value * unit.Scale, nil
		}
	}

	value, err := strconv.ParseInt(normalized, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("解析 server.%s 失败: %w", field, err)
	}
	if value <= 0 {
		return 0, fmt.Errorf("解析 server.%s 失败: 值必须大于 0", field)
	}
	return value, nil
}

// ValidateRuntimeConfig 在组件初始化前执行关键配置校验。
//
// 当前规则：
// - 仅 release 模式启用严格 fail-fast
// - 拒绝默认 JWT secret
// - 拒绝默认数据库名
// - 拒绝空数据库密码
func ValidateRuntimeConfig() error {
	cfg := GetViper()
	if cfg.GetString("server.mode") != "release" {
		return nil
	}

	if err := auth.ValidateConfig(cfg, true); err != nil {
		return err
	}

	dbName := cfg.GetString("database.dbname")
	if dbName == "" {
		return fmt.Errorf("database.dbname 不能为空")
	}
	if dbName == "test" {
		return fmt.Errorf("database.dbname 不能使用默认值")
	}

	if cfg.GetString("database.password") == "" {
		return fmt.Errorf("database.password 不能为空")
	}

	return nil
}

// ResetForTest 重置 config 包的全局配置状态，仅用于测试环境。
func ResetForTest() {
	mu.Lock()
	defer mu.Unlock()
	v = nil
}
