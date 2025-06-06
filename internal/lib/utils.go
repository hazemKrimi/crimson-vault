package lib

import (
	"os"
	"path/filepath"

	"github.com/gorilla/sessions"
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

func ConstructSession(session *sessions.Session, user types.User) {
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	}
	session.Values["id"] = user.ID
}
