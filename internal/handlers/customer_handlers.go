package handlers

import (
	"net/http"
	"nusantara_service/internal/data/services"
	"nusantara_service/internal/dto"
	"nusantara_service/internal/response"

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

	user, err := h.CustomerService.CheckPhone(c.Request().Context(), req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error(), nil)
	}

	if user != nil {
		return response.Success(c, http.StatusOK, "Phone Number Registered", map[string]interface{}{
			"action": "login",
			"user":   user,
		})
	}

	return response.Success(c, http.StatusOK, "Phone number not registered", map[string]interface{}{
		"action": "register",
	})
}

func (h *CustomerHandler) RegisterCustomer(c echo.Context) error {
	var req dto.RegisterCustomerRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request", err.Error())
	}

	data, err := h.CustomerService.RegisterCustomer(c.Request().Context(), req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error(), nil)
	}

	if data == nil {
		return response.Error(c, http.StatusBadRequest, "Internal Server Error", nil)
	}

	return response.Success(c, http.StatusCreated, "customer registered", data)
}
