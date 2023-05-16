package routes

import (
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/api/handler"
	"github.com/gin-gonic/gin"
)

func UserRoutes(
	api *gin.RouterGroup,
	userHandler *handler.UserHandler,
) {

	// User routes that don't require authentication
	api.POST("/signup", userHandler.UserSignUp)
	api.POST("/signup/otp/verify", userHandler.SignupOtpVerify)
	api.POST("/login/email", userHandler.UserLoginByEmail)

	// User routes that require authentication
	// api.Use(middleware.UserAuth)
	// {
	// 	api.GET("/profile", userHandler.UserProfile)
	// 	api.GET("/logout", userHandler.UserLogout)

	// }

	//sets up a route group for the "/user" endpoint but i think no need to group the user
	// signup := api.Group("/user")
	// {
	// 	signup.POST("/signup", userHandler.UserSignUp)
	// 	signup.POST("/signup/otp/verify", userHandler.SignupOtpVerify)
	// }

	// login := api.Group("/user")
	// {
	// 	login.POST("/login/email", userHandler.UserLoginByEmail)
	// }
}
