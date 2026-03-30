package handle

import (
	"github.com/gin-gonic/gin"
	r "go-mvc/pkg/response"
)

type TestQuery struct {
	Name string `form:"name" binding:"required,numeric"` // 👈 核心校验
}

func TestGet(c *gin.Context) {

	var req TestQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		r.Error(c, 400, err.Error())
		return
	}

	r.SuccessWithMessage(c, "555", c.Query("name"))
}
