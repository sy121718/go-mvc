package config

import (
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestGetServerParsesHTTPTimeouts(t *testing.T) {
	oldV := v
	t.Cleanup(func() {
		v = oldV
	})

	cfg := viper.New()
	cfg.Set("server.port", 8080)
	cfg.Set("server.mode", "test")
	cfg.Set("server.app_name", "go-mvc")
	cfg.Set("server.read_header_timeout", "3s")
	cfg.Set("server.read_timeout", "15s")
	cfg.Set("server.write_timeout", "30s")
	cfg.Set("server.idle_timeout", "60s")
	cfg.Set("server.request_body_limit", "2MB")
	cfg.Set("server.upload_body_limit", "32MB")
	cfg.Set("server.rate_limit_enabled", true)
	cfg.Set("server.rate_limit_limit", 120)
	cfg.Set("server.rate_limit_window", "1m")
	v = cfg

	serverCfg, err := GetServer()
	if err != nil {
		t.Fatalf("获取服务配置失败: %v", err)
	}

	if serverCfg.ReadHeaderTimeout != 3*time.Second {
		t.Fatalf("ReadHeaderTimeout 不正确: got=%s want=%s", serverCfg.ReadHeaderTimeout, 3*time.Second)
	}
	if serverCfg.ReadTimeout != 15*time.Second {
		t.Fatalf("ReadTimeout 不正确: got=%s want=%s", serverCfg.ReadTimeout, 15*time.Second)
	}
	if serverCfg.WriteTimeout != 30*time.Second {
		t.Fatalf("WriteTimeout 不正确: got=%s want=%s", serverCfg.WriteTimeout, 30*time.Second)
	}
	if serverCfg.IdleTimeout != 60*time.Second {
		t.Fatalf("IdleTimeout 不正确: got=%s want=%s", serverCfg.IdleTimeout, 60*time.Second)
	}
	if serverCfg.RequestBodyLimit != 2*1024*1024 {
		t.Fatalf("RequestBodyLimit 不正确: got=%d want=%d", serverCfg.RequestBodyLimit, 2*1024*1024)
	}
	if serverCfg.UploadBodyLimit != 32*1024*1024 {
		t.Fatalf("UploadBodyLimit 不正确: got=%d want=%d", serverCfg.UploadBodyLimit, 32*1024*1024)
	}
	if !serverCfg.RateLimitEnabled {
		t.Fatalf("RateLimitEnabled 不正确: got=%t want=%t", serverCfg.RateLimitEnabled, true)
	}
	if serverCfg.RateLimitLimit != 120 {
		t.Fatalf("RateLimitLimit 不正确: got=%d want=%d", serverCfg.RateLimitLimit, 120)
	}
	if serverCfg.RateLimitWindow != time.Minute {
		t.Fatalf("RateLimitWindow 不正确: got=%s want=%s", serverCfg.RateLimitWindow, time.Minute)
	}
}

func TestValidateRuntimeConfigFailsForReleaseDefaultJWTSecret(t *testing.T) {
	oldV := v
	t.Cleanup(func() {
		v = oldV
	})

	cfg := viper.New()
	setDefaults(cfg)
	cfg.Set("server.mode", "release")
	cfg.Set("database.dbname", "base")
	cfg.Set("database.password", "secret")
	v = cfg

	err := ValidateRuntimeConfig()
	if err == nil {
		t.Fatalf("release 模式下默认 JWT secret 应当校验失败")
	}
	if !strings.Contains(err.Error(), "jwt.secret") {
		t.Fatalf("错误信息应包含 jwt.secret, got=%v", err)
	}
}

func TestValidateRuntimeConfigFailsForReleaseDefaultDatabaseName(t *testing.T) {
	oldV := v
	t.Cleanup(func() {
		v = oldV
	})

	cfg := viper.New()
	setDefaults(cfg)
	cfg.Set("server.mode", "release")
	cfg.Set("jwt.secret", "custom-secret")
	cfg.Set("database.password", "secret")
	v = cfg

	err := ValidateRuntimeConfig()
	if err == nil {
		t.Fatalf("release 模式下默认数据库名应当校验失败")
	}
	if !strings.Contains(err.Error(), "database.dbname") {
		t.Fatalf("错误信息应包含 database.dbname, got=%v", err)
	}
}

func TestValidateRuntimeConfigPassesForReleaseCustomValues(t *testing.T) {
	oldV := v
	t.Cleanup(func() {
		v = oldV
	})

	cfg := viper.New()
	setDefaults(cfg)
	cfg.Set("server.mode", "release")
	cfg.Set("jwt.secret", "custom-secret")
	cfg.Set("database.dbname", "base")
	cfg.Set("database.password", "secret")
	v = cfg

	if err := ValidateRuntimeConfig(); err != nil {
		t.Fatalf("release 模式下自定义关键配置应通过校验: %v", err)
	}
}
