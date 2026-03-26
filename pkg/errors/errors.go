// Package errors /*
package errors

// 错误码定义
const (
	// 通用错误码
	Success      = 200
	BadRequest   = 400
	Unauthorized = 401
	Forbidden    = 403
	NotFound     = 404
	ServerError  = 500

	// 参数相关 1000-1999
	ParamInvalid     = 1000
	ParamMissing     = 1001
	ParamFormatError = 1002

	// 用户相关 2000-2999
	UserNotFound      = 2000
	UserPasswordError = 2001
	UserDisabled      = 2002
	UserAlreadyExists = 2003
	UserNotLogin      = 2004
	TokenInvalid      = 2005
	TokenExpired      = 2006
	TokenMissing      = 2007

	// 数据库相关 3000-3999
	DBError       = 3000
	DBInsertError = 3001
	DBUpdateError = 3002
	DBDeleteError = 3003
	DBQueryError  = 3004

	// 缓存相关 4000-4999
	CacheError    = 4000
	CacheSetError = 4001
	CacheGetError = 4002
	CacheDelError = 4003

	// 业务相关 5000-5999
	BusinessError = 5000
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

	UserNotFound:      "用户不存在",
	UserPasswordError: "密码错误",
	UserDisabled:      "用户已被禁用",
	UserAlreadyExists: "用户已存在",
	UserNotLogin:      "用户未登录",
	TokenInvalid:      "Token 无效",
	TokenExpired:      "Token 已过期",
	TokenMissing:      "Token 缺失",

	DBError:       "数据库错误",
	DBInsertError: "数据插入失败",
	DBUpdateError: "数据更新失败",
	DBDeleteError: "数据删除失败",
	DBQueryError:  "数据查询失败",

	CacheError:    "缓存错误",
	CacheSetError: "缓存设置失败",
	CacheGetError: "缓存获取失败",
	CacheDelError: "缓存删除失败",

	BusinessError: "业务错误",
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
