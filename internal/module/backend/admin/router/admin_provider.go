package adminrouter

import (
	adminhandle "go-mvc/internal/module/backend/admin/handle"
	adminmodel "go-mvc/internal/module/backend/admin/model"
	adminservice "go-mvc/internal/module/backend/admin/service"
	"go-mvc/pkg/database"
)

func newAdminHandle() (*adminhandle.Handle, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, err
	}

	model := adminmodel.NewAdminModel(db)
	service := adminservice.NewService(adminservice.Deps{
		AdminModel: model,
	})
	handle := adminhandle.NewHandle(adminhandle.Deps{
		AdminService: service,
	})
	return handle, nil
}
