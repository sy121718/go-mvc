package permissionhttp

// func newHandle() (*Handle, error) {
// 	db, err := database.GetDB()
// 	if err != nil {
// 		return nil, err
// 	}
// 	service := permissionservice.NewService(permissionodel.NewAdminModel(db))
// 	return NewHandle(service), nil
// }

// func SetupAdminRoutes(rg *gin.RouterGroup) {
// 	if rg == nil {
// 		return
// 	}

// 	handle, err := newHandle()
// 	if err != nil {
// 		return
// 	}

// 	admin := rg.Group("/permission")
// 	admin.POST("/login", handle.Login)

// 	auth := admin.Group("").Use(builtin.JWTAuthMiddleware())
// 	{

// 	}
// }
