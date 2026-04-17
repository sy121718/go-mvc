package cache_test

import (
	"context"
	"testing"
	"time"

	"go-mvc/pkg/cache"

	"github.com/alicebob/miniredis/v2"
	"github.com/spf13/viper"
)

func TestCacheInitAndReadWriteWithMiniRedis(t *testing.T) {
	t.Cleanup(func() {
		if err := cache.Close(); err != nil {
			t.Fatalf("关闭缓存失败: %v", err)
		}
	})

	miniRedis, err := miniredis.Run()
	if err != nil {
		t.Fatalf("启动 miniredis 失败: %v", err)
	}
	defer miniRedis.Close()

	cfg := viper.New()
	cfg.Set("redis.provider", "redis")
	cfg.Set("redis.addrs", []string{miniRedis.Addr()})
	cfg.Set("redis.password", "")
	cfg.Set("redis.db", 0)

	if err := cache.InitRedis(cfg); err != nil {
		t.Fatalf("初始化缓存失败: %v", err)
	}

	if !cache.IsInited() {
		t.Fatalf("缓存初始化状态错误: 期望=true 实际=false")
	}

	client, err := cache.GetRedis()
	if err != nil {
		t.Fatalf("获取缓存客户端失败: %v", err)
	}

	ctx := context.Background()
	if err := client.Set(ctx, "k:feature", "v:ok", time.Minute).Err(); err != nil {
		t.Fatalf("写入缓存失败: %v", err)
	}

	value, err := client.Get(ctx, "k:feature").Result()
	if err != nil {
		t.Fatalf("读取缓存失败: %v", err)
	}
	if value != "v:ok" {
		t.Fatalf("缓存值不正确: 期望=%s 实际=%s", "v:ok", value)
	}
}
