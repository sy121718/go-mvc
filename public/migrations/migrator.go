package migrations

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"gorm.io/gorm"
)

type Migration struct {
	Version   string
	TableName string
	SQL       string
}

type Seed struct {
	Version      string
	TableName    string
	ConditionSQL string
	SQL          string
}

var allMigrations []Migration
var allSeeds []Seed

func register(m Migration) {
	allMigrations = append(allMigrations, m)
}

func registerSeed(s Seed) {
	allSeeds = append(allSeeds, s)
}

func All() []Migration {
	sort.Slice(allMigrations, func(i, j int) bool {
		return allMigrations[i].Version < allMigrations[j].Version
	})
	return allMigrations
}

func AllSeeds() []Seed {
	sort.Slice(allSeeds, func(i, j int) bool {
		return allSeeds[i].Version < allSeeds[j].Version
	})
	return allSeeds
}

func Run(db *gorm.DB) error {
	for _, m := range All() {
		if err := apply(db, m); err != nil {
			return fmt.Errorf("迁移 %s (%s) 失败: %w", m.Version, m.TableName, err)
		}
	}
	return nil
}

func apply(db *gorm.DB, m Migration) error {
	var count int64
	db.Raw("SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?", m.TableName).Scan(&count)
	if count > 0 {
		log.Printf("跳过 %s (%s)：表已存在", m.Version, m.TableName)
		return nil
	}

	parts := strings.Split(m.SQL, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if err := db.Exec(part + ";").Error; err != nil {
			return fmt.Errorf("执行 %s (%s) 失败: %w", m.Version, m.TableName, err)
		}
	}

	log.Printf("完成 %s (%s)", m.Version, m.TableName)
	return nil
}

func RunSeeds(db *gorm.DB) error {
	for _, s := range AllSeeds() {
		if err := applySeed(db, s); err != nil {
			return fmt.Errorf("种子数据 %s (%s) 失败: %w", s.Version, s.TableName, err)
		}
	}
	return nil
}

func applySeed(db *gorm.DB, s Seed) error {
	var count int64
	if err := db.Raw(s.ConditionSQL).Scan(&count).Error; err != nil {
		return fmt.Errorf("检查种子数据 %s (%s) 失败: %w", s.Version, s.TableName, err)
	}
	if count > 0 {
		log.Printf("跳过种子 %s (%s)：数据已存在", s.Version, s.TableName)
		return nil
	}

	if err := db.Exec(s.SQL).Error; err != nil {
		return fmt.Errorf("执行种子 %s (%s) 失败: %w", s.Version, s.TableName, err)
	}

	log.Printf("完成种子 %s (%s)", s.Version, s.TableName)
	return nil
}
