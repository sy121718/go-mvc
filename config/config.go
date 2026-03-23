package config

import (
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

var (
	cfg     *Config
	once    sync.Once
)

// Config 总配置结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
}

// ServerConfig 服务配置
type ServerConfig struct {
	Port    int    `mapstructure:"port"`
	Mode    string `mapstructure:"mode"` // debug, release, test
	AppName string `mapstructure:"app_name"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	LazyInit     bool   `mapstructure:"lazy_init"` // 是否懒加载
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	LazyInit bool   `mapstructure:"lazy_init"` // 是否懒加载
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	ExpireTime int    `mapstructure:"expire_time"` // 过期时间（小时）
	Issuer     string `mapstructure:"issuer"`
	LazyInit   bool   `mapstructure:"lazy_init"` // 是否懒加载
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`       // debug, info, warn, error
	Filename   string `mapstructure:"filename"`    // 日志文件路径
	MaxSize    int    `mapstructure:"max_size"`    // 日志文件最大大小（MB）
	MaxBackups int    `mapstructure:"max_backups"` // 最多保留旧文件数
	MaxAge     int    `mapstructure:"max_age"`     // 最多保留天数
	Compress   bool   `mapstructure:"compress"`    // 是否压缩
}

// Init 初始化配置
func Init(configPath string) error {
	var err error
	once.Do(func() {
		err = initConfig(configPath)
	})
	return err
}

func initConfig(configPath string) error {
	v := viper.New()

	// 设置默认值
	setDefaults(v)

	// 设置配置文件
	v.SetConfigFile(configPath)

	// 尝试读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析到结构体
	cfg = &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return fmt.Errorf("解析配置失败: %v", err)
	}

	return nil
}

// setDefaults 设置默认值
func setDefaults(v *viper.Viper) {
	// 服务默认配置
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "debug")
	v.SetDefault("server.app_name", "go-mvc")

	// 数据库默认配置
	v.SetDefault("database.host", "127.0.0.1")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.user", "root")
	v.SetDefault("database.password", "")
	v.SetDefault("database.dbname", "go_mvc")
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.max_open_conns", 100)
	v.SetDefault("database.lazy_init", false)

	// Redis 默认配置
	v.SetDefault("redis.host", "127.0.0.1")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.lazy_init", false)

	// JWT 默认配置
	v.SetDefault("jwt.secret", "your-secret-key")
	v.SetDefault("jwt.expire_time", 24)
	v.SetDefault("jwt.issuer", "go-mvc")
	v.SetDefault("jwt.lazy_init", false)

	// 日志默认配置
	v.SetDefault("log.level", "info")
	v.SetDefault("log.filename", "logs/app.log")
	v.SetDefault("log.max_size", 100)
	v.SetDefault("log.max_backups", 10)
	v.SetDefault("log.max_age", 30)
	v.SetDefault("log.compress", false)
}

// Get 获取完整配置
func Get() *Config {
	return cfg
}

// GetServer 获取服务配置
func GetServer() ServerConfig {
	return cfg.Server
}

// GetDatabase 获取数据库配置
func GetDatabase() DatabaseConfig {
	return cfg.Database
}

// GetRedis 获取Redis配置
func GetRedis() RedisConfig {
	return cfg.Redis
}

// GetJWT 获取JWT配置
func GetJWT() JWTConfig {
	return cfg.JWT
}

// GetLog 获取日志配置
func GetLog() LogConfig {
	return cfg.Log
}