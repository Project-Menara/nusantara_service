package handlers

import (
	"net/http"
	"nusantara_service/internal/data/services"
	dto "nusantara_service/internal/dto/request"
	customerresponse "nusantara_service/internal/dto/responses/customer_response"
	shopresponse "nusantara_service/internal/dto/responses/shop_response"
	"nusantara_service/internal/response"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CustomerAddressHandler struct {
	CustomerAddressService services.CustomerAddressService
}

func NewCustomerAddressHandler(service services.CustomerAddressService) *CustomerAddressHandler {
	return &CustomerAddressHandler{CustomerAddressService: service}
}

func (h *CustomerAddressHandler) CreateAddress(c echo.Context) error {
	userToken := c.Get("user")
	if userToken == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
	}
	user := userToken.(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uuid.MustParse(claims["user_id"].(string))

	var req dto.CreateAddressRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "failed to bind request", err.Error())
	}

	if err := h.CustomerAddressService.CreateAddres(c.Request().Context(), userID, req); err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to create address")
	}

	return response.Success(c, http.StatusOK, "Address created successfully", nil)
}

func (h *CustomerAddressHandler) GetAllAddress(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uuid.MustParse(claims["user_id"].(string))

	addrs, err := h.CustomerAddressService.GetAllAddress(c.Request().Context(), userID)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get addresses")
	}

	data := make([]customerresponse.CustomerAddressResponse, len(addrs))
	for i, ca := range addrs {
		data[i] = customerresponse.ToCustomerAddressResponse(*ca)
	}

	return response.Success(c, http.StatusOK, "Addresses retrieved successfully", data)
}

func (h *CustomerAddressHandler) GetByIdAddress(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uuid.MustParse(claims["user_id"].(string))

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	addr, err := h.CustomerAddressService.GetByIdAddress(c.Request().Context(), id, userID)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get address")
	}

	return response.Success(c, http.StatusOK, "Address retrieved successfully", customerresponse.ToCustomerAddressResponse(*addr))
}

func (h *CustomerAddressHandler) UpdateAddress(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uuid.MustParse(claims["user_id"].(string))

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	var req dto.UpdateAddressRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "failed to bind request", err.Error())
	}

	if err := h.CustomerAddressService.UpdateAddress(c.Request().Context(), id, userID, req); err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to update address")
	}

	return response.Success(c, http.StatusOK, "Address updated successfully", nil)
}

func (h *CustomerAddressHandler) DeleteAddress(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uuid.MustParse(claims["user_id"].(string))

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	if err := h.CustomerAddressService.Delete(c.Request().Context(), id, userID); err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to delete address")
	}

	return response.Success(c, http.StatusOK, "Address deleted successfully", nil)
}

func (h *CustomerAddressHandler) GetDefaultAddress(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uuid.MustParse(claims["user_id"].(string))

	addr, err := h.CustomerAddressService.GetDefaultAddress(c.Request().Context(), userID)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get default address")
	}

	return response.Success(c, http.StatusOK, "Default address retrieved successfully", customerresponse.ToCustomerAddressResponse(*addr))
}

func (h *CustomerAddressHandler) SetDefaultAddress(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uuid.MustParse(claims["user_id"].(string))

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	if err := h.CustomerAddressService.SetDefaultAddress(c.Request().Context(), userID, id); err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to set default address")
	}

	return response.Success(c, http.StatusOK, "Default address updated successfully", nil)
}

func (h *CustomerAddressHandler) GetNearbyShops(c echo.Context) error {
	ctx := c.Request().Context()
	latStr := c.QueryParam("lat")
	lngStr := c.QueryParam("lng")

	var lat, lng float64
	var err error

	if latStr != "" && lngStr != "" {
		// Pakai current location
		lat, err = strconv.ParseFloat(latStr, 64)
		if err != nil {
			return response.Error(c, http.StatusBadRequest, "invalid lat", err.Error())
		}
		lng, err = strconv.ParseFloat(lngStr, 64)
		if err != nil {
			return response.Error(c, http.StatusBadRequest, "invalid lng", err.Error())
		}
	} else {
		val := c.Get("user")
		if val == nil {
			return response.Error(c, http.StatusUnauthorized, "you must be login", "no user found in context")
		}

		user, ok := val.(*jwt.Token)
		if !ok || user == nil {
			return response.Error(c, http.StatusUnauthorized, "you must be login", "invalid token type")
		}

		claims, ok := user.Claims.(jwt.MapClaims)
		if !ok {
			return response.Error(c, http.StatusUnauthorized, "you must be login", "invalid claims")
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok || userIDStr == "" {
			return response.Error(c, http.StatusUnauthorized, "you must be login", "user_id missing in claims")
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return response.Error(c, http.StatusBadRequest, "invalid user_id", err.Error())
		}
		addr, err := h.CustomerAddressService.GetDefaultAddress(ctx, userID)
		if err != nil {
			if customErr, ok := response.AsCustomErr(err); ok {
				return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
			}
			return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get default address")
		}
		if addr == nil {
			return response.Error(c, http.StatusNotFound, "no default address found", "please set default address")
		}
		lat, lng = addr.Lat, addr.Lng
	}

	// Radius fix: 10km
	const maxDistance = 10.0

	shops, distMap, err := h.CustomerAddressService.GetNearbyShops(ctx, lat, lng, maxDistance)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get nearby shops")
	}

	// Mapping ke response
	responses := []shopresponse.ShopNearbyResponse{}
	for _, shop := range shops {
		distance := distMap[shop.ID]
		res := shopresponse.ToShopNearbyResponse(*shop, distance)
		responses = append(responses, res)
	}

	return response.Success(c, http.StatusOK, "Nearby shops retrieved successfully", responses)
}
