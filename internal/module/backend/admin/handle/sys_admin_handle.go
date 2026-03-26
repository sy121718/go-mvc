package handle

import (
	"github.com/gin-gonic/gin"
	r "go-mvc/pkg/response"
	"strconv"
)

/*
控制器 - AdminHandle
===========================================
PHP 对比：
- Laravel: Controller 类中的方法
- Go: 普通函数，接收 *gin.Context 参数

知识点：
1. *gin.Context - 包含请求和响应的所有信息
2. c.JSON() - 返回 JSON 响应
3. c.ShouldBind() - 绑定请求参数
*/

// Test 测试接口
func Test(c *gin.Context) {

	var a = int(1)
	var b int
	b = 2
	s := "1"
	num, err := strconv.Atoi(s)
	if err != nil {

	}
	r.SuccessWithMessage(c, "555", a+b+num)
}

func Test2(c *gin.Context) {

}
