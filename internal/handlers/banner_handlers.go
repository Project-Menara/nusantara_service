package handlers

import (
	"net/http"
	"nusantara_service/internal/data/services"
	"nusantara_service/internal/dto"
	"nusantara_service/internal/response"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
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
