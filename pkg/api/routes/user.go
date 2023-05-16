package routes

import (
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/api/handler"
	"github.com/gin-gonic/gin"
)

func UserRoutes(
	api *gin.RouterGroup,
	userHandler *handler.UserHandler,
) {
	// user routes that don't require authentication
	api.POST("/signup", userHandler.UserSignUp)
	api.POST("/signup/otp/verify", userHandler.SignupOtpVerify)
	api.POST("/login/email", userHandler.UserLoginByEmail)

}
