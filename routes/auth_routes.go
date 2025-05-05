// routes/auth_routes.go
package routes

import (
	"github.com/birsennaydin/BankManagementSystem/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)
}
