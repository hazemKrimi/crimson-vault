package api

import (
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/hazemKrimi/crimson-vault/internal/lib"
	"github.com/hazemKrimi/crimson-vault/internal/models"
	"github.com/hazemKrimi/crimson-vault/internal/types"
)

type API struct {
	ConfigDirectory string
	instance        *echo.Echo
	db              *models.DB
	Logger          *slog.Logger
}

func (api *API) Initialize() {
	validator := validator.New(validator.WithRequiredStructEnabled())
	validator.RegisterValidation("password", PasswordValidator)

	db := &models.DB{}
	ech := echo.New()
	ech.Validator = &CustomValidator{validator: validator}
	ech.HTTPErrorHandler = func(err error, context echo.Context) {
		if context.Response().Committed {
			return
		}

		custom, ok := err.(types.Error)

		if ok {
			context.JSONPretty(custom.Code, map[string][]string{"errors": custom.Messages}, " ")
		}

		ech.DefaultHTTPErrorHandler(err, context)
	}

	db.Connect(api.ConfigDirectory)
	db.MigrateClients()
	db.MigrateUsers()

	api.instance = ech
	api.db = db

	api.instance.Use(api.LoggerMiddleware())
	// TODO: Update with appropriate origins when finishing v1
	api.instance.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	// TODO: Change and store the secret separately when finilizing v1.
	api.instance.Use(session.Middleware(sessions.NewCookieStore([]byte("SECRET"))))
	api.instance.Pre(middleware.AddTrailingSlash())

	api.ClientRoutes()
	api.UserRoutes()
	api.AuthRoutes()

	api.instance.Logger.Fatal(api.instance.Start(fmt.Sprintf(":%d", lib.DEFAULT_PORT)))
}
