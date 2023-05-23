package routes

import (
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/api/handler"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/api/middleware"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(
	api *gin.RouterGroup,
	adminHandler *handler.AdminHandler,
	productHandler *handler.ProductHandler,

) {

	api.POST("/login", adminHandler.AdminLogin)

	api.Use(middleware.AdminAuth) //Middleware functions in Gin are handlers that are executed before reaching the final handler for a particular route. They provide a way to perform common tasks, such as authentication, logging, or data preprocessing, for multiple routes or groups of routes.the AdminAuth middleware function is used to enforce authentication for certain routes under the /admins group.The purpose of using the AdminAuth middleware here is to ensure that only authenticated administrators can access the routes under the /admins group.
	{
		api.GET("/logout", adminHandler.AdminLogout)

		//user management
		userManagement := api.Group("/users")
		{
			userManagement.GET("/", adminHandler.ListAllUsers)
			userManagement.GET("/:id", adminHandler.FindUserByID)
			userManagement.PUT("/:id/block", adminHandler.BlockUser)
			userManagement.PUT("/unblock/:id", adminHandler.UnblockUser)
		}

		//admin management
		adminManagement := api.Group("/admins")
		{
			adminManagement.POST("/", adminHandler.CreateAdmin)
		}

		//category management routes
		categoryRoutes := api.Group("/categories")
		{
			categoryRoutes.POST("/", productHandler.CreateCategory)
		}

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
