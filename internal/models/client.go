package models

import (
	"github.com/hazemKrimi/crimson-vault/internal/types"
)

func (db *DB) MigrateClients() {
	db.instance.AutoMigrate(&types.Client{})
}

func (db *DB) CreateClient(body types.CreateClientRequestBody) types.Client {
	client := types.Client{
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

func (db *DB) GetClients() ([]types.Client, error) {
	var clients []types.Client

	result := db.instance.Find(&clients)

	if result.Error != nil {
		return nil, result.Error
	}

	return clients, nil
}

func (db *DB) GetClient(id uint32, client *types.Client) error {
	result := db.instance.Where("id = ?", id).First(client, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) UpdateClient(id uint32, body types.UpdateClientRequestBody, client *types.Client) error {
	result := db.instance.Where("id = ?", id).First(client, id)

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

func (db *DB) DeleteClient(id uint32) error {
	result := db.instance.Delete(&types.Client{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
