package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultBaseDir = "public/logs"
	defaultScene   = "default"
	defaultLogFile = "app.log"
)

type runtimeConfig struct {
	baseDir    string
	level      zapcore.Level
	maxSize    int
	maxBackups int
	maxAge     int
	compress   bool
}

// Entry 场景日志入口。
type Entry struct {
	scene            string
	autoSceneByLevel bool
	fields           []zap.Field
	initErr          error
}

var (
	mu           sync.RWMutex
	inited       bool
	cfg          runtimeConfig
	sceneLoggers sync.Map // key: scene, value: *zap.Logger
)
