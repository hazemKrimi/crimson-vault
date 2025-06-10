package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/hazemKrimi/crimson-vault/internal/lib"
	"github.com/hazemKrimi/crimson-vault/internal/types"
)

func (api *API) CreateUserHandler(context echo.Context) error {
	var body types.CreateUserRequestBody

	if err := context.Bind(&body); err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"Invalid JSON!"}}
	}

	if err := context.Validate(body); err != nil {
		return err
	}

	user, err := api.db.CreateUser(body)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error creating User!"}}
	}

	sess, err := session.Get("session", context)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error creating User session!"}}
	}

	if err := lib.CreateSession(sess, context, &user); err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error creating User session!"}}
	}

	return context.JSON(http.StatusOK, user)
}

func (api *API) GetAllUsersHandler(context echo.Context) error {
	users, err := api.db.GetUsers()

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error getting User!"}}
	}

	return context.JSON(http.StatusOK, users)
}

func (api *API) GetUserHandler(context echo.Context) error {
	userId, ok := context.Get("id").(string)

	if !ok {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Session ID not found after authorization."), Messages: []string{"Unexpected error getting User!"}}
	}

	id, err := uuid.Parse(userId)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error getting User!"}}
	}

	var user types.User

	if err := api.db.GetUserById(id, &user); err != nil {
		return types.Error{Code: http.StatusNotFound, Messages: []string{"User not found!"}}
	}

	return context.JSON(http.StatusOK, user)
}

func (api *API) UpdateUserHandler(context echo.Context) error {
	userId, ok := context.Get("id").(string)

	if !ok {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Session ID not found after authorization."), Messages: []string{"Unexpected error updating User!"}}
	}

	id, err := uuid.Parse(userId)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error updating User!"}}
	}

	var body types.UpdateUserRequestBody

	if err := context.Bind(&body); err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"Invalid JSON!"}}
	}

	empty := body == types.UpdateUserRequestBody{}

	if empty {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"You must update at least one field!"}}
	}

	if err := context.Validate(body); err != nil {
		return err
	}

	var user types.User

	if err := api.db.UpdateUser(id, body, &user); err != nil {
		return types.Error{Code: http.StatusNotFound, Messages: []string{"User not found!"}}
	}

	return context.JSON(http.StatusOK, user)
}

func (api *API) UpdateUserSecurityCredentialsHandler(context echo.Context) error {
	userId, ok := context.Get("id").(string)

	if !ok {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Session ID not found after authorization."), Messages: []string{"Unexpected error updating User security credentials!"}}
	}

	id, err := uuid.Parse(userId)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error updating User security credentials!"}}
	}

	var body types.UpdateUserSecurityCredentialsBody

	if err := context.Bind(&body); err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"Invalid JSON!"}}
	}

	if err := context.Validate(body); err != nil {
		return err
	}

	var user types.User

	if err := api.db.UpdateUserSecurityCredentials(id, body, &user); err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"Invalid JSON!"}}
	}

	return context.JSON(http.StatusOK, user)
}

func (api *API) UpdateUserLogoHandler(context echo.Context) error {
	userId, ok := context.Get("id").(string)

	if !ok {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Session ID not found after authorization."), Messages: []string{"Unexpected error updating User logo!"}}
	}

	id, err := uuid.Parse(userId)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error updating User logo!"}}
	}

	var user types.User

	if err := api.db.GetUserById(id, &user); err != nil {
		return types.Error{Code: http.StatusNotFound, Messages: []string{"User not found!"}}
	}

	if user.Username == "" {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"You have to add a username first for this User!"}}
	}

	file, err := context.FormFile("logo")

	if err != nil {
		return types.Error{Code: http.StatusUnsupportedMediaType, Messages: []string{"Image must be uploaded in form data!"}}
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".webp": true,
	}

	if !allowedExtensions[ext] {
		return types.Error{Code: http.StatusUnsupportedMediaType, Messages: []string{"Invalid file type, only image files are allowed!"}}
	}

	src, err := file.Open()

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error updating logo User logo!"}}
	}

	defer src.Close()

	data, err := io.ReadAll(src)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error updating logo User logo!"}}
	}

	filetype := http.DetectContentType(data)

	if !strings.HasPrefix(filetype, "image/") {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"Uploaded file is not a valid image!"}}
	}

	if err := os.MkdirAll(filepath.Join(api.ConfigDirectory, user.Username), os.ModePerm); err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error updating User logo!"}}
	}

	path, err := filepath.Abs(filepath.Join(api.ConfigDirectory, user.Username, fmt.Sprintf("logo%s", ext)))

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error updating User logo!"}}
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error updating User logo!"}}
	}

	if err := api.db.UpdateUserLogo(path, &user); err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error updating User logo!"}}
	}

	return context.JSON(http.StatusOK, user)
}

func (api *API) DeleteUserHandler(context echo.Context) error {
	userId, ok := context.Get("id").(string)

	if !ok {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Session ID not found after authorization."), Messages: []string{"Unexpected error updating User logo!"}}
	}

	id, err := uuid.Parse(userId)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error deleting User!"}}
	}

	if err := api.db.DeleteUser(id); err != nil {
		return types.Error{Code: http.StatusNotFound, Messages: []string{"User not found!"}}
	}

	return context.JSON(http.StatusOK, map[string]string{"message":"User deleted successfully!"})
}

func (api *API) DeleteUserLogoHandler(context echo.Context) error {
	userId, ok := context.Get("id").(string)

	if !ok {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Session ID not found after authorization."), Messages: []string{"Unexpected error deleting User logo!"}}
	}

	id, err := uuid.Parse(userId)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error deleting User logo!"}}
	}

	var user types.User

	if err := api.db.GetUserById(id, &user); err != nil {
		return types.Error{Code: http.StatusNotFound, Messages: []string{"User not found!"}}
	}

	os.Remove(user.Logo)

	if err := api.db.DeleteUserLogo(&user); err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error deleting User logo!"}}
	}

	return context.JSON(http.StatusOK, map[string]string{"message":"User logo deleted successfully!"})
}
