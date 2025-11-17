package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/TryHanger/digital_signage/backend/internal/model"
	"github.com/TryHanger/digital_signage/backend/internal/repository"
	"github.com/TryHanger/digital_signage/backend/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nyaruka/phonenumbers"
	"gorm.io/gorm"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

type RegisterRequest struct {
	Phone    string
	Email    string
	Password string
}

type LoginRequest struct {
	Phone    string
	Email    string
	Password string
}

func (s *UserService) Register(ctx context.Context, req *RegisterRequest) (uint, error) {
	if req.Phone == "" && req.Email == "" {
		return 0, ErrPhoneOrEmailRequired
	}

	password, err := utils.HashPassword(req.Password)
	if err != nil {
		return 0, err
	}

	if req.Phone != "" {
		if err := ValidateKZPhone(req.Phone); err != nil {
			return 0, err
		}
	}

	var email string
	if req.Email != "" {
		email = strings.ToLower(strings.ReplaceAll(req.Email, " ", ""))
	}

	token := uuid.New().String()

	user := &model.User{
		Phone:        req.Phone,
		Email:        email,
		PassHash:     password,
		IsConfirmed:  false,
		ConfirmToken: token,
	}

	userID, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		fmt.Printf("DB error: %v\n", err)
		if dupErr := parseDuplicateError(err); dupErr != nil {
			return 0, dupErr
		}
		return 0, fmt.Errorf("create user: %w", err)
	}

	return userID, nil
}

func parseDuplicateError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
		switch pgErr.ConstraintName {
		case "idx_users_email_unique":
			return ErrEmailExists
		case "idx_users_phone_unique":
			return ErrPhoneExists
		default:
			return analyzeErrorContent(pgErr.Error())
		}
	}

	// Дополнительно ловим текстовые ошибки (на всякий случай)
	errStr := strings.ToLower(err.Error())
	if strings.Contains(errStr, "duplicate") || strings.Contains(errStr, "already exists") {
		return analyzeErrorContent(errStr)
	}

	return nil
}

func analyzeErrorContent(errStr string) error {
	// Приоритет: сначала ищем конкретные поля
	switch {
	case strings.Contains(errStr, "email") || strings.Contains(errStr, "e-mail"):
		return ErrEmailExists
	case strings.Contains(errStr, "phone") || strings.Contains(errStr, "telephone"):
		return ErrPhoneExists
	case strings.Contains(errStr, "duplicate") || strings.Contains(errStr, "already exists"):
		return ErrUserExists
	default:
		return ErrUserExists
	}
}

func ValidateKZPhone(phone string) error {
	num, err := phonenumbers.Parse(phone, "KZ")
	if err != nil {
		return fmt.Errorf("invalid phone format: %w", err)
	}

	if !phonenumbers.IsValidNumber(num) {
		return errors.New("not a valid KZ phone number")
	}

	// Дополнительно: проверяем что номер именно KZ
	countryCode := phonenumbers.GetRegionCodeForNumber(num)
	if countryCode != "KZ" {
		return errors.New("not a KZ phone number")
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, req *LoginRequest) (*model.User, string, error) {
	var user *model.User
	var err error

	switch {
	case req.Phone != "":
		user, err = s.repo.FindByPhone(ctx, req.Phone)

	case req.Email != "":
		email := strings.ToLower(strings.ReplaceAll(req.Email, " ", ""))
		user, err = s.repo.FindByEmail(ctx, email)

	default:
		return nil, "", ErrInvalidLoginFormat
	}

	if err != nil {
		// repository now returns gorm.ErrRecordNotFound on not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", ErrUserNotFound
		}
		return nil, "", err
	}

	if !utils.VerifyPassword(user.PassHash, req.Password) {
		return nil, "", ErrInvalidPassword
	}

	return user, "", nil
}
