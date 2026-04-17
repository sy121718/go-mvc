package response

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"go-mvc/pkg/i18n"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构。
type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type langCandidate struct {
	lang string
	q    float64
	idx  int
}

func getLang(c *gin.Context) string {
	if headerLang := parseAcceptLanguage(c.GetHeader("Accept-Language")); headerLang != "" {
		return headerLang
	}

	if queryLang := normalizeLang(c.Query("lang")); queryLang != "" {
		return queryLang
	}

	return i18n.GetDefaultLang()
}

func parseAcceptLanguage(header string) string {
	header = strings.TrimSpace(header)
	if header == "" {
		return ""
	}

	parts := strings.Split(header, ",")
	candidates := make([]langCandidate, 0, len(parts))
	for i, part := range parts {
		section := strings.TrimSpace(part)
		if section == "" {
			continue
		}

		items := strings.Split(section, ";")
		lang := normalizeLang(items[0])
		if lang == "" {
			continue
		}

		q := 1.0
		for _, item := range items[1:] {
			item = strings.TrimSpace(item)
			if !strings.HasPrefix(strings.ToLower(item), "q=") {
				continue
			}
			value, err := strconv.ParseFloat(strings.TrimSpace(item[2:]), 64)
			if err == nil {
				q = value
			}
		}

		if q <= 0 {
			continue
		}
		candidates = append(candidates, langCandidate{lang: lang, q: q, idx: i})
	}

	if len(candidates) == 0 {
		return ""
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].q == candidates[j].q {
			return candidates[i].idx < candidates[j].idx
		}
		return candidates[i].q > candidates[j].q
	})

	return candidates[0].lang
}

func normalizeLang(lang string) string {
	lang = strings.TrimSpace(strings.ReplaceAll(lang, "_", "-"))
	if lang == "" {
		return ""
	}

	parts := strings.Split(lang, "-")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
		if parts[i] == "" {
			return ""
		}
	}

	parts[0] = strings.ToLower(parts[0])
	for i := 1; i < len(parts); i++ {
		parts[i] = strings.ToUpper(parts[i])
	}
	return strings.Join(parts, "-")
}

// Success 成功响应。
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

// SuccessWithMessage 成功响应（自定义消息码）。
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

// Error 错误响应（自动获取多语言消息和 HTTP 状态码）。
func Error(c *gin.Context, errCode string) {
	lang := getLang(c)
	result := i18n.Get(errCode, lang)

	c.JSON(result.HttpCode, Response{
		Code:    errCode,
		Message: result.Value,
	})
}

// ErrorWithMessage 错误响应（自定义消息）。
func ErrorWithMessage(c *gin.Context, errCode string, message string) {
	lang := getLang(c)
	result := i18n.Get(errCode, lang)

	c.JSON(result.HttpCode, Response{
		Code:    errCode,
		Message: message,
	})
}

// ParamError 参数错误。
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

// NotFound 404 响应。
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

// SuccessWithData 成功响应（兼容旧代码）。
func SuccessWithData(c *gin.Context, data interface{}) {
	Success(c, data)
}
