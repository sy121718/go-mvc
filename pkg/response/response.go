package response

import (
	"go-mvc/pkg/i18n"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// getLang 获取客户端语言
func getLang(c *gin.Context) string {
	// 优先从 Header 获取
	lang := c.GetHeader("Accept-Language")
	if lang != "" {
		return lang
	}
	// 其次从 Query 参数获取
	lang = c.Query("lang")
	if lang != "" {
		return lang
	}
	// 默认中文
	return "zh-CN"
}

// Success 成功响应
func Success(c *gin.Context, data ...interface{}) {
	lang := getLang(c)
	result := i18n.Get("msg_operation_success", lang)

	var responseData interface{}
	if len(data) > 0 {
		responseData = data[0]
	}

	c.JSON(http.StatusOK, Response{
		Code:    "0",
		Message: result.Value,
		Data:    responseData,
	})
}

// SuccessWithMessage 成功响应（自定义消息码）
func SuccessWithMessage(c *gin.Context, msgCode string, data ...interface{}) {
	lang := getLang(c)
	result := i18n.Get(msgCode, lang)

	var responseData interface{}
	if len(data) > 0 {
		responseData = data[0]
	}

	c.JSON(http.StatusOK, Response{
		Code:    "0",
		Message: result.Value,
		Data:    responseData,
	})
}

// Error 错误响应（自动获取多语言消息和HTTP状态码）
func Error(c *gin.Context, errCode string) {
	lang := getLang(c)
	result := i18n.Get(errCode, lang)

	c.JSON(result.HttpCode, Response{
		Code:    errCode,
		Message: result.Value,
	})
}

// ErrorWithMessage 错误响应（自定义消息）
func ErrorWithMessage(c *gin.Context, errCode string, message string) {
	lang := getLang(c)
	result := i18n.Get(errCode, lang)

	c.JSON(result.HttpCode, Response{
		Code:    errCode,
		Message: message,
	})
}

// ParamError 参数错误
func ParamError(c *gin.Context, msg ...string) {
	lang := getLang(c)
	result := i18n.Get("ErrInvalidParams", lang)

	message := result.Value
	if len(msg) > 0 {
		message = msg[0]
	}

	c.JSON(result.HttpCode, Response{
		Code:    "ErrInvalidParams",
		Message: message,
	})
}

// NotFound 404响应
func NotFound(c *gin.Context, msg ...string) {
	lang := getLang(c)
	result := i18n.Get("ErrNotFound", lang)

	message := result.Value
	if len(msg) > 0 {
		message = msg[0]
	}

	c.JSON(result.HttpCode, Response{
		Code:    "ErrNotFound",
		Message: message,
	})
}

// SuccessWithData 成功响应（兼容旧代码）
func SuccessWithData(c *gin.Context, data interface{}) {
	Success(c, data)
}
