/*
Redis 缓存组件包
===========================================
提供 Redis 连接管理功能

主要功能：
- Redis 连接初始化
- 懒加载支持（按需连接）
- 全局单例模式
- 缓存操作封装

配置说明（config.yaml）：
  redis:
    host: 127.0.0.1      # Redis 地址
    port: 6379           # 端口
    password: ""         # 密码
    db: 0                # 数据库编号
    lazy_init: false     # 是否懒加载

使用示例：
  // 在 main.go 中初始化
  cache.InitRedis(viper)

  // 在业务代码中使用
  rdb := cache.GetRedis(viper)
  rdb.Set(ctx, "key", "value", time.Hour)
  rdb.Get(ctx, "key")

PHP 对比：
  // Laravel Redis
  Redis::set('key', 'value');
  Redis::get('key');

  // Go
  rdb.Set(ctx, "key", "value", 0)
  rdb.Get(ctx, "key")
*/
package cache

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

/*
Redis 组件
===========================================
配置结构体定义在这里，自己解析配置
*/

// Config Redis配置
type Config struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	LazyInit bool   `mapstructure:"lazy_init"`
}

var (
	rdb  *redis.Client
	once sync.Once
)

// GetRedis 获取 Redis 客户端实例（懒加载）
func GetRedis(v *viper.Viper) *redis.Client {
	once.Do(func() {
		if err := initRedis(v); err != nil {
			panic(fmt.Sprintf("Redis 初始化失败: %v", err))
		}
	})
	return rdb
}

// initRedis 初始化 Redis
func initRedis(v *viper.Viper) error {
	// 自己解析配置
	var cfg Config
	if err := v.UnmarshalKey("redis", &cfg); err != nil {
		return fmt.Errorf("解析 Redis 配置失败: %v", err)
	}

	// 检查是否懒加载
	if cfg.LazyInit {
		return nil
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis 连接失败: %v", err)
	}

	log.Println("Redis 初始化成功")
	return nil
}

// InitRedis 手动初始化 Redis（用于非懒加载场景）
func InitRedis(v *viper.Viper) error {
	return initRedis(v)
}
