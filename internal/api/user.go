package api

import (
	"fmt"
	"io"
	"log"
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
		log.Println(fmt.Sprintf("Error creating User: %v.", err))
		return context.String(http.StatusBadRequest, "Invalid JSON!")
	}

	if err := context.Validate(body); err != nil {
		return err
	}

	user, err := api.db.CreateUser(body)

	if err != nil {
		log.Println(fmt.Sprintf("Error creating User: %v.", err))
		return context.String(http.StatusInternalServerError, "Unexpected error creating User!")
	}

	sess, err := session.Get("session", context)

	if err != nil {
		log.Println(fmt.Sprintf("Error creating User session: %v.", err))
		return context.String(http.StatusInternalServerError, "Unexpected error creating User session!")
	}

	if err := lib.CreateSession(sess, context, &user); err != nil {
		log.Println(fmt.Sprintf("Error creating User session: %v.", err))
		return context.String(http.StatusInternalServerError, "Unexpected error creating User session!")
	}

	log.Println(fmt.Sprintf("User created with ID %s.", user.ID))
	return context.JSON(http.StatusOK, user)
}

func (api *API) GetAllUsersHandler(context echo.Context) error {
	users, err := api.db.GetUsers()

	if err != nil {
		return context.String(http.StatusInternalServerError, "Unexpected error getting User!")
	}

	log.Println("Got all Users.")
	return context.JSON(http.StatusOK, users)
}

func (api *API) GetUserHandler(context echo.Context) error {
	userId, ok := context.Get("id").(string)

	if !ok {
		return context.String(http.StatusInternalServerError, "Unexpected error getting User!")
	}

	id, err := uuid.Parse(userId)

	if err != nil {
		return context.String(http.StatusInternalServerError, "Unexpected error getting User!")
	}

	var user types.User

	if err := api.db.GetUserById(id, &user); err != nil {
		return context.String(http.StatusNotFound, "User not found!")
	}

	log.Println(fmt.Sprintf("Got User with ID %s.", user.ID))
	return context.JSON(http.StatusOK, user)
}

func (api *API) UpdateUserHandler(context echo.Context) error {
	userId, ok := context.Get("id").(string)

	if !ok {
		return context.String(http.StatusInternalServerError, "Unexpected error updating User!")
	}

	id, err := uuid.Parse(userId)

	if err != nil {
		return context.String(http.StatusInternalServerError, "Unexpected error updating User!")
	}

	var body types.UpdateUserRequestBody

	if err := context.Bind(&body); err != nil {
		log.Println(fmt.Sprintf("Error updating User: %v.", err))
		return context.String(http.StatusBadRequest, "Invalid JSON!")
	}

	if err := context.Validate(body); err != nil {
		return err
	}

	var user types.User

	if err := api.db.UpdateUser(id, body, &user); err != nil {
		return context.String(http.StatusNotFound, "User not found!")
	}

	log.Println(fmt.Sprintf("Updated user with ID %s.", user.ID))
	return context.JSON(http.StatusOK, user)
}

func (api *API) UpdateUserSecurityCredentialsHandler(context echo.Context) error {
	userId, ok := context.Get("id").(string)

	if !ok {
		return context.String(http.StatusInternalServerError, "Unexpected error updating User security credentials!")
	}

	id, err := uuid.Parse(userId)

	if err != nil {
		return context.String(http.StatusInternalServerError, "Unexpected error updating User security credentials!")
	}

	var body types.UpdateUserSecurityCredentialsBody

	if err := context.Bind(&body); err != nil {
		log.Println(fmt.Sprintf("Error creating security details for User: %v.", err))
		return context.String(http.StatusBadRequest, "Invalid JSON!")
	}

	if err := context.Validate(body); err != nil {
		return err
	}

	var user types.User

	if err := api.db.UpdateUserSecurityCredentials(id, body, &user); err != nil {
		return context.String(http.StatusNotFound, "User not found!")
	}

	log.Println(fmt.Sprintf("Updated security details of user with ID %s.", user.ID))
	return context.JSON(http.StatusOK, user)
}

