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

type ProducHandlers struct {
	ProductService services.ProductService
}

func NewProductHandler(service services.ProductService) *ProducHandlers {
	return &ProducHandlers{ProductService: service}
}

func (p *ProducHandlers) CreateProduct(c echo.Context) error {
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

	var req dto.CreateProductRequest
	req.Name = c.FormValue("name")
	req.Code = c.FormValue("code")
	req.Unit = c.FormValue("unit")
	req.Description = c.FormValue("description")

	if v := c.FormValue("price"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			req.Price = n
		}
	}
	if v := c.FormValue("status"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			req.Status = n
		}
	}
	if v := c.FormValue("type_product_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			req.TypeProductID = id
		}
	}

	cover, _ := c.FormFile("cover")
	req.Cover = cover

	if form, err := c.MultipartForm(); err == nil && form.File != nil {
		if files, ok := form.File["gallery"]; ok {
			req.Gallery = files
		}
	}

	product, err := p.ProductService.CreateProduct(c.Request().Context(), uuid.MustParse(superAdminID), req)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to create product")
	}

	return response.Success(c, http.StatusCreated, "Product created successfully", product)
}

func (p *ProducHandlers) GetByIDProduct(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	product, err := p.ProductService.GetProductByID(c.Request().Context(), id)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get product")
	}

	return response.Success(c, http.StatusOK, "Product retrieved successfully", product)
}

func (p *ProducHandlers) GetAllProduct(c echo.Context) error {
	pageInt, limtiInt := utils.ParsePaginationParams(c, 10)
	search := c.QueryParam("search")

	products, total, err := p.ProductService.GetProductAll(c.Request().Context(), pageInt, limtiInt, search)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get product")
	}

	meta := utils.BuildPaginationMeta(c, pageInt, limtiInt, total)
	return response.PaginatedSuccess(c, http.StatusOK, "Products retrieved successfully", products, meta)
}

func (p *ProducHandlers) UpdateProduct(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	req := dto.UpdateProductRequest{ID: id}
	if v := c.FormValue("name"); v != "" {
		req.Name = &v
	}
	if v := c.FormValue("code"); v != "" {
		req.Code = &v
	}
	if v := c.FormValue("unit"); v != "" {
		req.Unit = &v
	}
	if v := c.FormValue("description"); v != "" {
		req.Description = &v
	}

	if v := c.FormValue("price"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			req.Price = &n
		}
	}
	if v := c.FormValue("type_product_id"); v != "" {
		if id2, err := uuid.Parse(v); err == nil {
			req.TypeProductID = &id2
		}
	}
	if c.FormValue("replace_gallery") == "true" {
		req.ReplaceGallery = true
	}

	req.NewCover, _ = c.FormFile("new_cover")
	if form, err := c.MultipartForm(); err == nil && form.File != nil {
		if files, ok := form.File["new_gallery"]; ok {
			req.NewGallery = files
		}
	}

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

	product, err := p.ProductService.UpdateProduct(c.Request().Context(), uuid.MustParse(superAdminID), req)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to update product")
	}

	return response.Success(c, http.StatusOK, "Product updated successfully", product)
}

func (p *ProducHandlers) DeleteProduct(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	if err := p.ProductService.Delete(c.Request().Context(), id); err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to delete product")
	}

	return response.Success(c, http.StatusOK, "Product deleted successfully", nil)
}
