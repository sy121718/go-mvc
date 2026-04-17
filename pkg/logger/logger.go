package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

const (
	defaultBaseDir = "public/logs"
	defaultScene   = "default"
)

type runtimeConfig struct {
	baseDir      string
	initialized  bool
	initializedM sync.RWMutex
}

var (
	cfg     = runtimeConfig{}
	writeMu sync.Mutex
)

// EventLevel 日志级别。
type EventLevel string

const (
	LevelInfo  EventLevel = "info"
	LevelWarn  EventLevel = "warn"
	LevelError EventLevel = "error"
	LevelDebug EventLevel = "debug"
)

type eventRecord struct {
	Time    string                 `json:"time"`
	Level   EventLevel             `json:"level"`
	Scene   string                 `json:"scene"`
	Message string                 `json:"message"`
	Error   string                 `json:"error,omitempty"`
	Fields  map[string]interface{} `json:"fields,omitempty"`
}

// Init 初始化日志配置。
func Init(v *viper.Viper) error {
	baseDir := defaultBaseDir
	if v != nil {
		baseDir = strings.TrimSpace(v.GetString("log.base_dir"))
	}
	if baseDir == "" {
		baseDir = defaultBaseDir
	}

	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return fmt.Errorf("创建日志目录失败: %w", err)
	}

	cfg.initializedM.Lock()
	cfg.baseDir = filepath.Clean(baseDir)
	cfg.initialized = true
	cfg.initializedM.Unlock()

	return nil
}

// Entry 场景日志入口。
type Entry struct {
	scene  string
	fields map[string]interface{}
}

// Scene 指定日志场景（目录）。
// scene 为空时，默认写入 default 目录。
func Scene(scene string) *Entry {
	return &Entry{
		scene:  normalizeScene(scene),
		fields: map[string]interface{}{},
	}
}

// With 添加单个字段。
func (e *Entry) With(key string, value interface{}) *Entry {
	if strings.TrimSpace(key) == "" {
		return e
	}
	e.fields[key] = value
	return e
}

// WithFields 添加多个字段。
func (e *Entry) WithFields(fields map[string]interface{}) *Entry {
	for k, v := range fields {
		if strings.TrimSpace(k) == "" {
			continue
		}
		e.fields[k] = v
	}
	return e
}

// Debug 记录 debug 日志。
func (e *Entry) Debug(msg string) {
	write(e.scene, LevelDebug, msg, nil, e.fields)
}

// Info 记录 info 日志。
func (e *Entry) Info(msg string) {
	write(e.scene, LevelInfo, msg, nil, e.fields)
}

// Warn 记录 warn 日志。
func (e *Entry) Warn(msg string) {
	write(e.scene, LevelWarn, msg, nil, e.fields)
}

// Error 记录 error 日志。
func (e *Entry) Error(err error, msg string) {
	write(e.scene, LevelError, msg, err, e.fields)
}

func write(scene string, level EventLevel, msg string, err error, fields map[string]interface{}) {
	scene = normalizeScene(scene)

	baseDir := defaultBaseDir
	cfg.initializedM.RLock()
	if cfg.initialized {
		baseDir = cfg.baseDir
	}
	cfg.initializedM.RUnlock()

	record := eventRecord{
		Time:    time.Now().Format(time.RFC3339Nano),
		Level:   level,
		Scene:   scene,
		Message: strings.TrimSpace(msg),
	}
	if err != nil {
		record.Error = err.Error()
	}
	if len(fields) > 0 {
		cloned := make(map[string]interface{}, len(fields))
		for k, v := range fields {
			cloned[k] = v
		}
		record.Fields = cloned
	}

	data, marshalErr := json.Marshal(record)
	if marshalErr != nil {
		log.Printf("日志序列化失败: %v", marshalErr)
		return
	}

	filePath := filepath.Join(baseDir, scene, time.Now().Format("2006-01-02")+".log")
	if err := appendLine(filePath, string(data)); err != nil {
		log.Printf("写日志失败 scene=%s path=%s err=%v", scene, filePath, err)
	}
}

func appendLine(filePath string, line string) error {
	writeMu.Lock()
	defer writeMu.Unlock()

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(line + "\n"); err != nil {
		return err
	}
	return nil
}

func normalizeScene(scene string) string {
	normalized := strings.TrimSpace(scene)
	if normalized == "" {
		return defaultScene
	}

	lower := strings.ToLower(normalized)
	switch lower {
	case "false", "off", "none", "null", "0", "disable", "disabled":
		return ""
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
