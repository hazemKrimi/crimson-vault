package models

import (
	"github.com/google/uuid"

	"github.com/hazemKrimi/crimson-vault/internal/types"
)

func (db *DB) MigrateClients() {
	db.instance.AutoMigrate(&types.Client{})
}

func (db *DB) CreateClient(userId uuid.UUID, body types.CreateClientRequestBody) types.Client {
	client := types.Client{
		ID:         uuid.New().String(),
		UserID:     userId.String(),
		Name:       body.Name,
		FiscalCode: body.FiscalCode,
		Address:    body.Address,
		Zip:        body.Zip,
		Country:    body.Country,
		Phone:      body.Phone,
		Email:      body.Email,
	}

	db.instance.Create(&client)
	return client
}

func (db *DB) GetClients(userId uuid.UUID) ([]types.Client, error) {
	var clients []types.Client

	result := db.instance.Where("user_id = ?", userId).Find(&clients)

	if result.Error != nil {
		return nil, result.Error
	}

	return clients, nil
}

func (db *DB) GetClientById(userId, id uuid.UUID, client *types.Client) error {
	result := db.instance.Where("user_id = ?", userId).Where("id = ?", id).First(client)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) UpdateClient(userId, id uuid.UUID, body types.UpdateClientRequestBody, client *types.Client) error {
	result := db.instance.Where("user_id = ?", userId).Where("id = ?", id).First(client)

	if result.Error != nil {
		return result.Error
	}

	result = db.instance.Model(client).Updates(types.Client{
		Name:       body.Name,
		FiscalCode: body.FiscalCode,
		Address:    body.Address,
		Zip:        body.Zip,
		Country:    body.Country,
		Phone:      body.Phone,
		Email:      body.Email,
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) DeleteClient(userId, id uuid.UUID) error {
	result := db.instance.Where("user_id = ?", userId).Where("id = ?", id).Delete(&types.Client{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}
