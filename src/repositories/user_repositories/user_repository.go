package user_repositories

import (
	"fmt"
	"subscriber-topic-stars/src/entities/users"
	"subscriber-topic-stars/src/helpers"
)

type UserRepositoryInterface interface {
	FindUserById(UserID uint64) (*users.User, error)
}

type userRepository struct{}

func NewUserRepository() userRepository {
	return userRepository{}
}

func (u userRepository) FindUserById(UserID uint64) (*users.User, error) {
	var user users.User
	if err := helpers.FindOneByField(&user, "id", UserID); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}