func (api *API) UpdateUserLogoHandler(context echo.Context) error {
	userId, ok := context.Get("id").(string)

	if !ok {
		return context.String(http.StatusInternalServerError, "Unexpected error updating User logo!")
	}

	id, err := uuid.Parse(userId)

	if err != nil {
		return context.String(http.StatusInternalServerError, "Unexpected error updating User logo!")
	}

	var user types.User

	if err := api.db.GetUserById(id, &user); err != nil {
		return context.String(http.StatusNotFound, "User not found!")
	}

	if user.Username == "" {
		return context.String(http.StatusBadRequest, "You have to add a username first for this User!")
	}

	file, err := context.FormFile("logo")

	if err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"No image has been uploaded!"}}
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
		return context.String(http.StatusBadRequest, "Invalid file type, only image files are allowed!")
	}

	src, err := file.Open()

	if err != nil {
		log.Println(fmt.Sprintf("Error updating logo for User: %v.", err))
		return context.String(http.StatusInternalServerError, "Unexpected error while updating logo for User!")
	}

	defer src.Close()

	data, err := io.ReadAll(src)

	if err != nil {
		return context.String(http.StatusInternalServerError, "Unexpected error while updating logo for User!")
	}

	filetype := http.DetectContentType(data)

	if !strings.HasPrefix(filetype, "image/") {
		return context.String(http.StatusBadRequest, "Uploaded file is not a valid image!")
	}

	if err := os.MkdirAll(filepath.Join(api.ConfigDirectory, user.Username), os.ModePerm); err != nil {
		log.Println(fmt.Sprintf("Error updating logo for User: %v.", err))
		return context.String(http.StatusInternalServerError, "Unexpected error while updating logo for User!")
	}

	path, err := filepath.Abs(filepath.Join(api.ConfigDirectory, user.Username, fmt.Sprintf("logo%s", ext)))

	if err != nil {
		log.Println(fmt.Sprintf("Error updating logo for User: %v.", err))
		return context.String(http.StatusInternalServerError, "Unexpected error while updating logo for User!")
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		log.Println(fmt.Sprintf("Error updating logo for User: %v.", err))
		return context.String(http.StatusInternalServerError, "Unexpected error while updating logo for User!")
	}

	if err := api.db.UpdateUserLogo(path, &user); err != nil {
		log.Println(fmt.Sprintf("Error updating logo for User: %v.", err))
		return context.String(http.StatusInternalServerError, "Unexpected error while updating logo for User!")
	}

	return context.JSON(http.StatusOK, user)
}

func (api *API) DeleteUserHandler(context echo.Context) error {
	userId, ok := context.Get("id").(string)

	if !ok {
		return context.String(http.StatusInternalServerError, "Unexpected error deleting User!")
	}

	id, err := uuid.Parse(userId)

	if err != nil {
		return context.String(http.StatusInternalServerError, "Unexpected error deleting User!")
	}

	if err := api.db.DeleteUser(id); err != nil {
		return context.String(http.StatusNotFound, "User not found!")
	}

	log.Println(fmt.Sprintf("Deleted User with ID %d.", id))
	return context.String(http.StatusOK, "User deleted successfully!")
}

func (api *API) DeleteUserLogoHandler(context echo.Context) error {
	userId, ok := context.Get("id").(string)

	if !ok {
		return context.String(http.StatusInternalServerError, "Unexpected error deleting User logo!")
	}

	id, err := uuid.Parse(userId)

	if err != nil {
		return context.String(http.StatusInternalServerError, "Unexpected error deleting User logo!")
	}

	var user types.User

	if err := api.db.GetUserById(id, &user); err != nil {
		return context.String(http.StatusNotFound, "User not found!")
	}

	os.Remove(user.Logo)

	if err := api.db.DeleteUserLogo(&user); err != nil {
		log.Println(fmt.Sprintf("Error deleting logo of User: %v.", err))
		return context.String(http.StatusInternalServerError, "Unexpected error deleting logo of User!")
	}

	log.Println(fmt.Sprintf("Deleted logo of User with ID %s.", user.ID))
	return context.String(http.StatusOK, "User logo deleted successfully!")
}
