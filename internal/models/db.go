package models

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	instance *gorm.DB
}

func (wrapper *DB) Connect() {
	db, err := gorm.Open(sqlite.Open("crimson_vault.db"), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	wrapper.instance = db
}
