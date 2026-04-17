package logger_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"go-mvc/pkg/logger"

	"github.com/spf13/viper"
)

func TestLoggerInitReturnsErrorWhenBaseDirIsFile(t *testing.T) {
	tmpDir := t.TempDir()
	basePath := filepath.Join(tmpDir, "not_dir")
	if err := os.WriteFile(basePath, []byte("x"), 0o644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	cfg := viper.New()
	cfg.Set("log.base_dir", basePath)
	if err := logger.Init(cfg); err == nil {
		t.Fatalf("日志初始化应当失败，但返回 nil")
	}
}

func TestLoggerConcurrentWriteAndSync(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "go-mvc-logger-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	cfg := viper.New()
	cfg.Set("log.base_dir", tmpDir)
	cfg.Set("log.level", "debug")

	if err := logger.Init(cfg); err != nil {
		t.Fatalf("日志初始化失败: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 40; i++ {
		idx := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				logger.Scene("concurrency").With("worker", idx).With("seq", j).Info("并发日志测试")
			}
		}()
	}
	wg.Wait()

	if err := logger.Sync(); err != nil {
		t.Fatalf("日志 Sync 失败: %v", err)
	}

	logPath := filepath.Join(tmpDir, "concurrency", "app.log")
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("读取日志文件失败: %v", err)
	}
	if len(content) == 0 {
		t.Fatalf("日志文件内容为空: %s", logPath)
	}
}
