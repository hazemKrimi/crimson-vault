package models

import (
	"log"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	instance *gorm.DB
}

func (wrapper *DB) Connect(configDir string) {
	db, err := gorm.Open(sqlite.Open(filepath.Join(configDir, "crimson_vault.db")), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	wrapper.instance = db
}
