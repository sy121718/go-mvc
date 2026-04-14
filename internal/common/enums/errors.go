package enums

// 系统级错误码（10xxx）
const (
	ErrSystemError   = "ErrSystemError"   // 系统异常
	ErrDBQueryError  = "ErrDBQueryError"  // 数据库查询错误
	ErrCacheError    = "ErrCacheError"    // 缓存错误
	ErrInvalidParams = "ErrInvalidParams" // 请求参数错误
	ErrInvalidBody   = "ErrInvalidBody"   // 请求体格式错误
)

// 认证错误码（90xxx）
const (
	ErrUnauthorized      = "ErrUnauthorized"      // 未登录或登录已过期
	ErrInvalidToken      = "ErrInvalidToken"      // Token无效
	ErrTokenExpired      = "ErrTokenExpired"      // Token已过期
	ErrPermissionDenied  = "ErrPermissionDenied"  // 无权限访问
)

// 用户模块错误码（20xxx）
const (
	ErrUserNotFound     = "ErrUserNotFound"     // 用户不存在
	ErrUserExists       = "ErrUserExists"       // 用户已存在
	ErrInvalidPassword  = "ErrInvalidPassword"  // 用户名或密码错误
)
