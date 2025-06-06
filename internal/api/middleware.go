package api

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func SessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		sess, err := session.Get("session", context)

		if err != nil {
			return context.String(http.StatusUnauthorized, "User not authenticated!")
		}

		cookie, err := context.Cookie("session")

		if err != nil {
			return context.String(http.StatusUnauthorized, "User not authenticated!")
		}

		if sess.IsNew || cookie.Value == "" || sess.Values["id"] == "" {
			return context.String(http.StatusUnauthorized, "User not authenticated!")
		}

		context.Set("id", sess.Values["id"])
		return next(context)
	}
}
