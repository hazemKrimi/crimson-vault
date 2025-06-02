package api

import (
	"github.com/hazemKrimi/crimson-vault/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type API struct {
	instance *echo.Echo
	db       *models.DB
}

func (api *API) Initialize() {
	db := &models.DB{}
	ech := echo.New()

	db.Connect()
	db.MigrateClients()

	api.instance = ech
	api.db = db

	api.ClientRoutes()
	api.instance.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	api.instance.Logger.Fatal(api.instance.Start(":5000"))
}
