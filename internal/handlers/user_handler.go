package handlers

import (
	"net/http"
	"nusantara_service/internal/data/services"
	"nusantara_service/internal/dto"
	"nusantara_service/internal/response"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	AuthService services.AuthService
}

func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: service}
}

func (h *AuthHandler) RegisterUser(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
	}

	if err := h.AuthService.Register(c.Request().Context(), req); err != nil {
		return response.Error(c, http.StatusInternalServerError, "internal server error", err.Error())
	}

	return response.Success(c, http.StatusCreated, "user registered", req)
}

func (h *AuthHandler) LoginUser(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
	}

	token, err := h.AuthService.Login(c.Request().Context(), req)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, "invalid login", err.Error())
	}

	return response.Success(c, http.StatusOK, "login successful", map[string]string{
		"token": token,
	})
}

func (h *AuthHandler) GetProfile(c echo.Context) error {
	userToken, ok := c.Get("user").(*jwt.Token)
	if !ok || userToken == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
	}

	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	return response.Success(c, http.StatusOK, "Get user successfully", map[string]string{
		"user_id": userID,
	})
}

func (h *AuthHandler) LogoutUser(c echo.Context) error {
	userToken, ok := c.Get("user").(*jwt.Token)
	if !ok || userToken == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
	}

	token := userToken.Raw

	if err := h.AuthService.Logout(c.Request().Context(), token); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Internal server error", err.Error())
	}

	return response.Success(c, http.StatusOK, "Logout successfully", nil)
}
