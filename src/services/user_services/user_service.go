package user_services

import (
	"fmt"
	"subscriber-topic-stars/src/helpers"
	"subscriber-topic-stars/src/repositories/user_repositories"

	"github.com/gin-gonic/gin"
)

type UserServiceInterface interface {
	GetUser(msg map[string]interface{}) (interface{}, error)
}

type UserService struct {
	userRepo user_repositories.UserRepositoryInterface
}

func NewUserService(repo user_repositories.UserRepositoryInterface) *UserService {
	return &UserService{userRepo: repo}
}

func (u *UserService) GetUser(msg map[string]interface{}) (interface{}, error) {
	tokenRaw, ok := msg["token"]
	if !ok {
		return map[string]interface{}{
			"success": false,
			"message": "Missing token",
		}, nil
	}

	tokenStr := fmt.Sprintf("%v", tokenRaw)

	// Validasi token
	claims, err := helpers.ParseAndValidateToken(tokenStr)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"message": fmt.Sprintf("Invalid token: %v", err),
		}, nil
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return map[string]interface{}{
			"success": false,
			"message": "Invalid user_id in token",
		}, nil
	}
	userID := uint64(userIDFloat)

	user, err := u.userRepo.FindUserById(userID)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"message": "User not found",
		}, nil
	}

	return gin.H{
		"uuid":     user.UUID,
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
	}, nil
}
