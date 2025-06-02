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

type CreateClientBody struct {
	Name       string         `json:"name"`
	FiscalCode string         `json:"fiscalCode"`
	Address    string         `json:"address"`
	Zip        string         `json:"zip"`
	Country    string         `json:"country"`
	Phone      string         `json:"phone"`
	Email      string         `json:"email"`
}

type UpdateClientBody struct {
	Name       string         `json:"name"`
	FiscalCode string         `json:"fiscalCode"`
	Address    string         `json:"address"`
	Zip        string         `json:"zip"`
	Country    string         `json:"country"`
	Phone      string         `json:"phone"`
	Email      string         `json:"email"`
}

func (db *DB) MigrateClients() {
	db.instance.AutoMigrate(&Client{})
}

func (db *DB) CreateClient(body CreateClientBody) Client {
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

func (db *DB) UpdateClient(id int, body UpdateClientBody, client *Client) error {
	result := db.instance.Where("id = ?", id).First(&client, id)

	if result.Error != nil {
		return result.Error
	}

	result = db.instance.Model(&client).Updates(Client{Name: body.Name, Country: body.Country, Phone: body.Phone})

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
