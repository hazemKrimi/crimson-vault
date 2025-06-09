package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

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
