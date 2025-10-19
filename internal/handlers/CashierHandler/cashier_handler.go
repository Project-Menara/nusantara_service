package cashierhandler

import (
	"net/http"
	cashierservice "nusantara_service/internal/data/services/CashierService"
	"nusantara_service/internal/response"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CashierHandler struct {
	service cashierservice.CashierService
}

func NewCashierHandler(service cashierservice.CashierService) *CashierHandler {
	return &CashierHandler{service: service}
}

func (cs *CashierHandler) GetAllShopName(c echo.Context) error {
	userToken := c.Get("user")
	if userToken == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	user, ok := userToken.(*jwt.Token)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	claims := user.Claims.(jwt.MapClaims)
	cashierID := claims["user_id"].(string)
	shopNames, err := cs.service.GetAllNameShop(c.Request().Context(), uuid.MustParse(cashierID))
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get shop")
	}

	return response.Success(c, http.StatusOK, "Get Shop Names Success", shopNames)
}

func (cs *CashierHandler) GetDetailShopCashier(c echo.Context) error {
	userToken := c.Get("user")
	if userToken == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	user, ok := userToken.(*jwt.Token)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	claims := user.Claims.(jwt.MapClaims)
	cashierID := claims["user_id"].(string)
	shopId, err := uuid.Parse(c.Param("shop_id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}
	shopCashier, err := cs.service.GetDetailShopCashier(c.Request().Context(), uuid.MustParse(cashierID), shopId)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get shop")
	}

	return response.Success(c, http.StatusOK, "Get Shop Cashier Success", shopCashier)
}

func (cs *CashierHandler) GetShopProductCashier(c echo.Context) error {
	userToken := c.Get("user")
	if userToken == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	user, ok := userToken.(*jwt.Token)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	claims := user.Claims.(jwt.MapClaims)
	cashierID := claims["user_id"].(string)
	shopId, err := uuid.Parse(c.Param("shop_id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}
	shopProductCashier, err := cs.service.GetDetailShopProduct(c.Request().Context(), uuid.MustParse(cashierID), shopId)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get shop")
	}

	return response.Success(c, http.StatusOK, "Get Cashier Shop Product Success", shopProductCashier)
}
