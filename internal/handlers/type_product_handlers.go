package handlers

import (
	"net/http"
	"nusantara_service/internal/data/services"
	"nusantara_service/internal/dto"
	"nusantara_service/internal/response"
	"nusantara_service/internal/utils"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TypeProductHandler struct {
	TypeProductService services.TypeProductService
}

func NewTypeProductHandler(service services.TypeProductService) *TypeProductHandler {
	return &TypeProductHandler{TypeProductService: service}
}

func (t *TypeProductHandler) CreateTypeProduct(c echo.Context) error {
	userToken := c.Get("user")
	if userToken == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	user, ok := userToken.(*jwt.Token)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	claims := user.Claims.(jwt.MapClaims)
	superAdminID := claims["user_id"].(string)

	var req dto.CreateTypeProductRequest
	req.Name = c.FormValue("name")
	status := c.FormValue("status")
	statusInt, err := strconv.Atoi(status)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid status format", err.Error())
	}
	req.Status = statusInt

	imageFile, err := c.FormFile("image")
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "failed to get photo file", err.Error())
	}

	newTypeProduct, err := t.TypeProductService.CreateTypeProduct(c.Request().Context(), superAdminID, req, imageFile)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to create type product")
	}

	return response.Success(c, http.StatusCreated, "Create Type Product Success", newTypeProduct)
}

func (t *TypeProductHandler) GetAllTypeProduct(c echo.Context) error {
	pageInt, limtiInt := utils.ParsePaginationParams(c, 10)
	search := c.QueryParam("search")

	typeProduct, total, err := t.TypeProductService.GetAllTypeProduct(c.Request().Context(), pageInt, limtiInt, search)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get type product")
	}

	meta := utils.BuildPaginationMeta(c, pageInt, limtiInt, total)

	return response.PaginatedSuccess(c, 200, "Get All Type Product Success", typeProduct, meta)

}

func (t *TypeProductHandler) GetByIdTypeProduct(c echo.Context) error {
	typeProductId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	typeProduct, err := t.TypeProductService.GetByIdTypeProduct(c.Request().Context(), typeProductId)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get type product")
	}

	return response.Success(c, http.StatusOK, "Get Type Product Success", typeProduct)
}

func (t *TypeProductHandler) UpdateTypeProduct(c echo.Context) error {
	userToken := c.Get("user")
	if userToken == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	user, ok := userToken.(*jwt.Token)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	claims := user.Claims.(jwt.MapClaims)
	superAdminID := claims["user_id"].(string)

	typeProductId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	var req dto.UpdateTypeProductRequest
	req.Name = c.FormValue("name")
	imageFile, _ := c.FormFile("image")

	updated, err := t.TypeProductService.UpdateTypeProduct(c.Request().Context(), superAdminID, typeProductId, req, imageFile)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to update type product")
	}

	return response.Success(c, http.StatusOK, "Upadated Type Product Success", updated)
}

func (t *TypeProductHandler) DeleteTypeProduct(c echo.Context) error {
	typeProductId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	if err := t.TypeProductService.DeleteTypeProduct(c.Request().Context(), typeProductId); err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to delete type product")
	}

	return response.Success(c, http.StatusOK, "Deleted Success", nil)

}

func (t *TypeProductHandler) UpdateStatusTypeProduct(c echo.Context) error {
	typeProductId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	var req dto.UpdateStatusTypeProductRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Failed to bind request", err.Error())
	}

	if err := t.TypeProductService.UpdateStatusTypeProduct(c.Request().Context(), typeProductId, req); err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to update status type product")
	}

	return response.Success(c, http.StatusOK, "Update Status Type Product Success", nil)

}

func (t *TypeProductHandler) GetAllTypeProductCustomer(c echo.Context) error {
	typeProduct, err := t.TypeProductService.GetAllTypeProductCustomer(c.Request().Context())
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get type product")
	}
	return response.Success(c, http.StatusOK, "Get All Type Product Success", typeProduct)
}

func (t *TypeProductHandler) GetByIdTypeProductCustomer(c echo.Context) error {
	typeProductId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	typeProduct, err := t.TypeProductService.GetByIdTypeProductCustomer(c.Request().Context(), typeProductId)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get type product")
	}

	return response.Success(c, http.StatusOK, "Get Type Product Success", typeProduct)
}
