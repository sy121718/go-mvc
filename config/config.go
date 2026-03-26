package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"sync"
)

/*
配置管理
===========================================
职责：
- 读取 config.yaml 文件
- 提供原始配置访问（通过 viper）
- 不定义具体业务配置结构体

配置结构体由各个 pkg 自己定义
*/

var (
	v    *viper.Viper
	once sync.Once
)

// ServerConfig 服务配置（核心配置，启动时加载）
type ServerConfig struct {
	Port    int    `mapstructure:"port"`
	Mode    string `mapstructure:"mode"`
	AppName string `mapstructure:"app_name"`
}

// Init 初始化配置文件
func Init(configPath string) error {
	var err error
	once.Do(func() {
		v = viper.New()
		v.SetConfigFile(configPath)

		// 设置默认值
		setDefaults()

		// 读取配置文件
		if err = v.ReadInConfig(); err != nil {
			err = fmt.Errorf("读取配置文件失败: %v", err)
			return
		}

		log.Printf("配置加载成功: %s", configPath)
	})
	return err
}

// GetViper 获取 viper 实例（供 pkg 使用）
func GetViper() *viper.Viper {
	return v
}

// GetServer 获取服务配置
func GetServer() ServerConfig {
	var cfg ServerConfig
	if err := v.UnmarshalKey("server", &cfg); err != nil {
		log.Fatalf("解析 Server 配置失败: %v", err)
	}
	return cfg
}

// setDefaults 设置默认值
func setDefaults() {
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
}
