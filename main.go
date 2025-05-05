package main

import (
	"github.com/birsennaydin/BankManagementSystem/database"
	"github.com/birsennaydin/BankManagementSystem/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	database.InitCassandra()

	router := gin.Default()

	routes.AuthRoutes(router)
	routes.ProfileRoutes(router)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.Run(":8080")
}
