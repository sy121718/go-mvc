package database

import (
	"fmt"
	"go-mvc/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
	"time"
)

var (
	db   *gorm.DB
	once sync.Once
)

// GetDB 获取数据库实例（懒加载）
func GetDB() *gorm.DB {
	once.Do(func() {
		if err := initDB(); err != nil {
			panic(fmt.Sprintf("数据库初始化失败: %v", err))
		}
	})
	return db
}

// initDB 初始化数据库
func initDB() error {
	cfg := config.GetDatabase()

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

	return nil
}

// InitDB 手动初始化数据库（用于懒加载场景）
func InitDB() error {
	return initDB()
}