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

type CustomerHandler struct {
	CustomerService services.CustomerService
}

func NewCustomerHandler(service services.CustomerService) *CustomerHandler {
	return &CustomerHandler{CustomerService: service}
}

func (h *CustomerHandler) CheckPhoneCustomer(c echo.Context) error {
	var req dto.CheckPhoneRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
	}

	result, err := h.CustomerService.CheckPhone(c.Request().Context(), req)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, "something went wrong", err.Error())
	}

	return response.Success(c, http.StatusOK, "Phone check result", map[string]interface{}{
		"action": result.Action,
		"user":   result.User,
	})
}

func (h *CustomerHandler) RegisterCustomer(c echo.Context) error {
	var req dto.RegisterCustomerRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
	}

	data, err := h.CustomerService.RegisterCustomer(c.Request().Context(), req)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusBadRequest, "something went wrong", err.Error())
	}

	if data == nil {
		return response.Error(c, http.StatusBadRequest, "Internal Server Error", nil)
	}

	return response.Success(c, http.StatusCreated, "customer registered", data)
}

func (h *CustomerHandler) ResendCodeOTPVerify(c echo.Context) error {
	var req dto.ResendOTPRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
	}

	err := h.CustomerService.ResendCodeOTPVerify(c.Request().Context(), req)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusBadRequest, "something went wrong", err.Error())
	}

	return response.Success(c, http.StatusCreated, "resend code success", map[string]string{
		"phone":   req.Phone,
		"message": "OTP code sent successfully",
	})
}

func (h *CustomerHandler) VerifyCodeOTP(c echo.Context) error {
	var req dto.VerifyOTPRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
	}

	err := h.CustomerService.VerifyCodeOTP(c.Request().Context(), req)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusBadRequest, "something went wrong", err.Error())
	}

	return response.Success(c, http.StatusOK, "OTP verification success", nil)
}

func (h *CustomerHandler) NewPin(c echo.Context) error {
	var req dto.NewPinRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewCustomError(response.ErrBadRequest, "invalid request", 400))
	}

	if err := h.CustomerService.NewPinCustomer(c.Request().Context(), req); err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusBadRequest, "something went wrong", err.Error())
	}

	return response.Success(c, http.StatusOK, "PIN saved temporarily", nil)
}

func (h *CustomerHandler) ConfirmationPin(c echo.Context) error {
	var req dto.ConfirmPinRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.NewCustomError(response.ErrBadRequest, "invalid request", 400))
	}

	user, token, err := h.CustomerService.ConfirmPinCustomer(c.Request().Context(), req)

	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusBadRequest, "something went wrong", err.Error())
	}

	return response.Success(c, http.StatusOK, "PIN confirmed successfully", map[string]interface{}{
		"token": token,
		"user":  user,
	})

}

func (h *CustomerHandler) GetProfileCustomer(c echo.Context) error {
	userToken := c.Get("user")
	if userToken == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	user, ok := userToken.(*jwt.Token)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	claims := user.Claims.(jwt.MapClaims)
	custID := claims["user_id"].(string)

	authHeader := c.Request().Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	data, err := h.CustomerService.GetProfileCustomer(context.Background(), custID, token)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusBadRequest, "something went wrong", err.Error())
	}

	return response.Success(c, http.StatusOK, "Get Profile Success", data)
}

func (h *CustomerHandler) LogoutCustomer(c echo.Context) error {
	userToken := c.Get("user")
	if userToken == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	user, ok := userToken.(*jwt.Token)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	claims := user.Claims.(jwt.MapClaims)

	role := claims["role"]
	if role != "customer" {
		return response.Error(c, http.StatusUnauthorized, "Invalid role access", nil)
	}

	custID := claims["user_id"].(string)

	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return response.Error(c, http.StatusUnauthorized, "Missing Authorization header", nil)
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	ctx := c.Request().Context()

	exists, err := h.CustomerService.CheckTokenBlacklisted(ctx, token)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error(), "internal server error")
	}
	if exists {
		return response.Error(c, http.StatusUnauthorized, "You are already logged out", nil)
	}

	if token != "" {
		if err := h.CustomerService.LogoutCustomer(ctx, custID, token); err != nil {
			if customErr, ok := response.AsCustomErr(err); ok {
				return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
			}
			return response.Error(c, http.StatusBadRequest, "something went wrong", err.Error())
		}
	} else {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: Token is empty", nil)
	}

	return response.Success(c, http.StatusOK, "Logout Success", nil)
}

func (h *CustomerHandler) LoginCustomer(c echo.Context) error {
	var req dto.LoginCustomerRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
	}

	token, err := h.CustomerService.LoginCustomer(c.Request().Context(), req)
	if err != nil {
		var cooldownErr *response.CooldownError
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
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
		return response.Error(c, http.StatusBadRequest, "something went wrong", err.Error())
	}

	return response.Success(c, http.StatusOK, "login success", token)
}
