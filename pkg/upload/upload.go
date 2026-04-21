package upload

import (
	"context"
	"errors"
	enums "go-mvc/pkg/enums"

	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	uploadprovider "go-mvc/pkg/upload/provider"

	"github.com/spf13/viper"
)

var (
	stateMu         sync.RWMutex
	runtimeMu       sync.RWMutex
	inited          bool
	configSource    *viper.Viper
	defaultProvider = "local"
	providers       = map[string]*providerEntry{}
	uploadRules     = validationRules{
		maxSize:           10 * 1024 * 1024,
		allowedExtensions: map[string]struct{}{},
		allowedMIMETypes:  map[string]struct{}{},
	}
)

type providerEntry struct {
	provider uploadprovider.Provider
	mu       sync.Mutex
	ready    bool
}

type validationRules struct {
	maxSize           int64
	allowedExtensions map[string]struct{}
	allowedMIMETypes  map[string]struct{}
}

func init() {
	_ = Register(uploadprovider.NewLocalProvider())
	_ = Register(uploadprovider.NewQiniuProvider())
}

// File 上传文件参数。
type File = uploadprovider.File

// Request 上传请求参数。
type Request = uploadprovider.Request

// RuntimeConfig 运行时上传配置。
type RuntimeConfig = uploadprovider.RuntimeConfig

// Result 上传结果。
type Result = uploadprovider.Result

// Client 上传客户端。
type Client struct {
	provider string
	runtime  *RuntimeConfig
}

// Init 初始化上传组件。
func Init(v *viper.Viper) error {
	if v == nil {
		return uploadprovider.NewErrorf(enums.ErrUploadConfigInvalid, "upload 初始化配置为空")
	}

	selected := normalizeProvider(v.GetString("upload.default_provider"))
	if selected == "" {
		selected = normalizeProvider(v.GetString("upload.provider"))
	}
	if selected == "" {
		selected = "local"
	}

	stateMu.Lock()
	configSource = v
	defaultProvider = selected
	rules, err := parseValidationRules(v)
	if err != nil {
		stateMu.Unlock()
		return err
	}
	uploadRules = rules
	inited = true
	_, ok := providers[selected]
	stateMu.Unlock()

	if !ok {
		stateMu.Lock()
		inited = false
		stateMu.Unlock()
		return uploadprovider.NewErrorf(enums.ErrUploadProviderNotFound, "默认 provider=%s", selected)
	}

	if err := ensureProviderReady(selected); err != nil {
		stateMu.Lock()
		inited = false
		stateMu.Unlock()
		return err
	}
	return nil
}

// Close 关闭上传组件。
func Close() error {
	runtimeMu.Lock()
	defer runtimeMu.Unlock()

	stateMu.RLock()
	if !inited {
		stateMu.RUnlock()
		return nil
	}
	entries := make([]*providerEntry, 0, len(providers))
	for _, entry := range providers {
		entries = append(entries, entry)
	}
	stateMu.RUnlock()

	var closeErr error
	for _, entry := range entries {
		entry.mu.Lock()
		if entry.ready {
			if err := entry.provider.Close(); err != nil {
				closeErr = errors.Join(closeErr, err)
			}
			entry.ready = false
		}
		entry.mu.Unlock()
	}

	stateMu.Lock()
	inited = false
	configSource = nil
	stateMu.Unlock()
	return closeErr
}

// IsInited 判断上传组件是否已初始化。
func IsInited() bool {
	stateMu.RLock()
	defer stateMu.RUnlock()
	return inited
}

// Register 注册上传 provider。
func Register(provider uploadprovider.Provider) error {
	if provider == nil {
		return uploadprovider.NewError(enums.ErrUploadConfigInvalid)
	}

	name := normalizeProvider(provider.Name())
	if name == "" {
		return uploadprovider.NewError(enums.ErrUploadConfigInvalid)
	}

	stateMu.Lock()
	providers[name] = &providerEntry{provider: provider}
	shouldInit := inited && configSource != nil
	stateMu.Unlock()

	if shouldInit {
		if err := ensureProviderReady(name); err != nil {
			return uploadprovider.WrapError(enums.ErrUploadSystemError, err, "初始化 provider=%s 失败", name)
		}
	}
	return nil
}

