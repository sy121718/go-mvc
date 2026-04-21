package middleware

import (
	"net/http"
	"strings"

	"go-mvc/pkg/enums"
	"go-mvc/pkg/response"

	"github.com/gin-gonic/gin"
)

// RequestBodyLimitMiddleware 限制请求体大小。
//
// 规则：
// - 普通请求使用 requestBodyLimit
// - 上传相关请求使用 uploadBodyLimit
// - 上传请求通过 URL 中包含 "/upload" 或 Content-Type 为 multipart/form-data 判断
func RequestBodyLimitMiddleware(requestBodyLimit int64, uploadBodyLimit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := resolveRequestBodyLimit(c, requestBodyLimit, uploadBodyLimit)
		if limit <= 0 {
			c.Next()
			return
		}

		if c.Request.ContentLength > limit {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, response.Response{
				Code:    enums.ErrRequestEntityTooLarge,
				Message: "请求体过大",
			})
			return
		}

		c.Next()
	}
}

func resolveRequestBodyLimit(c *gin.Context, requestBodyLimit int64, uploadBodyLimit int64) int64 {
	if c == nil || c.Request == nil {
		return requestBodyLimit
	}

	contentType := strings.ToLower(strings.TrimSpace(c.GetHeader("Content-Type")))
	path := strings.ToLower(strings.TrimSpace(c.Request.URL.Path))
	if strings.Contains(path, "/upload") || strings.HasPrefix(contentType, "multipart/form-data") {
		return uploadBodyLimit
	}
	return requestBodyLimit
}
