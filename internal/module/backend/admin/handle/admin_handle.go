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

// 新增管理员控制器
// h 代表指针类型的控制器服务，c是必备的上下文go没有全局共享
func (h *Handle) Create(c *gin.Context) {
	//变量 接收dto规则
	var req admindto.CreateReq
	//绑定规则，进行校验
	if err := c.ShouldBindJSON(&req); err != nil {
		r.ErrorWithMessage(c, 400, "请求参数错误："+err.Error())
		return
	}
	//h是控制器服务，as是adminService服务，Create是具体方法，然后传递具体参数
	res, err := h.as.Create(c.Request.Context(), &req)
	// 捕获错误
	if err != nil {
		r.ErrorWithMessage(c, 500, err.Error())
		return
	}
	//最终返回给用户信息
	r.Success(c, res)

}
