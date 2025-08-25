package auth_repositories

import (
	"fmt"
	"strings"
	"subscriber-topic-stars/src/dtos/auth_dtos"
	"subscriber-topic-stars/src/entities/auth"
	"subscriber-topic-stars/src/entities/users"
	"subscriber-topic-stars/src/helpers"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthRepositoryInterface interface {
	Register(req auth_dtos.RequestRegisterDto) (map[string]interface{}, error)
	FindByEmail(email string) (*users.User, error)
	FindByUsername(username string) (*users.User, error)
	CreateUser(user *users.User) error
	SaveTokens(userID uint64, accessToken string, accessExp time.Time, refreshToken string, refreshExp time.Time) error
	FindRefreshToken(token string) (*auth.RefreshToken, error)
	MarkRefreshTokenAsUsed(id uint64) error
	MarkTokenAsRevoked(tokenID uint64) error
	FindTokenByUserIDAndToken(userID uint64, tokenString string) (*auth.AccessToken, error)
	FindRefreshTokenByAccessTokenID(tokenID uint64) (*auth.RefreshToken, error)
}

type authRepository struct{}

func NewAuthRepository() authRepository {
	return authRepository{}
}

func (r authRepository) Register(req auth_dtos.RequestRegisterDto) (map[string]interface{}, error) {

	if _, err := r.FindByEmail(req.Email); err == nil {
		return nil, fmt.Errorf("email already in use %w", err)
	}

	if _, err := r.FindByUsername(req.Username); err == nil {
		return nil, fmt.Errorf("username already in use %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("could not hash password: %w", err)
	}

	newUser := users.User{
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashedPassword),
		Name:     req.Name,
	}

	if err := helpers.InsertModel(&newUser); err != nil {
		return nil, fmt.Errorf("could not insert user: %w", err)
	}

	// Return a map with user data
	response := map[string]interface{}{
		"id":       newUser.ID,
		"email":    newUser.Email,
		"username": newUser.Username,
		"name":     newUser.Name,
	}
	return response, nil
}

func (r authRepository) FindByEmail(email string) (*users.User, error) {
	var user users.User
	err := helpers.FindOneByField(&user, "email", email)
	if err != nil {
		return nil, fmt.Errorf("email not found: %w", err)

	}
	return &user, nil
}

func (r authRepository) FindByUsername(username string) (*users.User, error) {
	var user users.User
	err := helpers.FindOneByField(&user, "username", username)
	if err != nil {
		return nil, fmt.Errorf("username not found: %w", err)
	}
	return &user, nil
}

func (r authRepository) CreateUser(user *users.User) error {
	return helpers.InsertModel(user)
}

func (r authRepository) SaveTokens(userID uint64, accessToken string, accessExp time.Time, refreshToken string, refreshExp time.Time) error {
	const maxRetries = 3

	var access auth.AccessToken
	var err error

	for i := 0; i < maxRetries; i++ {
		access = auth.AccessToken{
			UserID:    userID,
			Token:     accessToken,
			ExpiresAt: accessExp,
		}
		err = helpers.InsertModel(&access)
		if err != nil && strings.Contains(err.Error(), "duplicate key value") {
			accessToken = helpers.GenerateRandomToken() // regenerate
			continue
		}
		break
	}
	if err != nil {
		return fmt.Errorf("failed insert access token: %w", err)
	}

	var refresh auth.RefreshToken
	for i := 0; i < maxRetries; i++ {
		refresh = auth.RefreshToken{
			UserID:        userID,
			AccessTokenID: access.ID,
			Token:         refreshToken,
			ExpiresAt:     refreshExp,
		}
		err = helpers.InsertModel(&refresh)
		if err != nil && strings.Contains(err.Error(), "duplicate key value") {
			refreshToken = helpers.GenerateRandomToken()
			continue
		}
		break
	}
	if err != nil {
		return fmt.Errorf("failed insert refresh token: %w", err)
	}

	return nil
}

func (r authRepository) FindRefreshToken(token string) (*auth.RefreshToken, error) {
	var refresh auth.RefreshToken
	if err := helpers.FindOneByField(&refresh, "token", token); err != nil {
		return nil, fmt.Errorf("token not found: %w", err)
	}
	return &refresh, nil
}

func (r authRepository) FindRefreshTokenByAccessTokenID(accessTokenID uint64) (*auth.RefreshToken, error) {
	var refresh auth.RefreshToken
	err := helpers.FindOneByField(&refresh, "access_token_id", accessTokenID)
	if err != nil {
		return nil, fmt.Errorf("refresh token not found: %w", err)
	}
	return &refresh, nil
}

func (r authRepository) MarkRefreshTokenAsUsed(id uint64) error {
	refresh := auth.RefreshToken{
		Claimed: true,
	}
	return helpers.UpdateModelByID(&refresh, id)
}

func (r authRepository) MarkTokenAsRevoked(tokenID uint64) error {
	updatedFields := map[string]interface{}{
		"revoked": true,
	}

	return helpers.UpdateModelByIDWithMap[auth.AccessToken](updatedFields, tokenID)
}

func (r authRepository) FindTokenByUserIDAndToken(userID uint64, tokenString string) (*auth.AccessToken, error) {
	var token auth.AccessToken
	tokenString = strings.TrimSpace(tokenString)
	err := helpers.FindOneByField(&token, "user_id", userID, "token", tokenString, "revoked", false)
	if err != nil {
		return nil, fmt.Errorf("token not found or already revoked: %w", err)
	}
	return &token, nil
}
