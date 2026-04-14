package i18n

import "sync"

var once sync.Once

// Init 初始化多语言配置中心（懒加载）
func Init() {
	once.Do(func() {
		if err := LoadCache(); err != nil {
			panic("Failed to load i18n cache: " + err.Error())
		}
	})
}

// Get 获取多语言信息（返回完整结构）
func Get(key, lang string) *I18nResult {
	Init()
	return cache.Get(key, lang)
}

// Reload 手动重新加载缓存
func Reload() error {
	return LoadCache()
}
