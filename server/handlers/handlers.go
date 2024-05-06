package handlers

import (
	//"encoding/json"
	//"fmt"
	//"log"
	"net/http"
	//"strings"

	"github.com/dchote/survey-engine/config"

	"github.com/labstack/echo"
)

// JSONResponse is a standard map interface
type JSONResponse map[string]interface{}

// Health simple response
func Health() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, JSONResponse{
			"status": "OK",
		})
	}
}

// Config returns our parsed config
func Config() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, config.Config.Survey)
	}
}
