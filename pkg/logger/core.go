package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func newSceneLogger(scene string) (*zap.Logger, io.Closer, error) {
	mu.RLock()
	current := cfg
	currentInited := inited
	mu.RUnlock()

	if !currentInited {
		current = runtimeConfig{
			baseDir:          defaultBaseDir,
			level:            zapcore.InfoLevel,
			sceneLevels:      map[string]zapcore.Level{},
			sampleEnabled:    false,
			sampleInitial:    100,
			sampleThereafter: 100,
			maxSize:          100,
			maxBackups:       10,
			maxAge:           30,
			compress:         false,
		}
	}

	sceneDir := filepath.Join(current.baseDir, scene)
	if err := os.MkdirAll(sceneDir, 0o755); err != nil {
		return nil, nil, fmt.Errorf("创建日志目录失败, scene=%s, dir=%s: %w", scene, sceneDir, err)
	}

	filePath := filepath.Join(sceneDir, defaultLogFile)
	rotator := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    current.maxSize,
		MaxBackups: current.maxBackups,
		MaxAge:     current.maxAge,
		Compress:   current.compress,
	}
	sink := zapcore.AddSync(rotator)

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.LevelKey = "level"
	encoderCfg.MessageKey = "message"
	encoderCfg.CallerKey = "caller"
	encoderCfg.StacktraceKey = "stacktrace"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder

	level := current.level
	if sceneLevel, ok := current.sceneLevels[scene]; ok {
		level = sceneLevel
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		sink,
		level,
	)
	if current.sampleEnabled {
		core = zapcore.NewSamplerWithOptions(core, time.Second, current.sampleInitial, current.sampleThereafter)
	}
	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)), rotator, nil
}

func normalizeScene(scene string) string {
	normalized := strings.TrimSpace(scene)
	if normalized == "" {
		return defaultScene
	}

	normalized = strings.ReplaceAll(normalized, "\\", "/")
	normalized = strings.Trim(normalized, "/")
	for strings.Contains(normalized, "..") {
		normalized = strings.ReplaceAll(normalized, "..", "")
	}
	normalized = strings.Trim(normalized, "/")
	if normalized == "" {
		return defaultScene
	}
	return normalized
}

func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return zapcore.DebugLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