// Providers 获取已注册 provider 列表。
func Providers() []string {
	stateMu.RLock()
	defer stateMu.RUnlock()

	result := make([]string, 0, len(providers))
	for name := range providers {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}

// Upload 使用默认 provider 上传。
func Upload(ctx context.Context, file File, req Request) (Result, error) {
	return uploadWithProvider(ctx, "", RuntimeConfig{}, file, req)
}

// UploadWithConfig 使用外部传入运行时配置上传。
func UploadWithConfig(ctx context.Context, runtime RuntimeConfig, file File, req Request) (Result, error) {
	return uploadWithProvider(ctx, runtime.Provider, runtime, file, req)
}

// Use 指定 provider 上传。
func Use(provider string) Client {
	return Client{provider: normalizeProvider(provider)}
}

// UseCfg 使用外部传入运行时配置上传。
func UseCfg(runtime RuntimeConfig) Client {
	cfg := runtime
	return Client{provider: normalizeProvider(cfg.Provider), runtime: &cfg}
}

// Upload 执行上传。
func (c Client) Upload(ctx context.Context, file File, req Request) (Result, error) {
	if c.runtime != nil {
		return uploadWithProvider(ctx, c.provider, *c.runtime, file, req)
	}
	return uploadWithProvider(ctx, c.provider, RuntimeConfig{}, file, req)
}

func uploadWithProvider(ctx context.Context, providerName string, runtime RuntimeConfig, file File, req Request) (Result, error) {
	runtimeMu.RLock()
	defer runtimeMu.RUnlock()

	provider, name, err := getProvider(providerName)
	if err != nil {
		return Result{}, err
	}

	runtime.Provider = name
	if name != "local" && !hasOnlineRuntimeConfig(runtime) {
		return Result{}, uploadprovider.NewError(enums.ErrUploadConfigMissing)
	}
	if err := validateFile(file); err != nil {
		return Result{}, err
	}

	result, err := provider.Upload(ctx, runtime, file, req)
	if err != nil {
		return Result{}, err
	}
	if strings.TrimSpace(result.Provider) == "" {
		result.Provider = name
	}
	return result, nil
}

func getProvider(providerName string) (uploadprovider.Provider, string, error) {
	stateMu.RLock()
	initialized := inited
	name := normalizeProvider(providerName)
	if name == "" {
		name = defaultProvider
	}
	entry, ok := providers[name]
	stateMu.RUnlock()

	if !initialized {
		return nil, "", uploadprovider.NewError(enums.ErrUploadNotInitialized)
	}
	if !ok {
		return nil, "", uploadprovider.NewErrorf(enums.ErrUploadProviderNotFound, "provider=%s", name)
	}
	if err := ensureProviderReady(name); err != nil {
		return nil, "", err
	}
	return entry.provider, name, nil
}

func ensureProviderReady(providerName string) error {
	stateMu.RLock()
	initialized := inited
	cfg := configSource
	entry, ok := providers[providerName]
	stateMu.RUnlock()

	if !initialized {
		return uploadprovider.NewError(enums.ErrUploadNotInitialized)
	}
	if !ok {
		return uploadprovider.NewErrorf(enums.ErrUploadProviderNotFound, "provider=%s", providerName)
	}
	if cfg == nil {
		return uploadprovider.NewError(enums.ErrUploadConfigInvalid)
	}

	entry.mu.Lock()
	defer entry.mu.Unlock()

	if entry.ready {
		return nil
	}

	if err := entry.provider.Init(cfg); err != nil {
		return uploadprovider.WrapError(enums.ErrUploadSystemError, err, "初始化 provider=%s 失败", providerName)
	}
	entry.ready = true
	return nil
}

func hasOnlineRuntimeConfig(cfg RuntimeConfig) bool {
	if strings.TrimSpace(cfg.Endpoint) != "" {
		return true
	}
	if strings.TrimSpace(cfg.Bucket) != "" {
		return true
	}
	if strings.TrimSpace(cfg.Region) != "" {
		return true
	}
	if strings.TrimSpace(cfg.BaseURL) != "" {
		return true
	}
	if strings.TrimSpace(cfg.AccessKey) != "" {
		return true
	}
	if strings.TrimSpace(cfg.SecretKey) != "" {
		return true
	}
	return len(cfg.Extra) > 0
}

func normalizeProvider(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func parseValidationRules(v *viper.Viper) (validationRules, error) {
	rules := validationRules{
		maxSize:           10 * 1024 * 1024,
		allowedExtensions: map[string]struct{}{},
		allowedMIMETypes:  map[string]struct{}{},
	}
	if v == nil {
		return rules, nil
	}

	if raw := strings.TrimSpace(v.GetString("upload.max_size")); raw != "" {
		size, err := parseByteSize(raw)
		if err != nil {
			return validationRules{}, uploadprovider.NewErrorf(enums.ErrUploadConfigInvalid, "upload.max_size 配置无效: %v", err)
		}
		rules.maxSize = size
	}

	for _, ext := range v.GetStringSlice("upload.allowed_extensions") {
		normalized := strings.ToLower(strings.TrimSpace(ext))
		if normalized == "" {
			continue
		}
		if !strings.HasPrefix(normalized, ".") {
			normalized = "." + normalized
		}
		rules.allowedExtensions[normalized] = struct{}{}
	}

	for _, mimeType := range v.GetStringSlice("upload.allowed_mime_types") {
		normalized := strings.ToLower(strings.TrimSpace(mimeType))
		if normalized == "" {
			continue
		}
		rules.allowedMIMETypes[normalized] = struct{}{}
	}

	return rules, nil
}

func validateFile(file File) error {
	if file.Size > 0 && uploadRules.maxSize > 0 && file.Size > uploadRules.maxSize {
		return uploadprovider.NewErrorf(enums.ErrUploadConfigInvalid, "上传文件大小超限: max=%d current=%d", uploadRules.maxSize, file.Size)
	}

	if len(uploadRules.allowedExtensions) > 0 {
		ext := strings.ToLower(filepath.Ext(strings.TrimSpace(file.Filename)))
		if _, ok := uploadRules.allowedExtensions[ext]; !ok {
			return uploadprovider.NewErrorf(enums.ErrUploadConfigInvalid, "上传扩展名不允许: %s", ext)
		}
	}

	if len(uploadRules.allowedMIMETypes) > 0 {
		contentType := strings.ToLower(strings.TrimSpace(file.ContentType))
		if _, ok := uploadRules.allowedMIMETypes[contentType]; !ok {
			return uploadprovider.NewErrorf(enums.ErrUploadConfigInvalid, "上传 MIME 类型不允许: %s", file.ContentType)
		}
	}

	return nil
}

func parseByteSize(raw string) (int64, error) {
	normalized := strings.ToUpper(strings.TrimSpace(raw))
	units := []struct {
		suffix string
		scale  int64
	}{
		{suffix: "KB", scale: 1024},
		{suffix: "MB", scale: 1024 * 1024},
		{suffix: "GB", scale: 1024 * 1024 * 1024},
		{suffix: "B", scale: 1},
	}

	for _, unit := range units {
		if strings.HasSuffix(normalized, unit.suffix) {
			text := strings.TrimSpace(strings.TrimSuffix(normalized, unit.suffix))
			var value int64
			_, err := fmt.Sscanf(text, "%d", &value)
			if err != nil {
				return 0, err
			}
			if value <= 0 {
				return 0, fmt.Errorf("值必须大于 0")
			}
			return value * unit.scale, nil
		}
	}

	var value int64
	_, err := fmt.Sscanf(normalized, "%d", &value)
	if err != nil {
		return 0, err
	}
	if value <= 0 {
		return 0, fmt.Errorf("值必须大于 0")
	}
	return value, nil
}
