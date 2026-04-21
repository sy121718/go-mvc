package queue_test

import (
	"context"
	"encoding/json"
	"net"
	"strconv"
	"testing"
	"time"

	"go-mvc/pkg/queue"

	"github.com/alicebob/miniredis/v2"
	"github.com/spf13/viper"
)

func TestQueueInitAndConsumeWithMiniRedis(t *testing.T) {
	t.Cleanup(func() {
		if err := queue.Shutdown(); err != nil {
			t.Fatalf("关闭队列失败: %v", err)
		}
	})

	miniRedis, err := miniredis.Run()
	if err != nil {
		t.Fatalf("启动 miniredis 失败: %v", err)
	}
	defer miniRedis.Close()

	host, port, err := splitHostPort(miniRedis.Addr())
	if err != nil {
		t.Fatalf("解析 miniredis 地址失败: %v", err)
	}

	cfg := viper.New()
	cfg.Set("queue.provider", "asynq")
	cfg.Set("queue.concurrency", 1)
	cfg.Set("queue.redis.host", host)
	cfg.Set("queue.redis.port", port)
	cfg.Set("queue.redis.password", "")
	cfg.Set("queue.redis.db", 0)

	if err := queue.Init(cfg); err != nil {
		t.Fatalf("初始化队列失败: %v", err)
	}

	resultCh := make(chan string, 1)
	queue.Register("feature:echo", func(_ context.Context, payload []byte) error {
		var msg struct {
			Value string `json:"value"`
		}
		if err := json.Unmarshal(payload, &msg); err != nil {
			return err
		}
		resultCh <- msg.Value
		return nil
	})

	if err := queue.Start(); err != nil {
		t.Fatalf("启动队列失败: %v", err)
	}

	if err := queue.Enqueue("feature:echo", map[string]string{"value": "ok"}); err != nil {
		t.Fatalf("入队任务失败: %v", err)
	}

	select {
	case value := <-resultCh:
		if value != "ok" {
			t.Fatalf("队列消费结果不正确: 期望=%s 实际=%s", "ok", value)
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("等待队列消费超时")
	}
}

func TestQueueInitStartsWorkerAndKeepsRegistrationsWhenConfigured(t *testing.T) {
	miniRedis, err := miniredis.Run()
	if err != nil {
		t.Fatalf("启动 miniredis 失败: %v", err)
	}
	t.Cleanup(func() {
		miniRedis.Close()
	})
	t.Cleanup(func() {
		if err := queue.Close(); err != nil {
			t.Fatalf("关闭队列失败: %v", err)
		}
	})

	host, port, err := splitHostPort(miniRedis.Addr())
	if err != nil {
		t.Fatalf("解析 miniredis 地址失败: %v", err)
	}

	resultCh := make(chan string, 1)
	queue.Register("feature:init-start", func(_ context.Context, payload []byte) error {
		var msg struct {
			Value string `json:"value"`
		}
		if err := json.Unmarshal(payload, &msg); err != nil {
			return err
		}
		resultCh <- msg.Value
		return nil
	})

	cfg := viper.New()
	cfg.Set("queue.provider", "asynq")
	cfg.Set("queue.run_worker", true)
	cfg.Set("queue.concurrency", 1)
	cfg.Set("queue.redis.host", host)
	cfg.Set("queue.redis.port", port)
	cfg.Set("queue.redis.password", "")
	cfg.Set("queue.redis.db", 0)

	if err := queue.Init(cfg); err != nil {
		t.Fatalf("初始化队列失败: %v", err)
	}

	if err := queue.Enqueue("feature:init-start", map[string]string{"value": "started"}); err != nil {
		t.Fatalf("入队任务失败: %v", err)
	}

	select {
	case value := <-resultCh:
		if value != "started" {
			t.Fatalf("任务消费结果不正确: 期望=%s 实际=%s", "started", value)
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("等待任务消费超时")
	}
}

func splitHostPort(address string) (string, int, error) {
	host, portText, err := net.SplitHostPort(address)
	if err != nil {
		return "", 0, err
	}

	port, err := strconv.Atoi(portText)
	if err != nil {
		return "", 0, err
	}

	return host, port, nil
}
