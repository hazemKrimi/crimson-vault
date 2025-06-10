package api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/hazemKrimi/crimson-vault/internal/lib"
	"github.com/hazemKrimi/crimson-vault/internal/types"
)

func (api *API) LoginHandler(context echo.Context) error {
	var body types.LoginRequestBody

	if err := context.Bind(&body); err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"Invalid JSON!"}}
	}

	if err := context.Validate(body); err != nil {
		return err
	}

	var user types.User

	if err := api.db.GetUserByUsername(body.Username, &user); err != nil {
		return types.Error{Code: http.StatusNotFound, Messages: []string{"User not found!"}}
	}

	if match := lib.CheckPasswordHash(body.Password, user.Password); !match {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"Invalid credentials!"}}
	}

	sess, err := session.Get("session", context)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error creating User session!"}}
	}

	if err := api.db.UpdateUserSessionID(&user); err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error creating User session!"}}
	}

	if err := lib.CreateSession(sess, context, &user); err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error creating User session!"}}
	}

	return context.JSON(http.StatusOK, user)
}

func (api *API) LogoutHandler(context echo.Context) error {
	sessionId, ok := context.Get("sessionId").(string)

	if !ok {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Session ID not found after authorization."), Messages: []string{"Unexpected error deleting User session!"}}
	}

	if err := api.db.DeleteUserSessionID(sessionId); err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error deleting User session!"}}
	}

	sess, err := session.Get("session", context)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error deleting User session!"}}
	}

	if err := lib.DeleteSession(sess, context); err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error deleting User session!"}}
	}

	return context.JSON(http.StatusOK, map[string]string{"message":"Logged out successfully!"})
}
