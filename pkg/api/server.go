package http

import (
	"github.com/ajujacob88/go-ecommerce-gin-clean-arch/pkg/api/middleware"
	"github.com/gin-gonic/gin"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(userHandler *handler.userHandler) {
	engine := gin.New()

	// Use logger from Gin
	engine.Use(gin.Logger())

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

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	sh.engine.Run(":3000")
}
