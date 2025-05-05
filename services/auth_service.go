package services

import (
	"errors"
	"time"

	"github.com/birsennaydin/BankManagementSystem/database"
	"github.com/birsennaydin/BankManagementSystem/models"
	"github.com/birsennaydin/BankManagementSystem/utils"
	"github.com/gocql/gocql"
)

// Kullanıcı kayıt işlemi
func RegisterUser(user models.User) error {
	// Kullanıcı var mı kontrol et
	var existing string
	err := database.Session.Query("SELECT username FROM users WHERE username = ?", user.Username).Scan(&existing)
	if err == nil {
		return errors.New("user already exists")
	} else if err != gocql.ErrNotFound {
		return err
	}

	// Şifreyi hashle
	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	// Cassandra'ya ekle
	if err := database.Session.Query("INSERT INTO users (username, password) VALUES (?, ?)",
		user.Username, hashed).Exec(); err != nil {
		return err
	}

	return nil
}

func LoginUser(user models.User) (string, string, error) {
	var storedHash string

	// Fetch the hashed password for the given username from Cassandra
	err := database.Session.Query("SELECT password FROM users WHERE username = ?", user.Username).Scan(&storedHash)
	if err != nil {
		return "", "", errors.New("user not found")
	}

	// Verify that the provided password matches the stored hashed password
	if !utils.CheckPasswordHash(user.Password, storedHash) {
		return "", "", errors.New("invalid credentials")
	}

	// Generate access and refresh tokens
	accessToken, _ := utils.GenerateJWT(user.Username)
	refreshToken, _ := utils.GenerateRefreshToken(user.Username)

	// Calculate the expiration time for the refresh token (optional logging purpose)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	// Store the refresh token in Cassandra with a TTL of 7 days (604800 seconds)
	if err := database.Session.Query(`
	INSERT INTO refresh_tokens (refresh_token, username, expires_at)
	VALUES (?, ?, ?) USING TTL 604800`,
		refreshToken, user.Username, expiresAt).Exec(); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
