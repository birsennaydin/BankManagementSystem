// services/auth_service.go
package services

import (
	"errors"

	"github.com/birsennaydin/BankManagementSystem/models"
	"github.com/birsennaydin/BankManagementSystem/utils"
)

var users = map[string]string{} // username: hashedPassword

func RegisterUser(user models.User) error {
	if _, exists := users[user.Username]; exists {
		return errors.New("user already exists")
	}
	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	users[user.Username] = hashed
	return nil
}

func LoginUser(user models.User) (string, error) {
	stored, exists := users[user.Username]
	if !exists || !utils.CheckPasswordHash(user.Password, stored) {
		return "", errors.New("invalid credentials")
	}
	return utils.GenerateJWT(user.Username)
}
