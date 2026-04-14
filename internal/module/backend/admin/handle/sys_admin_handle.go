package handle

import (
	"github.com/gin-gonic/gin"
	m "go-mvc/internal/module/backend/admin/model"
	r "go-mvc/pkg/response"
	v "go-mvc/pkg/validate"
)

type TestQuery struct {
	adminModel m.Admin
}

func TestGet(c *gin.Context) {

	var req TestQuery
	// ShouldBindQuery 将 URL query 参数按 form tag 绑定到结构体，
	// 并根据 binding tag 进行校验，失败时返回 error
	if err := c.ShouldBindQuery(&req); err != nil {
		r.ParamError(c, v.Msg(err))
		return
	}

	//if err := c.ShouldBindQuery(&); err != nil {
	//	r.ParamError(c, v.Msg(err))
	//}

	r.SuccessWithMessage(c, "555", c.Query("name"))
}
