package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"time"

	"go-mvc/pkg/enums"

	"go-mvc/pkg/crypto"
	"go-mvc/pkg/response"

	"github.com/gin-gonic/gin"
)

// SignatureMiddleware 签名验证中间件
func SignatureMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		timestampStr := c.GetHeader("X-Timestamp")
		if timestampStr == "" {
			timestampStr = c.Query("timestamp")
		}
		if timestampStr == "" {
			response.ParamError(c, "缺少时间戳参数")
			c.Abort()
			return
		}

		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			response.ParamError(c, "时间戳格式错误")
			c.Abort()
			return
		}

		now := time.Now().Unix()
		if now-timestamp > 300 || timestamp-now > 300 {
			response.ErrorWithMessage(c, enums.ErrInvalidParams, "请求已过期")
			c.Abort()
			return
		}

		signature := c.GetHeader("X-Signature")
		if signature == "" {
			signature = c.Query("signature")
		}
		if signature == "" {
			response.ParamError(c, "缺少签名参数")
			c.Abort()
			return
		}

		params := make(map[string]interface{})

		query := c.Request.URL.Query()
		for k, v := range query {
			if k == "signature" || k == "timestamp" || len(v) == 0 {
				continue
			}
			params[k] = v[0]
		}

		if c.Request.Method == "POST" {
			bodyParams, readErr := readBodyParams(c)
			if readErr != nil {
				response.ParamError(c, "请求体格式错误")
				c.Abort()
				return
			}
			for k, v := range bodyParams {
				if k == "signature" || k == "timestamp" {
					continue
				}
				params[k] = v
			}
		}

		if err := crypto.VerifySignature(params, timestamp, signature); err != nil {
			response.ErrorWithMessage(c, enums.ErrInvalidParams, "签名验证失败")
			c.Abort()
			return
		}

		c.Next()
	}
}

func readBodyParams(c *gin.Context) (map[string]interface{}, error) {
	if c.Request == nil || c.Request.Body == nil {
		return map[string]interface{}{}, nil
	}

	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(raw))

	if len(raw) == 0 {
		return map[string]interface{}{}, nil
	}

	var body map[string]interface{}
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, err
	}
	return body, nil
}
