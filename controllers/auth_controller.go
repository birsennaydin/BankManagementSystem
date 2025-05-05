// controllers/auth_controller.go
package controllers

import (
	"net/http"

	"github.com/birsennaydin/BankManagementSystem/database"
	"github.com/birsennaydin/BankManagementSystem/models"
	"github.com/birsennaydin/BankManagementSystem/services"
	"github.com/birsennaydin/BankManagementSystem/utils"
	"github.com/golang-jwt/jwt"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := services.RegisterUser(user); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	accessToken, refreshToken, err := services.LoginUser(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": accessToken, "refresh_token": refreshToken})
}

func Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Refresh token geçerli mi?
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(req.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("my_secret_key"), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	username := claims["username"].(string)

	// Cassandra'dan refresh token kontrolü
	var dbUsername string
	err = database.Session.Query("SELECT username FROM refresh_tokens WHERE refresh_token = ?", req.RefreshToken).Scan(&dbUsername)
	if err != nil || dbUsername != username {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found"})
		return
	}

	// Yeni access token üret
	accessToken, _ := utils.GenerateJWT(username)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

func Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Correct table name: refresh_tokens
	err := database.Session.Query("DELETE FROM refresh_tokens WHERE refresh_token = ?", req.RefreshToken).Exec()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Logout failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
