package api

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/hazemKrimi/crimson-vault/internal/lib"
	"github.com/hazemKrimi/crimson-vault/internal/models"
)

type API struct {
	ConfigDirectory string
	instance  *echo.Echo
	db        *models.DB
}

func (api *API) Initialize() {
	validator := validator.New(validator.WithRequiredStructEnabled())
	validator.RegisterValidation("password", PasswordValidator)

	db := &models.DB{}
	ech := echo.New()
	ech.Validator = &CustomValidator{validator: validator}

	db.Connect(api.ConfigDirectory)
	db.MigrateClients()
	db.MigrateUsers()

	api.instance = ech
	api.db = db

	api.ClientRoutes()
	api.UserRoutes()
	api.instance.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	api.instance.Logger.Fatal(api.instance.Start(fmt.Sprintf(":%d", lib.DEFAULT_PORT)))
}
