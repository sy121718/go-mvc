package i18n

import (
	"log"
	"time"

	"go-mvc/pkg/database"
)

// LoadCache 从数据库加载多语言数据到内存
func LoadCache() error {
	db := database.GetDB()

	var rows []struct {
		Key      string
		Lang     string
		Value    string
		HttpCode *int // 可为 NULL
	}

	err := db.Table("sys_i18n").
		Select("item_key as `key`, lang, item_value as value, http_code").
		Where("status = ?", 1).
		Scan(&rows).Error

	if err != nil {
		return err
	}

	newData := make(map[string]map[string]string)
	newHttpCodes := make(map[string]int)

	for _, r := range rows {
		// 构建多语言数据
		if newData[r.Key] == nil {
			newData[r.Key] = make(map[string]string)
		}
		newData[r.Key][r.Lang] = r.Value

		// 构建 HTTP 响应码映射（只需存储一次）
		if r.HttpCode != nil && newHttpCodes[r.Key] == 0 {
			newHttpCodes[r.Key] = *r.HttpCode
		}
	}

	// 更新缓存
	cache.Update(newData, newHttpCodes)

	log.Printf("[i18n] Loaded %d keys, %d records, version: %d", len(newData), len(rows), cache.GetVersion())

	return nil
}

// StartAutoRefresh 启动自动刷新（每10秒检查一次）
func StartAutoRefresh() {
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if err := LoadCache(); err != nil {
				log.Printf("[i18n] Auto refresh failed: %v", err)
			}
		}
	}()
}
