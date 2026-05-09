// Package admindto 管理员模块数据传输对象
//
// 该包定义了管理员模块在接口层与业务层之间传递的请求和响应数据结构。
// 每个结构体对应一个业务操作的输入输出格式，包含参数校验规则（binding/validate tag），
// 确保进入 service 层的数据是合法、完整的。
package admindto

import adminmodel "go-mvc/internal/module/backend/admin/model"

// ListResp 管理员列表查询响应
//
// 返回符合查询条件的总记录数和当前页的数据列表。
type ListResp struct {
	Total int64                    `json:"total"` // 符合条件的总记录数，用于前端分页组件
	List  []adminmodel.AdminEntity `json:"list"`  // 当前页的管理员数据列表
}

// CreateResp 管理员新增响应
type CreateResp struct {
	ID       uint64 `json:"id"` // 新增的管理员 ID，表示添加成功
	Username string `json:"username"`
}
