package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/TryHanger/digital_signage/backend/internal/model"
	"github.com/TryHanger/digital_signage/backend/internal/repository"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type TokenService struct {
	repo       *repository.AuthRepository
	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type RefreshTokenResult struct {
	RawToken  string
	TokenHash string
	JTI       string
	ExpiresAt time.Time
}

type TokenClaims struct {
	UserID uint   `json:"user_id"`
	JTI    string `json:"jti"`
	jwt.RegisteredClaims
}

func NewTokenService(repo *repository.AuthRepository, jwtSecret string) *TokenService {
	return &TokenService{
		repo:       repo,
		jwtSecret:  []byte(jwtSecret),
		accessTTL:  15 * time.Minute,
		refreshTTL: 7 * 24 * time.Hour,
	}
}

func (s *TokenService) GenerateAccess(userID uint, jti string) (string, error) {
	sub := strconv.FormatUint(uint64(userID), 10)
	claims := TokenClaims{
		UserID: userID,
		JTI:    jti,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *TokenService) GenerateRefresh() (*RefreshTokenResult, error) {
	rawToken := make([]byte, 32)
	if _, err := rand.Read(rawToken); err != nil {
		return nil, err
	}
	rawTokenHex := hex.EncodeToString(rawToken)

	jtiBytes := make([]byte, 16)
	if _, err := rand.Read(jtiBytes); err != nil {
		return nil, err
	}

	jti := hex.EncodeToString(jtiBytes)

	tokenHash, err := bcrypt.GenerateFromPassword([]byte(rawTokenHex), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	expireAt := time.Now().Add(s.refreshTTL)

	// Pack JTI into raw token so we can lookup session by JTI later
	packed := jti + ":" + rawTokenHex

	return &RefreshTokenResult{
		RawToken:  packed,
		TokenHash: string(tokenHash),
		JTI:       jti,
		ExpiresAt: expireAt,
	}, nil
}

// validateRefreshToken verifies provided packed refresh token (format: jti:rawHex)
func (s *TokenService) validateRefreshToken(packedToken string) (*model.RefreshSession, error) {
	parts := strings.SplitN(packedToken, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid refresh token format")
	}
	jti := parts[0]
	rawHex := parts[1]

	session, err := s.repo.GetSessionByJTI(context.Background(), jti)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Compare stored bcrypt hash with raw token
	if err := bcrypt.CompareHashAndPassword([]byte(session.TokenHash), []byte(rawHex)); err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Check expiry
	if time.Now().After(session.ExpiresAt) {
		return nil, fmt.Errorf("refresh token expired")
	}

	return session, nil
}

// ValidateAccessToken parses and validates JWT access token and returns claims
func (s *TokenService) ValidateAccessToken(tokenStr string) (*TokenClaims, error) {
	if tokenStr == "" {
		return nil, fmt.Errorf("empty token")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) && ve.Errors&jwt.ValidationErrorExpired != 0 {
			return nil, fmt.Errorf("token expired: %w", err)
		}
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token claims")
}

func (s *TokenService) CreateSession(userID uint, ip, userAgent string) (*RefreshTokenResult, error) {
	refreshResult, err := s.GenerateRefresh()
	if err != nil {
		return nil, err
	}

	session := &model.RefreshSession{
		UserID:    userID,
		JTI:       refreshResult.JTI,
		TokenHash: refreshResult.TokenHash,
		IPAddress: ip,
		UserAgent: userAgent,
		ExpiresAt: refreshResult.ExpiresAt,
		Revoked:   false,
	}

	err = s.repo.CreateRefreshSession(context.Background(), session)
	if err != nil {
		return nil, fmt.Errorf("Failed to create refresh session: %v", err)
	}
	return refreshResult, nil
}

func (s *TokenService) GenerateTokenPair(userID uint, ip, userAgent string) (*TokenPair, error) {
	refreshResult, err := s.CreateSession(userID, ip, userAgent)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.GenerateAccess(userID, refreshResult.JTI)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshResult.RawToken,
		ExpiresAt:    time.Now().Add(s.accessTTL),
	}, nil
}

func (s *TokenService) RefreshTokensWithRotation(oldRefreshToken string, ip, userAgent string) (*TokenPair, error) {
	ctx := context.Background()

	oldSession, err := s.validateRefreshToken(oldRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to validate old refresh token: %v", err)
	}

	if oldSession.Revoked {
		// revoke all active sessions for user
		_ = s.repo.RevokeAllSessionsByUser(ctx, oldSession.UserID)
		return nil, fmt.Errorf("security alert: token reuse detected - all sessions revoked")
	}

	err = s.repo.RevokeSession(ctx, oldSession.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to revoke old session: %w", err)
	}

	newTokenPair, err := s.GenerateTokenPair(oldSession.UserID, ip, userAgent)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	return newTokenPair, nil
}

func (s *TokenService) ValidateAccessTokenWithRotation(accessToken, refreshToken string, ip, userAgent string) (*TokenClaims, *TokenPair, error) {
	claims, err := s.ValidateAccessToken(accessToken)
	if err != nil {
		return claims, nil, nil
	}

	if refreshToken == "" {
		return nil, nil, fmt.Errorf("access token expired and no refresh token provided")
	}

	newTokenPair, err := s.RefreshTokensWithRotation(refreshToken, ip, userAgent)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to refresh tokens: %w", err)
	}

	newClaims, err := s.ValidateAccessToken(newTokenPair.AccessToken)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to validate new access token: %w", err)
	}

	return newClaims, newTokenPair, nil
}

func (s *TokenService) cleanupExpiredSessions() {
	ctx := context.Background()
	_ = s.repo.CleanupExpiredSessions(ctx)
}

// RevokeSessionByPackedToken validates packed refresh token and revokes the corresponding session
func (s *TokenService) RevokeSessionByPackedToken(packedToken string) error {
	ctx := context.Background()
	session, err := s.validateRefreshToken(packedToken)
	if err != nil {
		return err
	}
	return s.repo.RevokeSession(ctx, session.ID)
}
