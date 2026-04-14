package handle

import (
	"github.com/gin-gonic/gin"
	"go-mvc/internal/common/enums"
	"go-mvc/pkg/response"
)

// TestI18nHandle 测试多语言接口
type TestI18nHandle struct{}

// NewTestI18nHandle 创建测试处理器
func NewTestI18nHandle() *TestI18nHandle {
	return &TestI18nHandle{}
}

// TestSuccess 测试成功响应
func (h *TestI18nHandle) TestSuccess(c *gin.Context) {
	response.Success(c)
}

// TestError 测试错误响应
func (h *TestI18nHandle) TestError(c *gin.Context) {
	// 测试不同的错误码
	errCode := c.Query("code")
	if errCode == "" {
		errCode = enums.ErrSystemError
	}
	response.Error(c, errCode)
}

// TestData 测试数据响应
func (h *TestI18nHandle) TestData(c *gin.Context) {
	data := gin.H{
		"id":   1,
		"name": "测试数据",
	}
	response.SuccessWithData(c, data)
}
