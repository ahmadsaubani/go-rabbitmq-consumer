package auth_dtos

import (
	"subscriber-topic-stars/src/entities/users"
	"time"
)

type AccessTokenDto struct {
	ID        uint64      `json:"id"`
	UserID    uint64      `json:"user_id"`
	User      *users.User `json:"user"`
	Token     string      `json:"token"`
	ExpiresAt time.Time   `json:"expires_at"`
	Revoked   bool        `json:"revoked"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type TokenResultDto struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	AccessExpiresAt  time.Time `json:"access_expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}

type RequestLoginDto struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

type RequestRefreshTokenDto struct {
	RefreshToken string `form:"refresh_token" json:"refresh_token" binding:"required"`
}

type RequestRegisterDto struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required,min=8"`
}
