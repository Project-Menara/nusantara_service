package handlers

import (
	"net/http"
	"nusantara_service/internal/data/services"
	dto "nusantara_service/internal/dto/request"
	"nusantara_service/internal/response"
	"nusantara_service/internal/utils"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type BannerHandler struct {
	BannerService services.BannerService
}

func NewBannerHandler(service services.BannerService) *BannerHandler {
	return &BannerHandler{BannerService: service}
}

func (b *BannerHandler) CreateBanner(c echo.Context) error {
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

	// Manual parsing form-data fields
	var req dto.CreateBannerRequest
	req.Name = c.FormValue("name")
	req.Description = c.FormValue("description")
	status := c.FormValue("status")
	statusInt, err := strconv.Atoi(status)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid status format", err.Error())
	}
	req.Status = statusInt

	// Get image
	imageFile, err := c.FormFile("image")
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "failed to get photo file", err.Error())
	}

	// Call service
	newBanner, err := b.BannerService.CreateBanner(c.Request().Context(), superAdminID, req, imageFile)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to create banner")
	}

	return response.Success(c, http.StatusCreated, "Create Banner Success", newBanner)
}

func (b *BannerHandler) GetAllBanner(c echo.Context) error {
	pageInt, limtiInt := utils.ParsePaginationParams(c, 10)
	search := c.QueryParam("search")

	banners, total, err := b.BannerService.GetAllBanner(c.Request().Context(), pageInt, limtiInt, search)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get banner")
	}

	meta := utils.BuildPaginationMeta(c, pageInt, limtiInt, total)

	return response.PaginatedSuccess(c, 200, "Get All Banner Success", banners, meta)
}

func (b *BannerHandler) GetByIdBanner(c echo.Context) error {
	bannerId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	banner, err := b.BannerService.GetByIdBanner(c.Request().Context(), bannerId)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to create banner")
	}

	return response.Success(c, http.StatusOK, "Get Banner Success", banner)
}

func (b *BannerHandler) UpdateBanner(c echo.Context) error {
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

	bannerId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	var req dto.UpdateBannerRequest
	req.Name = c.FormValue("name")
	req.Description = c.FormValue("description")
	imageFile, _ := c.FormFile("image")

	updatedBanner, err := b.BannerService.UpdateBanner(c.Request().Context(), superAdminID, bannerId, req, imageFile)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to create banner")
	}

	return response.Success(c, http.StatusOK, "Updated Banner Success", updatedBanner)
}

func (b *BannerHandler) DeleteBanner(c echo.Context) error {
	bannerId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	if err := b.BannerService.DeleteBanner(c.Request().Context(), bannerId); err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to create banner")
	}

	return response.Success(c, http.StatusOK, "Deleted Success", nil)

}

func (b *BannerHandler) UpdateStatusBanner(c echo.Context) error {
	bannerId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	var req dto.UpdateStatusBannerRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Failed to bind request", err.Error())
	}

	if err := b.BannerService.UpdateStatusBanner(c.Request().Context(), bannerId.String(), req); err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to update status banner")
	}

	return response.Success(c, http.StatusOK, "Update Status Banner Success", nil)
}

func (b *BannerHandler) GetAllBannerCustomer(c echo.Context) error {
	banners, err := b.BannerService.GetAllBannerCustomer(c.Request().Context())
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to Get banner")
	}

	return response.Success(c, http.StatusOK, "Get All Banner Success", banners)
}

func (b *BannerHandler) GetByIdBannerCustomer(c echo.Context) error {
	bannerId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	banner, err := b.BannerService.GetByIdBannerCustomer(c.Request().Context(), bannerId)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get type product")
	}

	return response.Success(c, http.StatusOK, "Get Banner Success", banner)
}
