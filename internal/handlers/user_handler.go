package handlers

import (
	"context"
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
		return response.Error(c, http.StatusInternalServerError, "internal server error", err.Error())
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
		return response.Error(c, http.StatusUnauthorized, "invalid login", err.Error())
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
			return response.Error(c, http.StatusInternalServerError, "internal server error", err.Error())
		}
	} else {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: Token is empty", nil) // More specific error
	}

	return response.Success(c, http.StatusOK, "Logout Success", nil)
}
