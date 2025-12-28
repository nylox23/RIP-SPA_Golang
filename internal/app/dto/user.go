package dto

import (
	"web_service_auth/internal/app/role"

	"github.com/google/uuid"
)

type UserRegisterRequest struct {
	Login    string `json:"login" binding:"required,min=3,max=25"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserRegisterResponse struct {
	Ok bool `json:"ok"`
}

type UserLoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserUpdateRequest struct {
	Login    string `json:"login" binding:"min=3,max=25"`
	Password string `json:"password" binding:"min=6,max=100"`
}

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Login string    `json:"login"`
	Role  role.Role `json:"role"`
}

type LoginResponse struct {
	ExpiresIn   float32 `json:"expires_in"`
	AccessToken string  `json:"access_token"`
	TokenType   string  `json:"token_type"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
