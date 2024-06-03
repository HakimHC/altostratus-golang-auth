package routes

import (
	"github.com/HakimHC/altostratus-golang-auth/controllers"
	"github.com/labstack/echo/v4"
)

func AuthRoute(e *echo.Echo) {
	api := e.Group("/api/v1/auth", serverHeader)

	api.POST("/register", controllers.Register)
	api.POST("/login", controllers.Login)

	api.GET("/health", controllers.HealthCheck)
}

func serverHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("X-Version", "1.0.0")
		return next(c)
	}
}
