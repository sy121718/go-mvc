/*
Redis 缓存组件包
===========================================
提供 Redis 连接管理功能

主要功能：
- Redis 连接初始化
- 全局单例模式
- 缓存操作封装

配置说明（config.yaml）：

	redis:
	  host: 127.0.0.1
	  port: 6379
	  password: ""
	  db: 0
	  enabled: true
*/
package cacheprovider

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// Config Redis配置
type Config struct {
	Host     string   `mapstructure:"host"`
	Port     int      `mapstructure:"port"`
	Password string   `mapstructure:"password"`
	DB       int      `mapstructure:"db"`
	Addrs    []string `mapstructure:"addrs"`
}

type redisProvider struct {
	rdb    redis.UniversalClient
	mu     sync.Mutex
	inited bool
}

func getDefaultConfig() Config {
	return Config{
		Host:     "127.0.0.1",
		Port:     6379,
		Password: "",
		DB:       0,
	}
}

func normalizeConfig(cfg Config) Config {
	defaultCfg := getDefaultConfig()
	if cfg.Host == "" {
		cfg.Host = defaultCfg.Host
	}
	if cfg.Port == 0 {
		cfg.Port = defaultCfg.Port
	}
	if len(cfg.Addrs) == 0 {
		cfg.Addrs = []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)}
	}
	return cfg
}

func (p *redisProvider) initRedis(v *viper.Viper) (redis.UniversalClient, error) {
	var cfg Config
	if err := v.UnmarshalKey("redis", &cfg); err != nil {
		log.Printf("解析 Redis 配置失败，使用默认配置: %v", err)
		cfg = getDefaultConfig()
	}
	cfg = normalizeConfig(cfg)

	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    cfg.Addrs,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("Redis 连接失败: %w", err)
	}

	log.Println("Redis 初始化成功")
	return client, nil
}

func (p *redisProvider) Init(v *viper.Viper) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.inited {
		return nil
	}

	client, err := p.initRedis(v)
	if err != nil {
		return fmt.Errorf("Redis 初始化失败: %w", err)
	}

	p.rdb = client
	p.inited = true
	return nil
}

func (p *redisProvider) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.rdb == nil {
		return nil
	}

	if err := p.rdb.Close(); err != nil {
		return err
	}

	p.rdb = nil
	p.inited = false
	return nil
}

func (p *redisProvider) Client() redis.UniversalClient {
	if p.rdb == nil {
		panic("Redis 未初始化，请先调用 cache.InitRedis()")
	}
	return p.rdb
}

func (p *redisProvider) IsInited() bool {
	return p.inited && p.rdb != nil
}

// NewRedisProvider 创建 Redis 实现
func NewRedisProvider() Provider {
	return &redisProvider{}
}
