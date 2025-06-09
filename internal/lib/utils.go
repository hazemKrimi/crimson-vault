package lib

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/hazemKrimi/crimson-vault/internal/types"
)

func GetConfigDirectory() (string, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	config, err := filepath.Abs(filepath.Join(home, DEFAULT_CONFIG_DIRECTORY))

	return config, nil
}

func SaveSession(session *sessions.Session, context echo.Context) error {
	if err := session.Save(context.Request(), context.Response()); err != nil {
		return context.String(http.StatusInternalServerError, "Unexpected error saving User session!")
	}

	return nil
}

func CreateSession(session *sessions.Session, context echo.Context, user *types.User) error {
	if err := uuid.Validate(user.SessionID); err != nil {
		return context.String(http.StatusInternalServerError, "Unexpected error saving User session!")
	}

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	}
	session.Values["id"] = user.ID
	session.Values["sessionId"] = user.SessionID
	session.Values["username"] = user.Username

	if err := SaveSession(session, context); err != nil {
		return context.String(http.StatusInternalServerError, "Unexpected error saving User session!")
	}

	return nil
}

func DeleteSession(session *sessions.Session, context echo.Context) error {
	session.Options.MaxAge = -1

	if err := SaveSession(session, context); err != nil {
		return context.String(http.StatusInternalServerError, "Unexpected error saving User session!")
	}

	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
