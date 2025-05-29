package models

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBWrapper struct {
	db *gorm.DB
}

func (wrapper *DBWrapper) Connect() {
	db, err := gorm.Open(sqlite.Open("crimson_vault.db"), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	wrapper.db = db
}
