package permissionservice

import (
	permissioncontract "go-mvc/internal/module/backend/permission/contract"
	permissionmodel "go-mvc/internal/module/backend/permission/model"
)

var _ permissioncontract.PermissionService = (*Service)(nil)

type Service struct {
	pm *permissionmodel.PermissionModel
}

func NewService(pm *permissionmodel.PermissionModel) *Service {

	return &Service{
		pm: pm,
	}
}
