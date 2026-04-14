package middleware

import (
	"github.com/gin-gonic/gin"
	"go-mvc/pkg/crypto"
	"go-mvc/pkg/errors"
	"go-mvc/pkg/response"
	"strconv"
	"time"
)

// SignatureMiddleware 签名验证中间件
func SignatureMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取时间戳
		timestampStr := c.GetHeader("X-Timestamp")
		if timestampStr == "" {
			timestampStr = c.Query("timestamp")
		}

		if timestampStr == "" {
			response.ErrorCode(c, errors.ParamMissing)
			c.Abort()
			return
		}

		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			response.ErrorCode(c, errors.ParamFormatError)
			c.Abort()
			return
		}

		// 2. 验证时间戳（防重放）
		now := time.Now().Unix()
		if now-timestamp > 300 || timestamp-now > 300 { // 5分钟容差
			response.Error(c, errors.BadRequest, "请求已过期")
			c.Abort()
			return
		}

		// 3. 获取签名
		signature := c.GetHeader("X-Signature")
		if signature == "" {
			signature = c.Query("signature")
		}

		if signature == "" {
			response.ErrorCode(c, errors.ParamMissing)
			c.Abort()
			return
		}

		// 4. 收集参数
		params := make(map[string]interface{})

		// Query 参数
		query := c.Request.URL.Query()
		for k, v := range query {
			if k != "signature" && k != "timestamp" {
				params[k] = v[0]
			}
		}

		// POST 参数（如果是 JSON）
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			var body map[string]interface{}
			if err := c.ShouldBindJSON(&body); err == nil {
				for k, v := range body {
					if k != "signature" && k != "timestamp" {
						params[k] = v
					}
				}
			}
		}

		// 5. 验证签名
		if err := crypto.VerifySignature(params, timestamp, signature); err != nil {
			response.Error(c, errors.BadRequest, "签名验证失败")
			c.Abort()
			return
		}

		c.Next()
	}
}
