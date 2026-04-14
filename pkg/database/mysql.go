// Package database /*
package database

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
)

/*
数据库组件
===========================================
配置结构体定义在这里，自己解析配置
*/

// Config 数据库配置
type Config struct {
	Host         string `mapstructure:"host"`     //地址
	Port         int    `mapstructure:"port"`     //端口
	User         string `mapstructure:"user"`     //用户名
	Password     string `mapstructure:"password"` //密码
	DBName       string `mapstructure:"dbname"`   //数据库名称
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	LazyInit     bool   `mapstructure:"lazy_init"` //是否懒加载
}

var (
	db     *gorm.DB  // 数据库连接实例（全局单例）
	mu     sync.Mutex // 确保并发安全
	inited bool      // 标记是否已初始化
)

// GetDB 获取数据库实例（懒加载，全局单例）
// 必须先调用 InitDB 初始化，否则返回 nil
func GetDB() *gorm.DB {
	if db == nil {
		panic("数据库未初始化，请先调用 database.InitDB()")
	}
	return db
}

// InitDB 初始化数据库（可重试）
// 内部处理错误，致命错误会直接退出程序
func InitDB(v *viper.Viper) {
	mu.Lock()
	defer mu.Unlock()

	if inited {
		return
	}

	if err := initDB(v); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	inited = true
}

// IsInited 检查是否已初始化
func IsInited() bool {
	return inited
}

// getDefaultConfig 获取默认配置
func getDefaultConfig() Config {
	return Config{
		Host:         "127.0.0.1",
		Port:         3306,
		User:         "root",
		Password:     "",
		DBName:       "test",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		LazyInit:     false,
	}
}

// initDB 初始化数据库
func initDB(v *viper.Viper) error {
	// 自己解析配置
	var cfg Config
	if err := v.UnmarshalKey("database", &cfg); err != nil {
		log.Printf("解析数据库配置失败，使用默认配置: %v", err)
		cfg = getDefaultConfig()
	}

	// 配置兜底：如果关键字段为空，使用默认值
	defaultCfg := getDefaultConfig()
	if cfg.Host == "" {
		cfg.Host = defaultCfg.Host
	}
	if cfg.Port == 0 {
		cfg.Port = defaultCfg.Port
	}
	if cfg.User == "" {
		cfg.User = defaultCfg.User
	}
	if cfg.DBName == "" {
		cfg.DBName = defaultCfg.DBName
	}
	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = defaultCfg.MaxIdleConns
	}
	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = defaultCfg.MaxOpenConns
	}

	// 检查是否懒加载
	if cfg.LazyInit {
		log.Println("数据库懒加载模式，跳过初始化")
		return nil
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 获取底层 sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %v", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	log.Println("数据库初始化成功")
	return nil
}

// Close 关闭数据库连接
func Close() error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %v", err)
	}
	return sqlDB.Close()
}