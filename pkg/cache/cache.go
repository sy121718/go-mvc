package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	cacheprovider "go-mvc/pkg/cache/provider"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"golang.org/x/sync/singleflight"
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
var loadGroup singleflight.Group

// Init 初始化 Redis 缓存组件。
func Init(v *viper.Viper) error {
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

// SetJSON 将任意结构体序列化为 JSON 后写入缓存。
func SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error {
	client, err := GetRedis()
	if err != nil {
		return err
	}

	payload, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化缓存 JSON 失败: %w", err)
	}

	if err := client.Set(ctx, key, payload, ttl).Err(); err != nil {
		return fmt.Errorf("写入缓存失败: %w", err)
	}
	return nil
}

// GetJSON 从缓存中读取 JSON 并反序列化为指定类型。
func GetJSON[T any](ctx context.Context, key string) (T, error) {
	var zero T

	client, err := GetRedis()
	if err != nil {
		return zero, err
	}

	payload, err := client.Get(ctx, key).Bytes()
	if err != nil {
		return zero, err
	}

	var result T
	if err := json.Unmarshal(payload, &result); err != nil {
		return zero, fmt.Errorf("反序列化缓存 JSON 失败: %w", err)
	}
	return result, nil
}

// RememberJSON 优先从缓存读取 JSON，缓存缺失时通过 loader 回源加载，并使用 singleflight 合并并发请求。
//
// jitter 表示在 ttl 基础上的随机抖动上限，用于降低同一批 key 同时过期带来的压力。
func RememberJSON[T any](
	ctx context.Context,
	key string,
	ttl time.Duration,
	jitter time.Duration,
	loader func(context.Context) (T, error),
) (T, error) {
	var zero T

	result, err := GetJSON[T](ctx, key)
	if err == nil {
		return result, nil
	}
	if err != redis.Nil {
		return zero, err
	}

	value, loadErr, _ := loadGroup.Do(key, func() (any, error) {
		loaded, err := loader(ctx)
		if err != nil {
			return zero, err
		}
		cacheTTL := addTTLJitter(ttl, jitter)
		if err := SetJSON(ctx, key, loaded, cacheTTL); err != nil {
			return zero, err
		}
		return loaded, nil
	})
	if loadErr != nil {
		return zero, loadErr
	}

	typed, ok := value.(T)
	if !ok {
		return zero, fmt.Errorf("缓存加载结果类型断言失败")
	}
	return typed, nil
}

// IsInited 检查 Redis 是否已初始化。
func IsInited() bool {
	return defaultProvider.IsInited()
}

// Ready 检查缓存组件是否已初始化。
func Ready() error {
	if !IsInited() {
		return fmt.Errorf("缓存组件未初始化")
	}
	return nil
}

// Close 关闭 Redis 连接。
func Close() error {
	return defaultProvider.Close()
}

func addTTLJitter(ttl time.Duration, jitter time.Duration) time.Duration {
	if ttl <= 0 || jitter <= 0 {
		return ttl
	}
	extra := time.Duration(rand.Int63n(int64(jitter)))
	return ttl + extra
}
