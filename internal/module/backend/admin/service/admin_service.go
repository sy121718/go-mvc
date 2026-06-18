package adminservice

import (
	admincontract "go-mvc/internal/module/backend/admin/contract"
	adminmodel "go-mvc/internal/module/backend/admin/model"
)

// 编译期断言：确保 *Service 实现了 AdminService 接口
var _ admincontract.AdminService = (*Service)(nil)

// Service 定义了 Admin 模块的业务逻辑
type Service struct {
	am *adminmodel.AdminModel
}

// NewService 创建一个 Service 实例对象。
func NewService(am *adminmodel.AdminModel) *Service {
	return &Service{
		am: am,
	}
}
