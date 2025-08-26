package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HealthCheckHandler(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func EnqueueBarcodeHandler(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
