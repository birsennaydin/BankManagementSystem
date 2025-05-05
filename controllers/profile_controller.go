// controllers/profile_controller.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Profile(c *gin.Context) {
	username := c.GetString("username")

	c.JSON(http.StatusOK, gin.H{
		"message":  "Welcome to your profile!",
		"username": username,
	})
}
