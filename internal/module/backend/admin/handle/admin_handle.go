package adminhandle

import (
	admindto "go-mvc/internal/module/backend/admin/dto"
	adminservice "go-mvc/internal/module/backend/admin/service"
	r "go-mvc/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handle struct {
	as *adminservice.Service
}

type Deps struct {
	AdminService *adminservice.Service
}

func NewHandle(deps Deps) *Handle {
	return &Handle{
		as: deps.AdminService,
	}
}

func (h *Handle) List(c *gin.Context) {
	var req admindto.ListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		r.ErrorWithMessage(c, 400, "请求参数错误:"+err.Error())
		return
	}

	res, err := h.as.List(c.Request.Context(), &req)
	if err != nil {
		r.ErrorWithMessage(c, 500, err.Error())
		return
	}

	r.Success(c, res)
}
