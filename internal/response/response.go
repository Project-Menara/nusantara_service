package response

import "github.com/labstack/echo/v4"

type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Error      interface{} `json:"error,omitempty"`
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
