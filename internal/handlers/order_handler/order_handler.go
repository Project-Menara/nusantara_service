package orderhandler

import (
	"net/http"
	orderservice "nusantara_service/internal/data/services/order_service"
	"nusantara_service/internal/response"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	orderService orderservice.OrderService
}

func NewOrderHandler(orderService orderservice.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (o *OrderHandler) CreateOrder(c echo.Context) error {
	userToken := c.Get("user")
	if userToken == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	user, ok := userToken.(*jwt.Token)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	claims := user.Claims.(jwt.MapClaims)
	custID := claims["user_id"].(string)

	err := o.orderService.CreateOrder(c.Request().Context(), uuid.MustParse(custID))
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusBadRequest, "something went wrong", err.Error())
	}

	return response.Success(c, http.StatusOK, "Create Order Success", nil)
}
