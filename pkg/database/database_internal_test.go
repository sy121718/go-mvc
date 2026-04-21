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

func TestDatabaseConfigReservesResolverStructure(t *testing.T) {
	cfg := viper.New()
	cfg.Set("database.driver", "mysql")
	cfg.Set("database.dbname", "main")
	cfg.Set("database.resolver.enabled", true)
	cfg.Set("database.resolver.policy", "random")
	cfg.Set("database.resolver.sources", []string{"db-master:3306"})
	cfg.Set("database.resolver.replicas", []string{"db-replica-1:3306", "db-replica-2:3306"})

	var parsed Config
	if err := cfg.UnmarshalKey("database", &parsed); err != nil {
		t.Fatalf("解析数据库配置失败: %v", err)
	}

	if !parsed.Resolver.Enabled {
		t.Fatalf("resolver.enabled 不正确")
	}
	if parsed.Resolver.Policy != "random" {
		t.Fatalf("resolver.policy 不正确: got=%s want=%s", parsed.Resolver.Policy, "random")
	}
	if len(parsed.Resolver.Sources) != 1 || parsed.Resolver.Sources[0] != "db-master:3306" {
		t.Fatalf("resolver.sources 不正确: %+v", parsed.Resolver.Sources)
	}
	if len(parsed.Resolver.Replicas) != 2 {
		t.Fatalf("resolver.replicas 不正确: %+v", parsed.Resolver.Replicas)
	}
}
