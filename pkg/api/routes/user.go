package routes

import (
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/api/handler"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/api/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(
	api *gin.RouterGroup,
	userHandler *handler.UserHandler,
	productHandler *handler.ProductHandler,
	cartHandler *handler.CartHandler,
	orderHandler *handler.OrderHandler,
	paymentHandler *handler.PaymentHandler,
) {

	// User routes that don't require authentication
	//sets up a route group for the "/user" endpoint
	signup := api.Group("/user")
	{
		signup.POST("/signup", userHandler.UserSignUp)
		signup.POST("/signup/otp/verify", userHandler.SignupOtpVerify)
	}

	login := api.Group("/user")
	{
		login.POST("/login/email", userHandler.UserLoginByEmail)
	}

	//User routes that require authentication

	home := api.Group("/user")
	{
		//AuthorizationMiddleware as middleware to perform authorization checks for users accessing the "/user" endpoint.
		home.Use(middleware.UserAuth)
		{
			home.GET("/home", userHandler.Homehandler)
			home.GET("/logout", userHandler.LogoutHandler)
			home.GET("/products", productHandler.ListAllProducts)
			home.GET("/products/:id", productHandler.FindProductByID)
			home.POST("/cart/add/:product_details_id", cartHandler.AddToCart)
			home.DELETE("/cart/remove/:product_details_id", cartHandler.RemoveFromCart)
			home.GET("/cart", cartHandler.ViewCart)
			home.POST("/addresses", userHandler.AddAddress)
			home.PATCH("/addresses/edit/:address_id", userHandler.UpdateAddress)
			home.DELETE("/addresses/:address_id", userHandler.DeleteAddress)
			home.GET("/addresses", userHandler.ListAddress)
			home.POST("/cart/placeorder", orderHandler.PlaceOrderFromCart)
			home.GET("/payments/razorpay/:order_id", paymentHandler.RazorpayCheckout)
			home.POST("/payments/success", paymentHandler.RazorpayVerify)
			//cart routes
			// cart := api.Group("/cart")
			// {
			// 	cart.POST("/add/:product_details_id", cartHandler.AddToCart)
			// 	cart.POST("/remove/:product_details_id", cartHandler.RemoveFromCart)

			// }
		}
	}

}

// api.POST("/signup", userHandler.UserSignUp)
// api.POST("/signup/otp/verify", userHandler.SignupOtpVerify)
// api.POST("/login/email", userHandler.UserLoginByEmail)

// api.Use(middleware.AuthorizationMiddleware("user"))
// {
// 	api.GET("/home", userHandler.UserProfile)
// 	api.GET("/logout", userHandler.UserLogout)

// }
