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
	db   *gorm.DB  // 数据库连接实例（全局单例）
	once sync.Once // 确保只初始化一次（并发安全）
)

// GetDB 获取数据库实例（懒加载）
func GetDB(v *viper.Viper) *gorm.DB {
	once.Do(func() {
		if err := initDB(v); err != nil {
			panic(fmt.Sprintf("数据库初始化失败: %v", err))
		}
	})
	return db
}

// initDB 初始化数据库
func initDB(v *viper.Viper) error {
	// 自己解析配置
	var cfg Config
	if err := v.UnmarshalKey("database", &cfg); err != nil {
		return fmt.Errorf("解析数据库配置失败: %v", err)
	}

	// 检查是否懒加载
	if cfg.LazyInit {
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

// InitDB 手动初始化数据库（用于非懒加载场景）
func InitDB(v *viper.Viper) error {
	return initDB(v)
}
