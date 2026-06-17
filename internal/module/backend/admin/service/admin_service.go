package adminservice

import (
	adminmodel "go-mvc/internal/module/backend/admin/model"
)

// Service 定义了 Admin 模块的业务逻辑
type Service struct {
	am *adminmodel.AdminModel
}

// Deps 定义了 Service 依赖的组件
type Deps struct {
	AdminModel *adminmodel.AdminModel
}

// NewService 创建一个 Service 实例对象
func NewService(deps Deps) *Service {
	return &Service{
		am: deps.AdminModel,
	}
}
