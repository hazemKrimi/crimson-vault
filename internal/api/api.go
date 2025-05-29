package api

import "github.com/hazemKrimi/crimson-vault/internal/models"

type APIWrapper struct {
	dbWrapper *models.DBWrapper
}

func (api *APIWrapper) Initialize() {
	wrapper := models.DBWrapper{}

	wrapper.Connect()
	wrapper.MigrateClients()

	api.dbWrapper = &wrapper;
}
