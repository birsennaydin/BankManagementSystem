// routes/profile_routes.go
package routes

import (
	"github.com/birsennaydin/BankManagementSystem/controllers"
	"github.com/birsennaydin/BankManagementSystem/middleware"
	"github.com/gin-gonic/gin"
)

func ProfileRoutes(router *gin.Engine) {
	protected := router.Group("/profile")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("", controllers.Profile)
}
