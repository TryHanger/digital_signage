package repository

import (
	"context"
	"time"

	"github.com/TryHanger/digital_signage/backend/internal/model"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateRefreshSession(ctx context.Context, session *model.RefreshSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *AuthRepository) GetSessionByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshSession, error) {
	var session model.RefreshSession
	err := r.db.WithContext(ctx).
		Where("token_hash = ?", tokenHash).
		First(&session).Error
	return &session, err
}

func (r *AuthRepository) GetSessionByJTI(ctx context.Context, jti string) (*model.RefreshSession, error) {
	var session model.RefreshSession
	err := r.db.WithContext(ctx).
		Where("jti = ?", jti).
		First(&session).Error
	return &session, err
}

func (r *AuthRepository) RevokeSession(ctx context.Context, sessionID uint) error {
	return r.db.WithContext(ctx).
		Model(&model.RefreshSession{}).
		Where("id = ?", sessionID).
		Update("revoked", true).Error
}

func (r *AuthRepository) RevokeSessionByJTI(ctx context.Context, jti string) error {
	return r.db.WithContext(ctx).
		Model(&model.RefreshSession{}).
		Where("jti = ?", jti).
		Update("revoked", true).Error
}

func (r *AuthRepository) RevokeAllSessionsByUser(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Model(&model.RefreshSession{}).
		Where("user_id = ? AND revoked = ?", userID, false).
		Update("revoked", true).Error
}

// CleanupExpiredSessions deletes sessions which have expired
func (r *AuthRepository) CleanupExpiredSessions(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&model.RefreshSession{}).Error
}
