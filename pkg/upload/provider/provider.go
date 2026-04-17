package uploadprovider

import (
	"context"
	"fmt"
	enums "go-mvc/pkg/enums"

	"go-mvc/pkg/i18n"
	"io"

	"github.com/spf13/viper"
)

// File 上传文件数据。
type File struct {
	Filename    string
	Reader      io.Reader
	Size        int64
	ContentType string
}

// Request 上传请求参数。
type Request struct {
	Route        string
	Directory    string
	ObjectKey    string
	PreserveName bool
}

// RuntimeConfig 按标识加载的线上配置。
type RuntimeConfig struct {
	Mark      string
	Provider  string
	Endpoint  string
	Bucket    string
	Region    string
	BaseURL   string
	AccessKey string
	SecretKey string
	Extra     map[string]string
}

// Result 上传结果。
type Result struct {
	Provider string
	Key      string
	URL      string
	Size     int64
}

// Provider 上传实现接口。
type Provider interface {
	Name() string
	Init(v *viper.Viper) error
	Close() error
	Upload(ctx context.Context, cfg RuntimeConfig, file File, req Request) (Result, error)
}

// Msg 获取错误码对应的字典文案。
func Msg(code string) string {
	if code == "" {
		code = enums.ErrUploadSystemError
	}
	result := i18n.Get(code, i18n.GetDefaultLang())
	if result == nil || result.Value == "" {
		return code
	}
	return result.Value
}

// NewError 创建基于字典文案的错误。
func NewError(code string) error {
	return fmt.Errorf("%s", Msg(code))
}

// NewErrorf 创建带上下文信息的错误。
func NewErrorf(code string, format string, args ...interface{}) error {
	if format == "" {
		return NewError(code)
	}
	return fmt.Errorf("%s: %s", Msg(code), fmt.Sprintf(format, args...))
}

// WrapError 包装底层错误。
func WrapError(code string, cause error, format string, args ...interface{}) error {
	if cause == nil {
		return NewErrorf(code, format, args...)
	}

	if format == "" {
		return fmt.Errorf("%s: %w", Msg(code), cause)
	}
	return fmt.Errorf("%s: %s: %w", Msg(code), fmt.Sprintf(format, args...), cause)
}
