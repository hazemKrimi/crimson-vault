package api

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
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

	api.instance.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	// TODO: Change and store the secret separately when finilizing v1.
	api.instance.Use(session.Middleware(sessions.NewCookieStore([]byte("SECRET"))))
	
	api.ClientRoutes()
	api.UserRoutes()
	api.AuthRoutes()
	
	api.instance.Logger.Fatal(api.instance.Start(fmt.Sprintf(":%d", lib.DEFAULT_PORT)))
}
