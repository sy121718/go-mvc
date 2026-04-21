package i18n_test

import (
	"testing"
	"time"

	"go-mvc/pkg/database"
	"go-mvc/pkg/i18n"

	"github.com/spf13/viper"
)

type testI18nRecord struct {
	ID        uint   `gorm:"primaryKey"`
	ItemKey   string `gorm:"column:item_key"`
	Lang      string `gorm:"column:lang"`
	ItemValue string `gorm:"column:item_value"`
	HttpCode  int    `gorm:"column:http_code"`
	Status    int    `gorm:"column:status"`
}

func (testI18nRecord) TableName() string {
	return "sys_i18n"
}

func TestI18nInitUsesConfigAndAutoRefresh(t *testing.T) {
	t.Cleanup(func() {
		if err := i18n.Close(); err != nil {
			t.Fatalf("关闭 i18n 失败: %v", err)
		}
		if err := database.Close(); err != nil {
			t.Fatalf("关闭数据库失败: %v", err)
		}
	})

	cfg := viper.New()
	cfg.Set("server.mode", "test")
	cfg.Set("database.driver", "sqlite")
	cfg.Set("database.dbname", ":memory:")
	cfg.Set("database.max_idle_conns", 1)
	cfg.Set("database.max_open_conns", 1)
	cfg.Set("i18n.default_lang", "en-US")
	cfg.Set("i18n.auto_refresh", true)
	cfg.Set("i18n.refresh_interval", "20ms")

	if err := database.InitDB(cfg); err != nil {
		t.Fatalf("初始化数据库失败: %v", err)
	}

	db, err := database.GetDB()
	if err != nil {
		t.Fatalf("获取数据库实例失败: %v", err)
	}

	if err := db.AutoMigrate(&testI18nRecord{}); err != nil {
		t.Fatalf("迁移 sys_i18n 失败: %v", err)
	}

	seed := []testI18nRecord{
		{ItemKey: "msg_operation_success", Lang: "zh-CN", ItemValue: "操作成功", HttpCode: 200, Status: 1},
		{ItemKey: "msg_operation_success", Lang: "en-US", ItemValue: "Operation successful", HttpCode: 200, Status: 1},
	}
	if err := db.Create(&seed).Error; err != nil {
		t.Fatalf("写入 i18n 测试数据失败: %v", err)
	}

	if err := i18n.Init(cfg); err != nil {
		t.Fatalf("初始化 i18n 失败: %v", err)
	}

	if got := i18n.GetDefaultLang(); got != "en-US" {
		t.Fatalf("默认语言不正确: got=%s want=%s", got, "en-US")
	}

	if got := i18n.GetText("msg_operation_success", ""); got != "Operation successful" {
		t.Fatalf("默认语言文案不正确: got=%s want=%s", got, "Operation successful")
	}

	if err := db.Model(&testI18nRecord{}).
		Where("item_key = ? AND lang = ?", "msg_operation_success", "en-US").
		Update("item_value", "Operation refreshed").Error; err != nil {
		t.Fatalf("更新 i18n 数据失败: %v", err)
	}

	waitForText(t, func() string {
		return i18n.GetText("msg_operation_success", "")
	}, "Operation refreshed")
}

func waitForText(t *testing.T, getter func() string, expected string) {
	t.Helper()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if getter() == expected {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}

	t.Fatalf("等待 i18n 自动刷新超时: want=%s got=%s", expected, getter())
}
