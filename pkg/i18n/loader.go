package i18n

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go-mvc/pkg/database"
)

var (
	refreshMu     sync.Mutex
	refreshCancel context.CancelFunc
	refreshWG     sync.WaitGroup
)

// LoadCache 从数据库加载多语言数据到内存
func LoadCache() error {
	db, err := database.GetDB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	var rows []struct {
		Key      string
		Lang     string
		Value    string
		HttpCode *int
	}

	err = db.Table("sys_i18n").
		Select("item_key AS key", "lang AS lang", "item_value AS value", "http_code AS http_code").
		Where("status = ?", 1).
		Scan(&rows).Error
	if err != nil {
		return err
	}

	newData := make(map[string]map[string]string)
	newHttpCodes := make(map[string]int)

	for _, r := range rows {
		if newData[r.Key] == nil {
			newData[r.Key] = make(map[string]string)
		}
		newData[r.Key][r.Lang] = r.Value

		if r.HttpCode != nil && newHttpCodes[r.Key] == 0 {
			newHttpCodes[r.Key] = *r.HttpCode
		}
	}

	cache.Update(newData, newHttpCodes)
	log.Printf("[i18n] Loaded %d keys, %d records, version: %d", len(newData), len(rows), cache.GetVersion())
	return nil
}

// StartAutoRefresh 启动自动刷新
func StartAutoRefresh(interval time.Duration) {
	if interval <= 0 {
		interval = 20 * time.Second
	}

	refreshMu.Lock()
	defer refreshMu.Unlock()

	if refreshCancel != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	refreshCancel = cancel
	refreshWG.Add(1)

	go func() {
		defer refreshWG.Done()

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("[i18n] Auto refresh stopped")
				return
			case <-ticker.C:
				if err := LoadCache(); err != nil {
					log.Printf("[i18n] Auto refresh failed: %v", err)
				}
			}
		}
	}()

	log.Printf("[i18n] Auto refresh started, interval: %s", interval)
}

// StopAutoRefresh 停止自动刷新
func StopAutoRefresh() {
	refreshMu.Lock()
	cancel := refreshCancel
	refreshCancel = nil
	refreshMu.Unlock()

	if cancel == nil {
		return
	}

	cancel()
	refreshWG.Wait()
}

// ValidateReady 检查 i18n 是否已经完成初始化
func ValidateReady() error {
	if !IsInited() {
		return fmt.Errorf("i18n 未初始化")
	}
	return nil
}
