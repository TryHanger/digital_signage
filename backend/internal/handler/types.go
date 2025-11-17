package handler

import "time"

// Request структуры
type RegisterRequest struct {
	Phone    string `json:"phone" binding:"omitempty,e164"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// Response структуры
type RegisterResponse struct {
	Message string `json:"message"`
	UserID  uint   `json:"user_id,omitempty"`
}

// AuthResponse возвращается при успешном входе/регистрации
type AuthResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type LoginRequest struct {
	Phone    string `json:"phone" blinding:"omitempty,e164"`
	Email    string `json:"email" blinding:"omitempty,email"`
	Password string `json:"password" blinding:"required,min=8"`
}
