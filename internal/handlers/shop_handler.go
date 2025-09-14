package handlers

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"nusantara_service/internal/data/services"
	dto "nusantara_service/internal/dto/request"
	shopresponse "nusantara_service/internal/dto/responses/shop_response"
	"nusantara_service/internal/response"
	"nusantara_service/internal/utils"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ShopHandler struct {
	ShopService services.ShopService
}

func NewShopHandler(service services.ShopService) *ShopHandler {
	return &ShopHandler{ShopService: service}
}

func (s *ShopHandler) CreateShop(c echo.Context) error {
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

	var req dto.CreateShopRequest
	req.Name = c.FormValue("name")
	req.Description = c.FormValue("description")
	req.FullAddress = c.FormValue("full_address")

	if v := c.FormValue("lat"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			req.Lat = &f
		}
	}
	if v := c.FormValue("lang"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			req.Lng = &f
		}
	}
	if v := c.FormValue("status"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			req.Status = &n
		}
	}

	// Cover
	if cover, err := c.FormFile("cover"); err == nil {
		req.Cover = cover
	}

	// Gallery
	gallery := []*multipart.FileHeader{}
	formGallery, _ := c.MultipartForm()
	if formGallery != nil {
		if files, ok := formGallery.File["gallery"]; ok {
			gallery = files
		}
	}
	req.Gallery = gallery

	// Cashier IDs (array of uuid strings)
	if cashierIDs := c.FormValue("cashier_ids"); cashierIDs != "" {
		var ids []string
		if err := json.Unmarshal([]byte(cashierIDs), &ids); err == nil {
			for _, id := range ids {
				if parsed, err := uuid.Parse(id); err == nil {
					req.CashierIDs = append(req.CashierIDs, parsed)
				}
			}
		}
	}

	// Products (array JSON)
	if products := c.FormValue("products"); products != "" {
		var parsed []dto.AssignProductRequest
		if err := json.Unmarshal([]byte(products), &parsed); err == nil {
			req.Products = parsed
		}
	}

	err := s.ShopService.CreateShop(c.Request().Context(), uuid.MustParse(superAdminID), req)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to create shop")
	}

	return response.Success(c, http.StatusOK, "Shop Created Successfully", nil)
}

func (s *ShopHandler) GetAllShop(c echo.Context) error {
	pageInt, limitInt := utils.ParsePaginationParams(c, 10)
	search := c.QueryParam("search")

	shops, total, err := s.ShopService.GetAllShop(c.Request().Context(), pageInt, limitInt, search)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get shops")
	}

	meta := utils.BuildPaginationMeta(c, pageInt, limitInt, total)
	data := make([]shopresponse.ShopResponse, len(shops))
	for i, s := range shops {
		data[i] = shopresponse.ToShopResponse(*s)
	}

	return response.PaginatedSuccess(c, http.StatusOK, "Shops Retrieved Successfully", data, meta)
}

func (s *ShopHandler) GetByIdShop(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	shop, err := s.ShopService.GetByIdShop(c.Request().Context(), id)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get shop")
	}

	return response.Success(c, http.StatusOK, "Shop Retrieved Successfully", shopresponse.ToShopResponse(*shop))
}

func (s *ShopHandler) UpdateShop(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	var req dto.UpdateShopRequest
	req.Name = c.FormValue("name")
	req.Description = c.FormValue("description")
	req.FullAddress = c.FormValue("full_address")

	if v := c.FormValue("lat"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			req.Lat = f
		}
	}
	if v := c.FormValue("lang"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			req.Lng = f
		}
	}
	if v := c.FormValue("replace_gallery"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			req.ReplaceGallery = b
		}
	}

	// New cover
	if newCover, err := c.FormFile("cover"); err == nil {
		req.NewCover = newCover
	}

	// New gallery
	newGallery := []*multipart.FileHeader{}
	formGallery, _ := c.MultipartForm()
	if formGallery != nil {
		if files, ok := formGallery.File["gallery"]; ok {
			newGallery = files
		}
	}
	req.NewGallery = newGallery

	// Cashier IDs
	if cashierIDs := c.FormValue("cashier_ids"); cashierIDs != "" {
		var ids []string
		if err := json.Unmarshal([]byte(cashierIDs), &ids); err == nil {
			for _, id := range ids {
				if parsed, err := uuid.Parse(id); err == nil {
					req.CashierIDs = append(req.CashierIDs, parsed)
				}
			}
		}
	}

	// Products
	if products := c.FormValue("products"); products != "" {
		var parsed []dto.AssignProductRequest
		if err := json.Unmarshal([]byte(products), &parsed); err == nil {
			req.Products = parsed
		}
	}

	err = s.ShopService.UpdateShop(c.Request().Context(), uuid.Nil, id, req)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to update shop")
	}

	return response.Success(c, http.StatusOK, "Shop Updated Successfully", nil)
}

func (s *ShopHandler) DeleteShop(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	if err := s.ShopService.Delete(c.Request().Context(), id); err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to delete shop")
	}

	return response.Success(c, http.StatusOK, "Shop Deleted Successfully", nil)
}

func (s *ShopHandler) UpdateStatusShop(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	var req dto.UpdateStatusShopRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "failed to bind request", err.Error())
	}

	if err := s.ShopService.UpdateStatusShop(c.Request().Context(), id, req); err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to delete shop")
	}

	return response.Success(c, http.StatusOK, "update status success", nil)
}
