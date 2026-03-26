// Package crypto /*
package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

var (
	secretKey          = "your-secret-key" // 签名密钥
	timestampTolerance = 300 * time.Second // 时间戳容差（5分钟）
)

// SetSecretKey 设置签名密钥
func SetSecretKey(key string) {
	secretKey = key
}

// GenerateSignature 生成签名
// params: 请求参数
// timestamp: 时间戳（秒）
func GenerateSignature(params map[string]interface{}, timestamp int64) string {
	// 1. 参数排序
	sortedParams := sortParams(params)

	// 2. 拼接字符串
	signStr := fmt.Sprintf("%s&timestamp=%d&key=%s", sortedParams, timestamp, secretKey)

	// 3. HMAC-SHA256 签名
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(signStr))
	return hex.EncodeToString(h.Sum(nil))
}

// VerifySignature 验证签名
// params: 请求参数
// timestamp: 时间戳（秒）
// signature: 前端传来的签名
func VerifySignature(params map[string]interface{}, timestamp int64, signature string) error {
	// 1. 验证时间戳（防重放）
	if err := verifyTimestamp(timestamp); err != nil {
		return err
	}

	// 2. 生成签名
	expectedSign := GenerateSignature(params, timestamp)

	// 3. 对比签名
	if !hmac.Equal([]byte(expectedSign), []byte(signature)) {
		return fmt.Errorf("签名验证失败")
	}

	return nil
}

// verifyTimestamp 验证时间戳
func verifyTimestamp(timestamp int64) error {
	now := time.Now().Unix()
	diff := now - timestamp

	// 时间差超过容差
	if diff > int64(timestampTolerance.Seconds()) || diff < -int64(timestampTolerance.Seconds()) {
		return fmt.Errorf("时间戳已过期")
	}

	return nil
}

// sortParams 参数排序
func sortParams(params map[string]interface{}) string {
	// 提取所有 key
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "sign" || k == "signature" {
			continue // 跳过签名字段
		}
		keys = append(keys, k)
	}

	// 排序
	sort.Strings(keys)

	// 拼接
	var builder strings.Builder
	for i, k := range keys {
		if i > 0 {
			builder.WriteString("&")
		}
		builder.WriteString(fmt.Sprintf("%s=%v", k, params[k]))
	}

	return builder.String()
}

// GenerateSignatureWithString 字符串签名（简化版）
func GenerateSignatureWithString(data string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// VerifySignatureWithString 字符串签名验证（简化版）
func VerifySignatureWithString(data, signature string) bool {
	expected := GenerateSignatureWithString(data)
	return hmac.Equal([]byte(expected), []byte(signature))
}
