/*
Package errors 系统级错误码定义

职责范围：
- HTTP 标准错误码
- 参数验证错误
- 数据库操作错误
- 缓存操作错误

业务错误码：
- 各业务模块应在自己的 errors 包中定义
- 错误码范围分配：
  - 系统级：0-4999
  - Admin 模块：5000-5999
  - User 模块：6000-6999
  - Order 模块：7000-7999
  - 其他模块：8000+
*/
package errors

// 错误码定义
const (
	// 通用错误码（HTTP 标准）
	Success      = 200 // 请求成功
	BadRequest   = 400 // 请求参数错误或业务逻辑错误
	Unauthorized = 401 // 未授权，需要登录或 Token 无效
	Forbidden    = 403 // 已登录但无权限访问
	NotFound     = 404 // 请求的资源不存在
	ServerError  = 500 // 服务器内部错误

	// 参数相关 1000-1999
	ParamInvalid     = 1000 // 参数无效（格式正确但不符合业务规则）
	ParamMissing     = 1001 // 缺少必需参数
	ParamFormatError = 1002 // 参数格式错误（类型、长度等）

	// 认证相关 2000-2999（系统级认证，非业务）
	TokenInvalid = 2005 // Token 格式错误或签名验证失败
	TokenExpired = 2006 // Token 已过期
	TokenMissing = 2007 // 请求头中缺少 Token

	// 数据库相关 3000-3999
	DBError       = 3000 // 数据库连接或执行错误
	DBInsertError = 3001 // 数据插入失败
	DBUpdateError = 3002 // 数据更新失败
	DBDeleteError = 3003 // 数据删除失败
	DBQueryError  = 3004 // 数据查询失败

	// 缓存相关 4000-4999
	CacheError    = 4000 // 缓存操作错误
	CacheSetError = 4001 // 缓存设置失败
	CacheGetError = 4002 // 缓存获取失败
	CacheDelError = 4003 // 缓存删除失败
)

// 错误消息
var errorMsg = map[int]string{
	Success:      "成功",
	BadRequest:   "请求参数错误",
	Unauthorized: "未授权",
	Forbidden:    "禁止访问",
	NotFound:     "资源不存在",
	ServerError:  "服务器内部错误",

	ParamInvalid:     "参数无效",
	ParamMissing:     "参数缺失",
	ParamFormatError: "参数格式错误",

	TokenInvalid: "Token 无效",
	TokenExpired: "Token 已过期",
	TokenMissing: "Token 缺失",

	DBError:       "数据库错误",
	DBInsertError: "数据插入失败",
	DBUpdateError: "数据更新失败",
	DBDeleteError: "数据删除失败",
	DBQueryError:  "数据查询失败",

	CacheError:    "缓存错误",
	CacheSetError: "缓存设置失败",
	CacheGetError: "缓存获取失败",
	CacheDelError: "缓存删除失败",
}

// GetMessage 获取错误消息
func GetMessage(code int) string {
	if msg, ok := errorMsg[code]; ok {
		return msg
	}
	return "未知错误"
}

// NewError 创建错误
func NewError(code int) *Error {
	return &Error{
		Code:    code,
		Message: GetMessage(code),
	}
}

// NewErrorWithMessage 创建错误（自定义消息）
func NewErrorWithMessage(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Error 错误结构
type Error struct {
	Code    int
	Message string
}

// Error 实现 error 接口
func (e *Error) Error() string {
	return e.Message
}
