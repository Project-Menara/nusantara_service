package handlers

import (
	"fmt"
	"net/http"
	"nusantara_service/internal/data/services"
	"nusantara_service/internal/dto"
	"nusantara_service/internal/response"
	"strings"

	"github.com/labstack/echo/v4"
)

type RoleHandler struct {
	RoleService services.RoleService
}

func NewRoleHandler(service services.RoleService) *RoleHandler {
	return &RoleHandler{RoleService: service}
}

func (r *RoleHandler) CreateRole(c echo.Context) error {
	var req dto.CreateRoleRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Failed to create role", err.Error())
	}

	if strings.TrimSpace(req.Name) == "" {
		return response.Error(c, 400, "validation error", "name is required")
	}

	newRole, err := r.RoleService.CreateRole(c.Request().Context(), req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "internal server error", err)
	}
	fmt.Printf("DEBUG newRole: %+v\n", newRole)

	return response.Success(c, 201, "created", newRole)
}
