package main

import (
	"github.com/HakimHC/altostratus-golang-auth/middleware"
	"github.com/HakimHC/altostratus-golang-auth/models"
	"github.com/HakimHC/altostratus-golang-auth/routes"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Validator = &models.CustomValidator{Validator: validator.New()}
	e.Use(middleware.LoggerMiddleware)
	routes.AuthRoute(e)
	e.Logger.Fatal(e.Start(":80"))
}
