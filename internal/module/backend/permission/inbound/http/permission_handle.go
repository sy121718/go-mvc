package permissionhttp

// import (
// 	admincontract "go-mvc/internal/module/backend/admin/contract"
// 	admindto "go-mvc/internal/module/backend/admin/dto"
// 	adminenums "go-mvc/internal/module/backend/admin/enums"
// 	r "go-mvc/pkg/response"

// 	"github.com/gin-gonic/gin"
// )

// type Handle struct {
// 	as admincontract.AdminService
// }

// func NewHandle(as admincontract.AdminService) *Handle {
// 	return &Handle{
// 		as: as,
// 	}
// }

// func (h *Handle) List(c *gin.Context) {
// 	var req admindto.ListReq
// 	if err := c.ShouldBindQuery(&req); err != nil {
// 		r.ErrorWithMessage(c, 400, adminenums.MsgBadRequest+":"+err.Error())
// 		return
// 	}

// 	res, err := h.as.List(c.Request.Context(), &req)
// 	if err != nil {
// 		r.ErrorWithMessage(c, 500, err.Error())
// 		return
// 	}

// 	r.Success(c, res)
// }

// // Login 管理员登录
// func (h *Handle) Login(c *gin.Context) {
// 	var req admindto.LoginReq
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		r.ErrorWithMessage(c, 400, adminenums.MsgBadRequest+"："+err.Error())
// 		return
// 	}

// 	res, err := h.as.Login(c.Request.Context(), &req, c.ClientIP())
// 	if err != nil {
// 		r.ErrorWithMessage(c, 400, err.Error())
// 		return
// 	}

// 	c.Header("X-New-Token", res.AccessToken)
// 	r.SuccessWithMessage(c, adminenums.MsgSuccess, res)
// }

// // Profile 获取当前登录用户信息。
// func (h *Handle) Profile(c *gin.Context) {
// 	userID, exists := c.Get("user_id")
// 	if !exists {
// 		r.ErrorWithMessage(c, 401, adminenums.MsgUnauthorized)
// 		return
// 	}

// 	uid, ok := userID.(int64)
// 	if !ok {
// 		r.ErrorWithMessage(c, 500, adminenums.MsgWrongUserType)
// 		return
// 	}

// 	res, err := h.as.Profile(c.Request.Context(), uint64(uid))
// 	if err != nil {
// 		r.ErrorWithMessage(c, 500, err.Error())
// 		return
// 	}

// 	r.Success(c, res)
// }

// // 新增管理员控制器
// // h 代表指针类型的控制器服务，c是必备的上下文go没有全局共享
// func (h *Handle) Create(c *gin.Context) {
// 	//变量 接收dto规则
// 	var req admindto.CreateReq
// 	//绑定规则，进行校验
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		r.ErrorWithMessage(c, 400, adminenums.MsgBadRequest+"："+err.Error())
// 		return
// 	}
// 	//h是控制器服务，as是adminService服务，Create是具体方法，然后传递具体参数
// 	res, err := h.as.Create(c.Request.Context(), &req)
// 	// 捕获错误
// 	if err != nil {
// 		r.ErrorWithMessage(c, 500, err.Error())
// 		return
// 	}
// 	//最终返回给用户信息
// 	r.Success(c, res)

// }

// func (h *Handle) Edit(c *gin.Context) {

// 	var req admindto.EditReq
// 	if err := c.ShouldBindJSON(&req); err != nil {

// 		r.ErrorWithMessage(c, 400, adminenums.MsgBadRequest+":"+err.Error())
// 		return
// 	}
// 	res, err := h.as.Edit(c.Request.Context(), &req)
// 	if err != nil {
// 		r.ErrorWithMessage(c, 500, err.Error())
// 		return
// 	}
// 	r.Success(c, res)
// }

// func (h *Handle) Detail(c *gin.Context) {
// 	var req admindto.DetailReq
// 	if err := c.ShouldBindQuery(&req); err != nil {
// 		r.ErrorWithMessage(c, 400, adminenums.MsgBadRequest+"："+err.Error())
// 		return
// 	}

// 	res, err := h.as.Detail(c.Request.Context(), &req)
// 	if err != nil {
// 		r.ErrorWithMessage(c, 500, err.Error())
// 		return
// 	}

// 	r.Success(c, res)
// }
