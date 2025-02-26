package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func GainCORSSetting() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.DefaultCORSConfig)
}
