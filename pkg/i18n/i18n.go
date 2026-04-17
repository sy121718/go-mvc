package i18n

import (
	"fmt"
	"strings"
	"sync"
)

const fallbackDefaultLang string = "zh-CN"

var (
	initMu      sync.Mutex
	inited      bool
	defaultLang = fallbackDefaultLang
)

// Init initializes i18n cache data.
func Init() error {
	initMu.Lock()
	defer initMu.Unlock()

	if inited {
		return nil
	}

	if err := LoadCache(); err != nil {
		return fmt.Errorf("failed to load i18n cache: %w", err)
	}

	inited = true
	return nil
}

// SetDefaultLang sets default language code.
func SetDefaultLang(lang string) {
	initMu.Lock()
	defer initMu.Unlock()

	lang = strings.TrimSpace(lang)
	if lang == "" {
		defaultLang = fallbackDefaultLang
		return
	}

	defaultLang = lang
}

// GetDefaultLang returns default language code.
func GetDefaultLang() string {
	initMu.Lock()
	defer initMu.Unlock()
	return defaultLang
}

// Get returns full i18n result.
//
// Example:
//
//	result := i18n.Get("ErrUploadConfigMissing", "zh-CN")
//	// result.Key      == "ErrUploadConfigMissing"
//	// result.Value    == "上传配置缺失"
//	// result.HttpCode == 400
//	// result.Lang     == "zh-CN"
//
// Fields:
//   - Key: code/text key
//   - Value: localized text
//   - Lang: matched language
//   - HttpCode: mapped HTTP status
//   - AllLangs: all language versions for this key
func Get(key, lang string) *I18nResult {
	lang = strings.TrimSpace(lang)
	if lang == "" {
		lang = GetDefaultLang()
	}
	return cache.Get(key, lang)
}

// GetText returns localized text only.
func GetText(key, lang string) string {
	result := Get(key, lang)
	if result == nil {
		return key
	}
	return result.Value
}

// GetHttpCode returns mapped HTTP status code for key.
func GetHttpCode(key string) int {
	result := Get(key, GetDefaultLang())
	if result == nil {
		return 200
	}
	return result.HttpCode
}

// Reload reloads i18n cache.
func Reload() error {
	if err := LoadCache(); err != nil {
		return err
	}

	initMu.Lock()
	inited = true
	initMu.Unlock()
	return nil
}

// IsInited reports whether i18n has been initialized.
func IsInited() bool {
	initMu.Lock()
	defer initMu.Unlock()
	return inited
}
