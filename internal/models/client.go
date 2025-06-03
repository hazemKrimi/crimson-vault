package models

import (
	"time"

	"gorm.io/gorm"
)

type Client struct {
	ID         uint32         `json:"id"`
	CreatedAt  time.Time      `json:"createAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	Name       string         `json:"name"`
	FiscalCode string         `json:"fiscalCode"`
	Address    string         `json:"address"`
	Zip        string         `json:"zip"`
	Country    string         `json:"country"`
	Phone      string         `json:"phone"`
	Email      string         `json:"email"`
}

type CreateClientRequestBody struct {
	Name       string         `json:"name" validate:"required"`
	FiscalCode string         `json:"fiscalCode"`
	Address    string         `json:"address" validate:"required"`
	Zip        string         `json:"zip" validate:"required"`
	Country    string         `json:"country" validate:"required"`
	Phone      string         `json:"phone" validate:"required,e164"`
	Email      string         `json:"email" validate:"required,email"`
}

type UpdateClientRequestBody struct {
	Name       string         `json:"name"`
	FiscalCode string         `json:"fiscalCode"`
	Address    string         `json:"address"`
	Zip        string         `json:"zip"`
	Country    string         `json:"country"`
	Phone      string         `json:"phone" validate:"omitempty,e164"`
	Email      string         `json:"email" validate:"omitempty,email"`
}

func (db *DB) MigrateClients() {
	db.instance.AutoMigrate(&Client{})
}

func (db *DB) CreateClient(body CreateClientRequestBody) Client {
	client := Client{Name: body.Name, Country: body.Country, Phone: body.Phone}

	db.instance.Create(&client)
	return client
}

func (db *DB) GetClients() ([]Client, error) {
	var clients []Client

	result := db.instance.Find(&clients)

	if result.Error != nil {
		return nil, result.Error
	}

	return clients, nil
}

func (db *DB) GetClient(id int, client *Client) error {
	result := db.instance.Where("id = ?", id).First(&client, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) UpdateClient(id int, body UpdateClientRequestBody, client *Client) error {
	result := db.instance.Where("id = ?", id).First(&client, id)

	if result.Error != nil {
		return result.Error
	}

	result = db.instance.Model(&client).Updates(Client{
		Name: body.Name,
		FiscalCode: body.FiscalCode,
		Address: body.Address,
		Zip: body.Zip,
		Country: body.Country,
		Phone: body.Phone,
		Email: body.Email,
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) DeleteClient(id int) error {
	result := db.instance.Delete(&Client{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
