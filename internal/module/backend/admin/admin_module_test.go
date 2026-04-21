package adminmodule

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterRoutesCanBeCalledFromModuleEntry(t *testing.T) {
	engine := gin.New()
	group := engine.Group("/api")

	RegisterRoutes(group)
}
