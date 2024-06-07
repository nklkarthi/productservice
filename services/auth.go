package services

import (
	"fmt"
	"productservice/models"
	"productservice/utils"
)

var users = map[string]models.User{
	"admin": {Username: "admin", Password: utils.HashPassword("adminpass"), Role: "admin"},
	"user":  {Username: "user", Password: utils.HashPassword("userpass"), Role: "user"},
}

func Authenticate(username, password string) (string, error) {
	user, exists := users[username]
	if !exists || !utils.CheckPasswordHash(password, user.Password) {
		return "", fmt.Errorf("Invalid credentials")
	}

	token, err := utils.GenerateJWT(user.Username, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
