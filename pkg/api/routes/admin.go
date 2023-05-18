package routes

import (
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/api/handler"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(
	api *gin.RouterGroup,
	adminHandler *handler.AdminHandler,
) {
	adminManagement := api.Group("/admins")
	{
		adminManagement.POST("/", adminHandler.CreateAdmin)
	}
}

// signUp := api.Group("/admin")
// 	{
// 		signUp.POST("/signup", adminHandler.CreateAdmin)
// 	}
// 	login := api.Group("/admin")
// 	{
// 		login.POST("/login", adminHandler.AdminLogin)
// 	}
