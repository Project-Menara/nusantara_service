package handlers

import (
	"encoding/json"
	"net/http"
	"nusantara_service/internal/data/services"
	"nusantara_service/internal/domain/entities"
	dto "nusantara_service/internal/dto/request"
	eventresponse "nusantara_service/internal/dto/responses/event_response"
	"nusantara_service/internal/response"
	"nusantara_service/internal/utils"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type EventHandler struct {
	service services.EventService
}

func NewEventHandler(service services.EventService) *EventHandler {
	return &EventHandler{service: service}
}

func (e *EventHandler) CreateEvent(c echo.Context) error {
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

	var req dto.CreateEventRequest
	req.Name = c.FormValue("name")
	req.Type = c.FormValue("type_event")
	if v := c.FormValue("status"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			req.Status = n
		}
	}

	const dateLayout = time.RFC3339
	if startAtStr := c.FormValue("start_date"); startAtStr != "" {
		startDate, err := time.Parse(dateLayout, startAtStr)
		if err != nil {
			return response.Error(c, http.StatusBadRequest, "Invalid start_date format. Expected format: "+dateLayout, err.Error())
		}
		req.StartDate = startDate
	} else {
		// Usecase requires StartDate, so handle missing value
		return response.Error(c, http.StatusBadRequest, "start at is required", nil)
	}

	if endAtStr := c.FormValue("end_date"); endAtStr != "" {
		endDate, err := time.Parse(dateLayout, endAtStr)
		if err != nil {
			return response.Error(c, http.StatusBadRequest, "Invalid end_date format. Expected format: "+dateLayout, err.Error())
		}
		req.EndDate = endDate
	} else {
		return response.Error(c, http.StatusBadRequest, "end at is required", nil)
	}

	if cover, err := c.FormFile("cover"); err == nil {
		req.Cover = cover
	}

	switch req.Type {
	case string(entities.EventTypeDiskon):
		if eventProducts := c.FormValue("event_products"); eventProducts != "" {
			var parsed []dto.AddEventProduct
			if err := json.Unmarshal([]byte(eventProducts), &parsed); err == nil {
				req.EventProducts = parsed
			}
		}
	case string(entities.EventTypeBundle):
		if eventBundleBuys := c.FormValue("event_bundle_buys"); eventBundleBuys != "" {
			var parsed []dto.AddEventBundleBuy
			if err := json.Unmarshal([]byte(eventBundleBuys), &parsed); err != nil {
				return response.Error(c, http.StatusBadRequest, "invalid event_bundle_buys format, must be JSON array", err.Error())
			}
			req.EventBundleBuys = parsed
		}

		if eventBundleRewards := c.FormValue("event_bundle_rewards"); eventBundleRewards != "" {
			var parsed []dto.AddEventBundleReward
			if err := json.Unmarshal([]byte(eventBundleRewards), &parsed); err != nil {
				return response.Error(c, http.StatusBadRequest, "invalid event_bundle_rewards format, must be JSON array", err.Error())
			}
			req.EventBundleRewards = parsed
		}
	}

	err := e.service.CreateEvent(c.Request().Context(), uuid.MustParse(superAdminID), req)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to create event")
	}

	return response.Success(c, http.StatusCreated, "created event success", nil)

}

func (e *EventHandler) GetAllEvents(c echo.Context) error {
	pageInt, limitInt := utils.ParsePaginationParams(c, 10)
	search := c.QueryParam("search")

	events, total, err := e.service.GetAllEvents(c.Request().Context(), pageInt, limitInt, search)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get events")
	}

	meta := utils.BuildPaginationMeta(c, pageInt, limitInt, total)
	data := make([]eventresponse.EventResponse, len(events))
	for i, e := range events {
		data[i] = eventresponse.ToEventResponse(*e)
	}

	return response.PaginatedSuccess(c, http.StatusOK, "Events Retrieved Successfully", data, meta)
}

func (e *EventHandler) GetEventById(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}
	event, err := e.service.GetEventById(c.Request().Context(), id)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get event")
	}
	return response.Success(c, http.StatusOK, "Event Retrieved Successfully", eventresponse.ToEventResponse(*event))
}

func (e *EventHandler) UpdateEvent(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	var req dto.UpdateEventRequest
	req.Name = c.FormValue("name")
	req.Type = c.FormValue("type_event")

	const dateLayout = time.RFC3339
	if startAtStr := c.FormValue("start_date"); startAtStr != "" {
		startDate, err := time.Parse(dateLayout, startAtStr)
		if err != nil {
			return response.Error(c, http.StatusBadRequest, "Invalid start_date format. Expected format: "+dateLayout, err.Error())
		}
		req.StartDate = startDate
	}

	if endAtStr := c.FormValue("end_date"); endAtStr != "" {
		endDate, err := time.Parse(dateLayout, endAtStr)
		if err != nil {
			return response.Error(c, http.StatusBadRequest, "Invalid end_date format. Expected format: "+dateLayout, err.Error())
		}
		req.EndDate = endDate
	}

	if newCover, err := c.FormFile("new_cover"); err == nil {
		req.NewCover = newCover
	}

	switch req.Type {
	case string(entities.EventTypeDiskon):
		if eventProducts := c.FormValue("event_products"); eventProducts != "" {
			var parsed []dto.AddEventProduct
			if err := json.Unmarshal([]byte(eventProducts), &parsed); err == nil {
				req.EventProducts = parsed
			}
		}
	case string(entities.EventTypeBundle):
		if eventBundleBuys := c.FormValue("event_bundle_buys"); eventBundleBuys != "" {
			var parsed []dto.AddEventBundleBuy
			if err := json.Unmarshal([]byte(eventBundleBuys), &parsed); err != nil {
				return response.Error(c, http.StatusBadRequest, "invalid event_bundle_buys format, must be JSON array", err.Error())
			}
			req.EventBundleBuys = parsed
		}

		if eventBundleRewards := c.FormValue("event_bundle_rewards"); eventBundleRewards != "" {
			var parsed []dto.AddEventBundleReward
			if err := json.Unmarshal([]byte(eventBundleRewards), &parsed); err != nil {
				return response.Error(c, http.StatusBadRequest, "invalid event_bundle_rewards format, must be JSON array", err.Error())
			}
			req.EventBundleRewards = parsed
		}
	}

	err = e.service.UpdateEvent(c.Request().Context(), uuid.Nil, id, req)
	if err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to update event")
	}

	return response.Success(c, http.StatusOK, "event upadted successfully", nil)
}

func (e *EventHandler) DeleteEvent(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	if err := e.service.DeleteEvent(c.Request().Context(), id); err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to delete event")
	}

	return response.Success(c, http.StatusOK, "Event deleted successfully", nil)
}

func (e *EventHandler) UpdateStatusEvent(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid uuid", err.Error())
	}

	var req dto.UpdateStatusEventRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "failed to bind request", err.Error())
	}

	if err := e.service.UpdateStatusEvent(c.Request().Context(), id, req); err != nil {
		if customErr, ok := response.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to update status event")
	}

	return response.Success(c, http.StatusOK, "update status success", nil)
}
