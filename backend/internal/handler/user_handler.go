package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/TryHanger/digital_signage/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service      *service.UserService
	tokenService *service.TokenService
}

func NewUserHandler(service *service.UserService, tokenService *service.TokenService) *UserHandler {
	return &UserHandler{service: service, tokenService: tokenService}
}

func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/auth")
	{
		group.POST("/register", h.Register)
		group.POST("/login", h.Login)
		group.POST("/refresh", h.Refresh)
		group.POST("/logout", h.Logout)
	}
}

// Refresh обновляет access-token используя refresh-token из куки (с ротацией)
func (h *UserHandler) Refresh(c *gin.Context) {
	// read refresh token from cookie
	refresh, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "no_refresh_token", Message: "Не найден refresh token"})
		return
	}

	pair, err := h.tokenService.RefreshTokensWithRotation(refresh, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		fmt.Printf("refresh error: %v\n", err)
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid_refresh", Message: "Не удалось обновить токены"})
		return
	}

	// set new refresh cookie
	maxAge := 7 * 24 * 3600
	c.SetCookie("refresh_token", pair.RefreshToken, maxAge, "/", "", false, true)

	c.JSON(http.StatusOK, AuthResponse{AccessToken: pair.AccessToken, ExpiresAt: time.Now().Add(15 * time.Minute)})
}

// Logout аннулирует сессию по refresh-token и очищает куку
func (h *UserHandler) Logout(c *gin.Context) {
	refresh, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "no_refresh_token", Message: "Не найден refresh token"})
		return
	}

	if err := h.tokenService.RevokeSessionByPackedToken(refresh); err != nil {
		fmt.Printf("logout revoke error: %v\n", err)
		// even if revoke failed, try to clear cookie
	}

	// clear cookie
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, map[string]string{"message": "logged_out"})
}

// Register регистрация нового пользователя
// @Summary      Регистрация пользователя
// @Description  Создает нового пользователя по email или телефону
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Данные для регистрации"
// @Success      201 {object} RegisterResponse
// @Failure      400 {object} ErrorResponse
// @Failure      409 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/v1/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("JSON parse error: %v\n", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid request",
			Message: "Неверный формат данных",
		})
		return
	}

	userID, err := h.service.Register(c.Request.Context(), &service.RegisterRequest{
		Phone:    req.Phone,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		// Обработка ошибок сервиса
		fmt.Printf("Register service error: %v\n", err)
		switch {
		case errors.Is(err, service.ErrPhoneOrEmailRequired):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "phone_or_email_required",
				Message: "Укажите телефон или email",
			})

		case errors.Is(err, service.ErrPasswordRequired):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "password_required",
				Message: "Пароль обязателен",
			})

		case errors.Is(err, service.ErrEmailExists):
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "email_exists",
				Message: "Email уже зарегистрирован",
			})

		case errors.Is(err, service.ErrPhoneExists):
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "phone_exists",
				Message: "Телефон уже зарегистрирован",
			})

		case errors.Is(err, service.ErrUserExists):
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "user_exists",
				Message: "Пользователь уже существует",
			})

		default:
			// Логируем неизвестные ошибки для дебага
			fmt.Printf("Register unexpected error: %v\n", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "internal_error",
				Message: "Внутренняя ошибка сервера",
			})
		}
		return
	}

	// Создаём токены и ставим куку с refresh-token
	pair, err := h.tokenService.GenerateTokenPair(userID, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		fmt.Printf("failed to generate tokens: %v\n", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal_error", Message: "Не удалось создать сессию"})
		return
	}

	// Set refresh token as HttpOnly cookie (MaxAge in seconds)
	maxAge := 7 * 24 * 3600 // 7 days
	c.SetCookie("refresh_token", pair.RefreshToken, maxAge, "/", "", false, true)

	c.JSON(http.StatusCreated, AuthResponse{
		AccessToken: pair.AccessToken,
		ExpiresAt:   time.Now().Add(15 * time.Minute),
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var input LoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid request",
			Message: "Неверный формат данных",
		})
		return
	}

	user, _, err := h.service.Login(c.Request.Context(), &service.LoginRequest{
		Phone:    input.Phone,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		fmt.Printf("Login service error: %v\n", err)
		switch {
		case errors.Is(err, service.ErrUserNotFound), errors.Is(err, service.ErrInvalidPassword):
			// Объединяем для безопасности
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "invalid_credentials",
				Message: "Неверный логин или пароль",
			})

		case errors.Is(err, service.ErrInvalidLoginFormat):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "invalid_login_format",
				Message: "Неверный формат логина",
			})

		case errors.Is(err, service.ErrUserNotConfirmed):
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error:   "user_not_confirmed",
				Message: "Подтвердите email/телефон для входа",
			})

		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "internal_error",
				Message: "Внутренняя ошибка сервера",
			})
		}
		return
	}
	// Create token pair and set refresh cookie
	pair, err := h.tokenService.GenerateTokenPair(user.ID, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		fmt.Printf("failed to generate tokens: %v\n", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "internal_error", Message: "Не удалось создать сессию"})
		return
	}

	maxAge := 7 * 24 * 3600 // 7 days
	c.SetCookie("refresh_token", pair.RefreshToken, maxAge, "/", "", false, true)

	c.JSON(http.StatusOK, AuthResponse{AccessToken: pair.AccessToken, ExpiresAt: time.Now().Add(15 * time.Minute)})
}
