package cache

import (
	"context"
	"fmt"
	"go-mvc/config"
	"github.com/redis/go-redis/v9"
	"sync"
)

var (
	rdb  *redis.Client
	once sync.Once
)

// GetRedis 获取 Redis 客户端实例（懒加载）
func GetRedis() *redis.Client {
	once.Do(func() {
		if err := initRedis(); err != nil {
			panic(fmt.Sprintf("Redis 初始化失败: %v", err))
		}
	})
	return rdb
}

// initRedis 初始化 Redis
func initRedis() error {
	cfg := config.GetRedis()

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

	return nil
}

// InitRedis 手动初始化 Redis（用于懒加载场景）
func InitRedis() error {
	return initRedis()
}