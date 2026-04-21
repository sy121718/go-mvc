package database

import (
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestParseRuntimeOptionsUsesConfiguredValues(t *testing.T) {
	cfg := viper.New()
	cfg.Set("database.prepare_stmt", true)
	cfg.Set("database.skip_default_transaction", true)
	cfg.Set("database.slow_threshold", "500ms")
	cfg.Set("server.mode", "test")

	options, err := parseRuntimeOptions(cfg)
	if err != nil {
		t.Fatalf("解析数据库运行时配置失败: %v", err)
	}

	if !options.prepareStmt {
		t.Fatalf("prepareStmt 不正确")
	}
	if !options.skipDefaultTransaction {
		t.Fatalf("skipDefaultTransaction 不正确")
	}
	if options.slowThreshold != 500*time.Millisecond {
		t.Fatalf("slowThreshold 不正确: got=%s want=%s", options.slowThreshold, 500*time.Millisecond)
	}
}
