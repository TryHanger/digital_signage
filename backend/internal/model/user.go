package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Phone        string `json:"phone" gorm:"uniqueIndex:idx_users_phone_unique;default:null"`
	Email        string `json:"email" gorm:"uniqueIndex:idx_users_email_unique;default:null"`
	PassHash     string
	IsConfirmed  bool
	ConfirmToken string
}

type RefreshSession struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null;index"`
	JTI       string `gorm:"uniqueIndex;not null"`
	TokenHash string `gorm:"uniqueIndex;not null"`
	IPAddress string
	UserAgent string
	ExpiresAt time.Time `gorm:"not null;index:idx_refresh_sessions_expires_at"`
	Revoked   bool      `gorm:"default:false;index:idx_refresh_sessions_revoked"`
	CreatedAt time.Time `gorm:"index:idx_refresh_sessions_created_at"`
	UpdatedAt time.Time
}
