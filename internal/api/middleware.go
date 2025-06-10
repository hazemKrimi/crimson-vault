package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/hazemKrimi/crimson-vault/internal/types"
)

func (api *API) AuthSessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		sess, err := session.Get("session", context)

		if err != nil || sess.IsNew {
			return context.String(http.StatusUnauthorized, "User not authenticated!")
		}

		id, ok := sess.Values["sessionId"].(string)

		if !ok || id == "" {
			return context.String(http.StatusUnauthorized, "User not authenticated!")
		}

		sessionId, err := uuid.Parse(id)

		if err != nil {
			return context.String(http.StatusUnauthorized, "User not authenticated!")
		}

		var user types.User

		if err := api.db.GetUserBySessionId(sessionId, &user); err != nil {
			return context.String(http.StatusUnauthorized, "User not authenticated!")
		}

		context.Set("id", sess.Values["id"])
		context.Set("sessionId", sess.Values["sessionId"])
		context.Set("name", sess.Values["name"])

		return next(context)
	}
}

func (api *API) LoggerMiddleware() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true,
		LogMethod:   true,
		LogValuesFunc: func(logContext echo.Context, values middleware.RequestLoggerValues) error {
			if values.Error == nil {
				api.Logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", values.URI),
					slog.Int("status", values.Status),
				)
			} else {
				api.Logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", values.URI),
					slog.Int("status", values.Status),
					slog.String("err", values.Error.Error()),
				)
			}
			return nil
		},
	})
}
