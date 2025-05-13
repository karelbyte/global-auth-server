package routes

import (
	"global-auth-server/controllers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/", controllers.Home)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Group of routes with prefix /API
	api := r.Group("/api")
	{
		api.POST("/auth/login", controllers.Login)
		api.POST("/auth/can-login", controllers.CanLogin)
	}
}
