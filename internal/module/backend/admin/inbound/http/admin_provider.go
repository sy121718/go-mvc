package adminhttp

import (
	adminmodel "go-mvc/internal/module/backend/admin/model"
	adminservice "go-mvc/internal/module/backend/admin/service"
	"go-mvc/pkg/database"
)

func newAdminHandle() (*Handle, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, err
	}

	model := adminmodel.NewAdminModel(db)
	service := adminservice.NewService(adminservice.Deps{
		AdminModel: model,
	})
	handle := NewHandle(Deps{
		AdminService: service,
	})
	return handle, nil
}
