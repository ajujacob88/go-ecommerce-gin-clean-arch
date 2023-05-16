package routes

import "github.com/gin-gonic/gin"

func AdminRoutes(
	api *gin.RouterGroup,
	adminHandler *handler.AdminHandler,
) {
	signUp := api.Group("/admin")
	{
		signUp.POST("/signup", adminHandler.AdminSignUp)
	}
	login := api.Group("/admin")
	{
		login.POST("/login", adminHandler.AdminLogin)
	}
}
