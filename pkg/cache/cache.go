package cache

import (
	"fmt"
	"strings"

	cacheprovider "go-mvc/pkg/cache/provider"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func buildProvider(v *viper.Viper) (cacheprovider.Provider, error) {
	providerName := "redis"
	if v != nil {
		providerName = strings.TrimSpace(strings.ToLower(v.GetString("redis.provider")))
		if providerName == "" {
			providerName = "redis"
		}
	}

	switch providerName {
	case "redis":
		return cacheprovider.NewRedisProvider(), nil
	default:
		return nil, fmt.Errorf("不支持的缓存 provider: %s", providerName)
	}
}

var defaultProvider cacheprovider.Provider = cacheprovider.NewRedisProvider()

// InitRedis 初始化 Redis 缓存。
func InitRedis(v *viper.Viper) error {
	provider, err := buildProvider(v)
	if err != nil {
		return err
	}
	defaultProvider = provider
	return defaultProvider.Init(v)
}

// GetRedis 获取 Redis 客户端实例。
func GetRedis() (redis.UniversalClient, error) {
	return defaultProvider.Client()
}

// IsInited 检查 Redis 是否已初始化。
func IsInited() bool {
	return defaultProvider.IsInited()
}

// Close 关闭 Redis 连接。
func Close() error {
	return defaultProvider.Close()
}
