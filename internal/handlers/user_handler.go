package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"nusantara_service/internal/data/services"
	"nusantara_service/internal/dto"
	"nusantara_service/internal/response"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserService services.UserService
}

func NewAuthHandler(service services.UserService) *UserHandler {
	return &UserHandler{UserService: service}
}

func (h *UserHandler) RegisterAdmin(c echo.Context) error {
	var req dto.RegisterAdminRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
	}

	if strings.TrimSpace(req.Username) == "" {
		return response.Error(c, http.StatusBadRequest, "name is required", nil)
	}
	if strings.TrimSpace(req.Email) == "" {
		return response.Error(c, http.StatusBadRequest, "email is required", nil)
	}
	if strings.TrimSpace(req.Password) == "" {
		return response.Error(c, http.StatusBadRequest, "password is required", nil)
	}
	if strings.TrimSpace(req.RoleID) == "" {
		return response.Error(c, http.StatusBadRequest, "role is required", nil)
	}

	data, err := h.UserService.RegisterAdmin(c.Request().Context(), req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error(), "internal server error")
	}

	return response.Success(c, http.StatusCreated, "admin registered", *data)
}

func (h *UserHandler) LoginAdmin(c echo.Context) error {
	var req dto.LoginAdminRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
	}

	if strings.TrimSpace(req.Email) == "" {
		return response.Error(c, http.StatusBadRequest, "email is required", nil)
	}

	if strings.TrimSpace(req.Password) == "" {
		return response.Error(c, http.StatusBadRequest, "password is required", nil)
	}

	token, err := h.UserService.LoginAdmin(c.Request().Context(), req)
	if err != nil {
		var cooldownErr *response.CooldownError
		if errors.As(err, &cooldownErr) {
			return response.Error(c, http.StatusTooManyRequests, cooldownErr.Message, map[string]interface{}{
				"retry_after_seconds": int(cooldownErr.RetryAfter.Seconds()),
				"retry_after_human":   fmt.Sprintf("%.0f seconds", cooldownErr.RetryAfter.Seconds()),
			})
		}
		var attemptErr *response.LoginAttemptError
		if errors.As(err, &attemptErr) {
			return response.Error(c, http.StatusUnauthorized, attemptErr.Message, map[string]interface{}{
				"remaining_attempts": fmt.Sprintf("%s. Sisa %d percobaan lagi.", attemptErr.Message, attemptErr.RemainingAttempts),
			})
		}
		return response.Error(c, http.StatusUnauthorized, err.Error(), "invalid login")
	}

	return response.Success(c, http.StatusOK, "login success", token)
}

func (h *UserHandler) GetProfileAdmin(c echo.Context) error {
	userToken := c.Get("user")
	if userToken == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	user, ok := userToken.(*jwt.Token)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	authHeader := c.Request().Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	profile, err := h.UserService.GetProfile(context.Background(), userID, token)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized", err.Error())
	}

	return response.Success(c, http.StatusOK, "Get Profile Success", profile)
}

func (h *UserHandler) LogoutAdmin(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return response.Error(c, http.StatusUnauthorized, "Missing Authorization header", nil)
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	ctx := c.Request().Context() // Use request context here

	// Check if token is already blacklisted
	exists, err := h.UserService.CheckTokenBlacklisted(ctx, token)
	if err != nil {
		// Only return 500 if it's a true Redis error (not redis.Nil)
		return response.Error(c, http.StatusInternalServerError, "internal server error", err.Error())
	}
	if exists {
		return response.Error(c, http.StatusUnauthorized, "You are already logged out", nil)
	}

	// Blacklist the token
	if token != "" {
		if err := h.UserService.LogoutAdmin(ctx, token); err != nil { // Use request context
			return response.Error(c, http.StatusInternalServerError, err.Error(), "internal server error")
		}
	} else {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: Token is empty", nil) // More specific error
	}

	return response.Success(c, http.StatusOK, "Logout Success", nil)
}

func (h *UserHandler) ChangePasswordSuperAdmin(c echo.Context) error {
	userToken := c.Get("user")
	if userToken == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	user, ok := userToken.(*jwt.Token)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	authHeader := c.Request().Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	var req dto.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Failed to bind request", err.Error())
	}

	if strings.TrimSpace(req.CurrentPassword) == "" {
		return response.Error(c, http.StatusBadRequest, "current password is required", nil)
	}

	if strings.TrimSpace(req.NewPassword) == "" {
		return response.Error(c, http.StatusBadRequest, "new password is required", nil)
	}
	if strings.TrimSpace(req.ConfirmationPassword) == "" {
		return response.Error(c, http.StatusBadRequest, "confirmation password is required", nil)
	}

	userUpdate, err := h.UserService.ChangePasswordSuperAdmin(c.Request().Context(), userID, token, req)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, "something went wrong", err.Error())
	}

	return response.Success(c, http.StatusOK, "change password success", userUpdate)
}

//CUSTOMER
