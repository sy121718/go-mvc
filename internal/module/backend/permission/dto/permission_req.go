// Package permissiondto 管理员模块数据传输对象
//
// 该包定义了管理员模块在接口层与业务层之间传递的请求和响应数据结构。
// 每个结构体对应一个业务操作的输入输出格式，包含参数校验规则（binding/validate tag），
// 确保进入 service 层的数据是合法、完整的。
package permissiondto

// PageReq 分页参数，所有列表接口共用。
// 前端不传时默认 page=1, limit=10。
type PageReq struct {
	Page  int `form:"page" json:"page" binding:"omitempty,gte=1" validate:"omitempty,gte=1"`
	Limit int `form:"limit" json:"limit" binding:"omitempty,gte=1,lte=100" validate:"omitempty,gte=1,lte=100"`
}

func (r *PageReq) GetPage() int {
	if r.Page < 1 {
		return 1
	}
	return r.Page
}

func (r *PageReq) GetLimit() int {
	if r.Limit < 1 {
		return 10
	}
	if r.Limit > 100 {
		return 100
	}
	return r.Limit
}

// ListReq 权限点列表查询请求参数
//
// 支持按模块筛选，可选分页。
type ListReq struct {
	PageReq
	Module string `form:"module" json:"module" binding:"omitempty,max=50" validate:"omitempty,max=50"`
}
