package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/hazemKrimi/crimson-vault/internal/lib"
	"github.com/hazemKrimi/crimson-vault/internal/types"
)

func (api *API) LoginHandler(context echo.Context) error {
	var body types.LoginRequestBody

	if err := context.Bind(&body); err != nil {
		log.Println(fmt.Sprintf("Error logging User in: %v.", err))
		return context.String(http.StatusBadRequest, "Invalid JSON!")
	}

	if err := context.Validate(body); err != nil {
		return err
	}

	var user types.User

	if err := api.db.GetUserByUsername(body.Username, &user); err != nil {
		return context.String(http.StatusNotFound, "User not found!")
	}

	if match := lib.CheckPasswordHash(body.Password, user.Password); !match {
		return context.String(http.StatusBadRequest, "Invalid credentials!")
	}

	sess, err := session.Get("session", context)

	if err != nil {
		log.Println(fmt.Sprintf("Error creating User session: %v.", err))
		return context.String(http.StatusInternalServerError, "Unexpected error creating User session!")
	}

	if err := api.db.UpdateUserSessionID(&user); err != nil {
		log.Println(fmt.Sprintf("Error creating User session: %v.", err))
		return context.String(http.StatusInternalServerError, "Unexpected error creating User session!")
	}

	if err := lib.CreateSession(sess, context, &user); err != nil {
		log.Println(fmt.Sprintf("Error creating User session: %v.", err))
		return context.String(http.StatusInternalServerError, "Unexpected error creating User session!")
	}

	log.Println(fmt.Sprintf("User with ID %s logged in.", user.ID))
	return context.JSON(http.StatusOK, user)
}

func (api *API) LogoutHandler(context echo.Context) error {
	sessionId, ok := context.Get("sessionId").(string)

	if !ok {
		return context.String(http.StatusInternalServerError, "Unexpected error deleting User session!")
	}

	if err := api.db.DeleteUserSessionID(sessionId); err != nil {
		log.Println(fmt.Sprintf("Error deleting User session: %v.", err))
		return context.String(http.StatusInternalServerError, "Unexpected error deleting User session!")
	}

	sess, err := session.Get("session", context)

	if err != nil {
		log.Println(fmt.Sprintf("Error deleting User session: %v.", err))
		return context.String(http.StatusInternalServerError, "Unexpected error deleting User session!")
	}

	if err := lib.DeleteSession(sess, context); err != nil {
		log.Println(fmt.Sprintf("Error deleting User session: %v.", err))
		return context.String(http.StatusInternalServerError, "Unexpected error deleting User session!")
	}

	log.Println(fmt.Sprintf("User with SessionID %s logged out.", sessionId))
	return context.String(http.StatusOK, "Logged out successfully!")
}
