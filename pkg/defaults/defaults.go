// Package defaults 统一管理基础框架的共享默认值。
//
// 目的：
// - 减少 config 与各 pkg 内重复定义默认值带来的漂移
// - 让关键默认值只有一个来源
// - 后续新增默认值时，优先在这里补，再由 config/pkg 引用
package defaults

const (
	DefaultDatabaseDriver       = "mysql"
	DefaultDatabaseHost         = "127.0.0.1"
	DefaultDatabasePort         = 3306
	DefaultDatabaseUser         = "root"
	DefaultDatabasePassword     = ""
	DefaultDatabaseName         = "test"
	DefaultDatabaseMaxIdleConns = 10
	DefaultDatabaseMaxOpenConns = 100

	DefaultRedisProvider = "redis"
	DefaultRedisHost     = "127.0.0.1"
	DefaultRedisPort     = 6379
	DefaultRedisPassword = ""
	DefaultRedisDB       = 0

	DefaultJWTSecret     = "default-secret-key-please-change-in-production"
	DefaultJWTExpireTime = 24
	DefaultJWTIssuer     = "go-mvc"

	DefaultQueueProvider    = "asynq"
	DefaultQueueConcurrency = 10
)
