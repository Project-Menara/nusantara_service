package handlers

import (
	"net/http"
	"nusantara_service/internal/data/services"
	dto "nusantara_service/internal/dto/request"
	cashierresponse "nusantara_service/internal/dto/responses/cashier_response"
	"nusantara_service/internal/response"
	"nusantara_service/internal/utils"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CashierHandlers struct {
	CashierService services.CashierService
}

func NewCashierHandler(service services.CashierService) *CashierHandlers {
	return &CashierHandlers{CashierService: service}
}

func (cs *CashierHandlers) CreateCashier(c echo.Context) error {
	var req dto.CreateCashierRequest
	req.Name = c.FormValue("name")
	req.Email = c.FormValue("email")
	req.Username = c.FormValue("username")
	req.Password = c.FormValue("password")
	if v := c.FormValue("status"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			req.Status = n
		}
	}
	image, err := c.FormFile("image")
	if err == nil {
		req.Photo = image
	}

	err = cs.CashierService.CreateCashier(c.Request().Context(), req)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to create cashier")
	}

	return response.Success(c, http.StatusOK, "Cashier Created", nil)
}

func (cs *CashierHandlers) GetCashierAll(c echo.Context) error {
	pageInt, limitInt := utils.ParsePaginationParams(c, 10)
	search := c.QueryParam("search")

	cashier, total, err := cs.CashierService.GetCashierAll(c.Request().Context(), pageInt, limitInt, search)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get cashier")
	}

	meta := utils.BuildPaginationMeta(c, pageInt, limitInt, total)

	data := make([]cashierresponse.CashierResponse, len(cashier))
	for i, u := range cashier {
		data[i] = cashierresponse.ToCashierResponse(*u)
	}
	return response.PaginatedSuccess(c, http.StatusOK, "Cashier retrieved successfully", data, meta)
}

func (cs *CashierHandlers) GetCashierById(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	cashier, err := cs.CashierService.GetCashierById(c.Request().Context(), id)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get cashier")
	}

	return response.Success(c, http.StatusOK, "Cashier Retrieved Successfully", cashierresponse.ToCashierResponse(*cashier))
}

func (cs *CashierHandlers) UpdateCashier(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	var req dto.UpdateCashierRequest
	if name := c.FormValue("name"); name != "" {
		req.Name = &name
	}
	if username := c.FormValue("username"); username != "" {
		req.Username = &username
	}
	if email := c.FormValue("email"); email != "" {
		req.Email = &email
	}
	if status := c.FormValue("status"); status != "" {
		if n, err := strconv.Atoi(status); err == nil {
			req.Status = &n
		}
	}
	image, err := c.FormFile("image")
	if err == nil {
		req.Photo = image
	}

	err = cs.CashierService.UpdateCashier(c.Request().Context(), id, req)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to update cashier")
	}

	return response.Success(c, http.StatusOK, "Cashier Update Successfully", nil)
}

func (cs *CashierHandlers) DeleteCashier(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	if err := cs.CashierService.Delete(c.Request().Context(), id); err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to delete cashier")
	}

	return response.Success(c, http.StatusOK, "Cashier Deleted Successfully", nil)
}
