package logger

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Init 初始化 zap 日志运行时配置。
func Init(v *viper.Viper) error {
	next := runtimeConfig{
		baseDir:    defaultBaseDir,
		level:      zapcore.InfoLevel,
		maxSize:    100,
		maxBackups: 10,
		maxAge:     30,
		compress:   false,
	}

	if v != nil {
		if baseDir := strings.TrimSpace(v.GetString("log.base_dir")); baseDir != "" {
			next.baseDir = baseDir
		}
		next.level = parseLevel(v.GetString("log.level"))
		if size := v.GetInt("log.max_size"); size > 0 {
			next.maxSize = size
		}
		if backups := v.GetInt("log.max_backups"); backups > 0 {
			next.maxBackups = backups
		}
		if age := v.GetInt("log.max_age"); age > 0 {
			next.maxAge = age
		}
		next.compress = v.GetBool("log.compress")
	}

	next.baseDir = filepath.Clean(next.baseDir)
	if err := os.MkdirAll(next.baseDir, 0o755); err != nil {
		return fmt.Errorf("创建日志目录失败: %w", err)
	}

	mu.Lock()
	cfg = next
	inited = true
	mu.Unlock()

	if err := clearSceneLoggers(); err != nil {
		return fmt.Errorf("清理旧 logger 失败: %w", err)
	}
	return nil
}

// Sync 刷新并落盘所有场景 logger。
func Sync() error {
	var syncErr error
	sceneLoggers.Range(func(key, value any) bool {
		lg, ok := value.(*zap.Logger)
		if !ok {
			return true
		}
		err := syncOneLogger(key, lg)
		syncErr = errors.Join(syncErr, err)
		return true
	})
	return syncErr
}

func getSceneLogger(scene string) (*zap.Logger, error) {
	scene = normalizeScene(scene)

	if value, ok := sceneLoggers.Load(scene); ok {
		if lg, castOK := value.(*zap.Logger); castOK {
			return lg, nil
		}
	}

	lg, err := newSceneLogger(scene)
	if err != nil {
		return nil, err
	}

	actual, loaded := sceneLoggers.LoadOrStore(scene, lg)
	if loaded {
		if err := syncOneLogger(scene, lg); err != nil {
			log.Printf("关闭重复创建的 logger 失败, scene=%s, err=%v", scene, err)
		}
		if existing, ok := actual.(*zap.Logger); ok {
			return existing, nil
		}
	}
	return lg, nil
}

func clearSceneLoggers() error {
	var syncErr error
	sceneLoggers.Range(func(key, value any) bool {
		if lg, ok := value.(*zap.Logger); ok {
			err := syncOneLogger(key, lg)
			syncErr = errors.Join(syncErr, err)
		}
		sceneLoggers.Delete(key)
		return true
	})
	return syncErr
}

func ensureInited() error {
	mu.RLock()
	currentInited := inited
	mu.RUnlock()
	if currentInited {
		return nil
	}

	return Init(nil)
}

func syncOneLogger(scene any, lg *zap.Logger) error {
	if lg == nil {
		return nil
	}
	if err := lg.Sync(); err != nil {
		return fmt.Errorf("刷新场景日志失败, scene=%v: %w", scene, err)
	}
	return nil
}
