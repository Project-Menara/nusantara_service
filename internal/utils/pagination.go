package utils

import (
	"fmt"
	"nusantara_service/internal/response"
	"strconv"

	"github.com/labstack/echo/v4"
)

func ParsePaginationParams(c echo.Context, defaultLimit int) (page int, limit int) {
	page = 1
	limit = defaultLimit

	pageParam := c.QueryParam("page")
	limitParam := c.QueryParam("limit")

	if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
		page = p
	}

	if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
		limit = l
	}

	return page, limit
}

func BuildPaginationMeta(c echo.Context, page, limit, total int) response.PaginationMeta {
	totalPages := (total + limit - 1) / limit
	baseURL := c.Scheme() + "://" + c.Request().Host + c.Path()

	var nextPage, prevPage string

	if page < totalPages {
		nextPage = fmt.Sprintf("%s?page=%d&limit=%d", baseURL, page+1, limit)
	}
	if page > 1 {
		prevPage = fmt.Sprintf("%s?page=%d&limit=%d", baseURL, page-1, limit)
	}

	return response.PaginationMeta{
		CurrentPage: page,
		PerPage:     limit,
		TotalData:   total,
		TotalPages:  totalPages,
		NextPageURL: nextPage,
		PrevPageURL: prevPage,
	}
}
