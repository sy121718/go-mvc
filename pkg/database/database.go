// Package database /*
package database

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	dbdriver "go-mvc/pkg/database/driver"

	"github.com/spf13/viper"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

/*
数据库组件
===========================================
配置结构体定义在这里，自己解析配置
*/

// Config 数据库配置
type Config struct {
	Driver       string `mapstructure:"driver"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	LogLevel     string `mapstructure:"log_level"`
}

var (
	db     *gorm.DB
	mu     sync.Mutex
	inited bool
)

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	if db == nil {
		panic("数据库未初始化，请先调用 database.InitDB()")
	}
	return db
}

// InitDB 初始化数据库
func InitDB(v *viper.Viper) error {
	mu.Lock()
	defer mu.Unlock()

	if inited {
		return nil
	}

	instance, err := initDB(v)
	if err != nil {
		return fmt.Errorf("数据库初始化失败: %w", err)
	}

	db = instance
	inited = true
	return nil
}

// IsInited 检查是否已初始化
func IsInited() bool {
	return inited && db != nil
}

func getDefaultConfig() Config {
	return Config{
		Driver:       "mysql",
		Host:         "127.0.0.1",
		Port:         3306,
		User:         "root",
		Password:     "",
		DBName:       "test",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		LogLevel:     "",
	}
}

func resolveLogLevel(serverMode string, dbLogLevel string) gormlogger.LogLevel {
	switch strings.ToLower(dbLogLevel) {
	case "silent":
		return gormlogger.Silent
	case "error":
		return gormlogger.Error
	case "warn", "warning":
		return gormlogger.Warn
	case "info":
		return gormlogger.Info
	}

	switch strings.ToLower(serverMode) {
	case "release":
		return gormlogger.Warn
	case "test":
		return gormlogger.Error
	default:
		return gormlogger.Info
	}
}

func toDriverConfig(cfg Config) dbdriver.Config {
	return dbdriver.Config{
		Driver:       cfg.Driver,
		Host:         cfg.Host,
		Port:         cfg.Port,
		User:         cfg.User,
		Password:     cfg.Password,
		DBName:       cfg.DBName,
		MaxIdleConns: cfg.MaxIdleConns,
		MaxOpenConns: cfg.MaxOpenConns,
	}
}

func buildDialector(cfg Config) (gorm.Dialector, error) {
	driverCfg := dbdriver.NormalizeConfig(toDriverConfig(cfg))
	switch driverCfg.Driver {
	case "mysql":
		return dbdriver.OpenMySQL(driverCfg), nil
	case "postgres", "postgresql":
		return dbdriver.OpenPostgres(driverCfg), nil
	case "sqlserver", "mssql":
		return dbdriver.OpenSQLServer(driverCfg), nil
	case "sqlite":
		return dbdriver.OpenSQLite(driverCfg), nil
	default:
		return nil, fmt.Errorf("不支持的数据库驱动: %s", driverCfg.Driver)
	}
}

// initDB 初始化数据库
func initDB(v *viper.Viper) (*gorm.DB, error) {
	var cfg Config
	if err := v.UnmarshalKey("database", &cfg); err != nil {
		log.Printf("解析数据库配置失败，使用默认配置: %v", err)
		cfg = getDefaultConfig()
	}

	driverCfg := dbdriver.NormalizeConfig(toDriverConfig(cfg))
	cfg.Driver = driverCfg.Driver
	cfg.Host = driverCfg.Host
	cfg.Port = driverCfg.Port
	cfg.User = driverCfg.User
	cfg.Password = driverCfg.Password
	cfg.DBName = driverCfg.DBName
	cfg.MaxIdleConns = driverCfg.MaxIdleConns
	cfg.MaxOpenConns = driverCfg.MaxOpenConns

	dialector, err := buildDialector(cfg)
	if err != nil {
		return nil, err
	}

	sqlScene := ""
	if v.GetBool("log.capture.sql") {
		sqlScene = "sql"
	}

	gormBaseLogger := gormlogger.Default.LogMode(resolveLogLevel(v.GetString("server.mode"), cfg.LogLevel))
	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: newSceneGormLogger(gormBaseLogger, sqlScene),
		// 启用方言错误翻译，便于后续通过 gorm.ErrDuplicatedKey / gorm.ErrForeignKeyViolated 做统一判断。
		TranslateError: true,
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接测试失败: %w", err)
	}

	log.Printf("数据库初始化成功，driver=%s", strings.ToLower(cfg.Driver))
	return gormDB, nil
}

// Close 关闭数据库连接
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return err
	}

	db = nil
	inited = false
	return nil
}
