package handlers

import (
	"net/http"
	"nusantara_service/internal/data/services"
	"nusantara_service/internal/dto"
	"nusantara_service/internal/response"
	"nusantara_service/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type VoucherHandler struct {
	VoucherService services.VoucherService
}

func NewVoucherHandler(service services.VoucherService) *VoucherHandler {
	return &VoucherHandler{VoucherService: service}
}

func (v *VoucherHandler) CreateVoucher(c echo.Context) error {
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

	var req dto.CreateVoucherRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error(), "invalid request")
	}

	newVoucher, err := v.VoucherService.CreateVoucher(c.Request().Context(), superAdminID, req)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to create voucher")
	}

	return response.Success(c, http.StatusCreated, "create voucher success", newVoucher)
}

func (v *VoucherHandler) GetAllVoucher(c echo.Context) error {
	pageInt, limtiInt := utils.ParsePaginationParams(c, 10)
	search := c.QueryParam("search")

	vouchers, total, err := v.VoucherService.GetAllVoucher(c.Request().Context(), pageInt, limtiInt, search)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get voucher")
	}

	meta := utils.BuildPaginationMeta(c, pageInt, limtiInt, total)

	return response.PaginatedSuccess(c, http.StatusOK, "Get All Voucher Success", vouchers, meta)
}

func (v *VoucherHandler) GetByIdVoucher(c echo.Context) error {
	voucherId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	voucher, err := v.VoucherService.GetByIdVoucher(c.Request().Context(), voucherId)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get voucher")
	}

	return response.Success(c, http.StatusOK, "Get Voucher Success", voucher)
}

func (v *VoucherHandler) UpdateVoucher(c echo.Context) error {
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

	var req dto.UpdateVoucherRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error(), "invalid request")
	}
	voucherId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	updated, err := v.VoucherService.UpdateVoucher(c.Request().Context(), superAdminID, voucherId, req)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to update voucher")
	}

	return response.Success(c, http.StatusOK, "updated voucher success", updated)
}

func (v *VoucherHandler) DeleteVoucher(c echo.Context) error {
	voucherId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	if err := v.VoucherService.DeleteVoucher(c.Request().Context(), voucherId); err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to delete voucher")
	}

	return response.Success(c, http.StatusOK, "deleted success", nil)
}

func (v *VoucherHandler) UpdateStatusVoucher(c echo.Context) error {
	voucherId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	var req dto.UpdateStatusVoucherRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error(), "invalid request")
	}

	if err := v.VoucherService.UpdateStatusVoucher(c.Request().Context(), voucherId, req); err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to update status voucher")
	}

	return response.Success(c, http.StatusOK, "update status voucher success", nil)
}

func (v *VoucherHandler) GetAllVoucherCustomer(c echo.Context) error {
	pageInt, limtiInt := utils.ParsePaginationParams(c, 10)

	vouchers, total, err := v.VoucherService.GetAllVoucherCustomer(c.Request().Context(), pageInt, limtiInt)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get voucher")
	}

	meta := utils.BuildPaginationMeta(c, pageInt, limtiInt, total)

	return response.PaginatedSuccess(c, http.StatusOK, "Get All Voucher Success", vouchers, meta)
}

func (v *VoucherHandler) GetByIdVoucherCustomer(c echo.Context) error {
	voucherId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	voucher, err := v.VoucherService.GetByIdVoucherCustomer(c.Request().Context(), voucherId)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get voucher")
	}

	return response.Success(c, http.StatusOK, "Get Voucher Success", voucher)
}
