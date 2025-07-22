package response

import "github.com/labstack/echo/v4"

type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Error      interface{} `json:"error,omitempty"`
}

type PaginationMeta struct {
	CurrentPage int    `json:"current_page"`
	PerPage     int    `json:"per_page"`
	TotalData   int    `json:"total_data"`
	TotalPages  int    `json:"total_pages"`
	NextPageURL string `json:"next_page_url,omitempty"`
	PrevPageURL string `json:"prev_page_url,omitempty"`
}

type PaginationResponse struct {
	StatusCode int            `json:"status_code"`
	Message    string         `json:"message"`
	Data       interface{}    `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// Success response
func Success(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}

func Error(c echo.Context, statusCode int, message string, err interface{}) error {
	return c.JSON(statusCode, Response{
		StatusCode: statusCode,
		Message:    message,
		Error:      err,
	})
}

func PaginatedSuccess(c echo.Context, statusCode int, message string, data interface{}, pagination PaginationMeta) error {
	return c.JSON(statusCode, PaginationResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}
