package http

import (
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/api/handler"
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/api/routes"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(userHandler *handler.UserHandler,
	adminHandler *handler.AdminHandler,
	productHandler *handler.ProductHandler,
	cartHandler *handler.CartHandler,
	orderHandler *handler.OrderHandler,
	paymentHandler *handler.PaymentHandler,
	couponHandler *handler.CouponHandler,

) *ServerHTTP {

	engine := gin.New()

	// Use logger from Gin
	engine.Use(gin.Logger())

	//add swagger - Swagger docs
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//setup routes
	routes.UserRoutes(engine.Group("/"), userHandler, productHandler, cartHandler, orderHandler, paymentHandler, couponHandler)
	routes.AdminRoutes(engine.Group("/admin"), adminHandler, productHandler, couponHandler, orderHandler)

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	sh.engine.LoadHTMLGlob("views/*.html") //for loading the html page of razor pay
	sh.engine.Run(":3000")
}

/* no need

// Swagger docs
//engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

// Request JWT
engine.POST("/login", middleware.LoginHandler)

// Auth middleware
api := engine.Group("/api", middleware.AuthorizationMiddleware)

api.GET("users", userHandler.FindAll)
api.GET("users/:id", userHandler.FindByID)
api.POST("users", userHandler.Save)
api.DELETE("users/:id", userHandler.Delete)

*/
