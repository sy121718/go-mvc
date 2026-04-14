package config

import (
	"fmt"
	"go-mvc/pkg/auth"
	"go-mvc/pkg/cache"
	"go-mvc/pkg/casbin"
	"go-mvc/pkg/database"
	"go-mvc/pkg/i18n"
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
	v      *viper.Viper
	mu     sync.Mutex
	inited bool
)

// ServerConfig 服务配置（核心配置，启动时加载）
type ServerConfig struct {
	Port    int    `mapstructure:"port"`
	Mode    string `mapstructure:"mode"`
	AppName string `mapstructure:"app_name"`
}

// Init 初始化配置文件
func Init(configPath string) error {
	mu.Lock()
	defer mu.Unlock()

	if inited {
		return nil
	}

	v = viper.New()
	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	inited = true
	log.Printf("配置加载成功: %s", configPath)
	return nil
}

// GetViper 获取 viper 实例（供 pkg 使用）
func GetViper() *viper.Viper {
	if v == nil {
		panic("配置未初始化，请先调用 config.Init()")
	}
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

// InitComponents 初始化所有组件
// 这是配置启动器，负责驱动 pkg 组件初始化
// pkg 内部会处理错误，致命错误会直接退出程序
func InitComponents(v *viper.Viper) {
	log.Println("开始初始化组件...")

	// 初始化数据库（内部处理错误）
	database.InitDB(v)

	// 数据库初始化后，初始化多语言配置中心（依赖 DB 连接）
	if database.IsInited() {
		log.Println("初始化多语言配置中心...")
		i18n.Init()
		// 启动自动刷新（每10秒）
		i18n.StartAutoRefresh()
	}

	// 数据库初始化后，初始化 Casbin（依赖 DB 连接）
	if database.IsInited() {
		log.Println("初始化 Casbin...")
		casbin.InitCasbin(database.GetDB())
	}

	// 初始化 Redis（内部处理错误）
	cache.InitRedis(v)

	// 初始化 JWT（内部处理错误）
	auth.InitJWT(v)

	log.Println("组件初始化完成")
}

// CloseComponents 关闭所有组件
func CloseComponents() {
	log.Println("开始关闭组件...")

	// 只关闭已初始化的组件
	if database.IsInited() {
		if err := database.Close(); err != nil {
			log.Printf("关闭数据库失败: %v", err)
		}
	}

	if err := cache.Close(); err != nil {
		log.Printf("关闭 Redis 失败: %v", err)
	}

	log.Println("组件关闭完成")
}
