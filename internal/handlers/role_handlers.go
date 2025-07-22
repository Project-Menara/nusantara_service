package handlers

import (
	"fmt"
	"net/http"
	"nusantara_service/internal/data/services"
	"nusantara_service/internal/dto"
	"nusantara_service/internal/response"
	"nusantara_service/internal/utils"
	"strings"

	"github.com/google/uuid"
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

func (r *RoleHandler) GetAllRole(c echo.Context) error {
	pageInt, limitInt := utils.ParsePaginationParams(c)

	roles, total, err := r.RoleService.GetAllRole(c.Request().Context(), pageInt, limitInt)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "failed to fetch roles", err)
	}

	meta := utils.BuildPaginationMeta(c, pageInt, limitInt, total)

	return response.PaginatedSuccess(c, 200, "Get All Role Success", roles, meta)

}

func (r *RoleHandler) GetByIdRole(c echo.Context) error {
	roleId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, 400, "invalid UUID", err.Error())
	}

	role, err := r.RoleService.GetByIdRole(c.Request().Context(), roleId)
	if err != nil {
		return response.Error(c, 404, "Role Not Found", err.Error())
	}

	return response.Success(c, 200, "Get Role Success", role)
}

func (r *RoleHandler) UpdateRole(c echo.Context) error {
	roleId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, 400, "invalid UUID", err.Error())
	}

	var req dto.UpdateRoleRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, 400, "Failed to bind request", err.Error())
	}

	if strings.TrimSpace(req.Name) == "" {
		return response.Error(c, 400, "Name is required", err)
	}

	updateRole, err := r.RoleService.UpdateRole(c.Request().Context(), roleId, req)
	if err != nil {
		return response.Error(c, 500, "Failed to update role", err.Error())
	}

	return response.Success(c, 200, "update success", updateRole)
}

func (r *RoleHandler) DeleteRole(c echo.Context) error {
	roleId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, 400, "invalid UUID", err.Error())
	}

	if err := r.RoleService.DeleteRole(c.Request().Context(), roleId); err != nil {
		return response.Error(c, 400, "failed to delete role", err.Error())
	}

	return response.Success(c, 200, "deleted success", nil)
}
