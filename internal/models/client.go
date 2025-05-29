package models

import (
	"time"

	"gorm.io/gorm"
)

type Client struct {
	ID        uint32         `json:"id"`
	CreatedAt time.Time      `json:"createAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	Name      string         `json:"name"`
	Country   string         `json:"country"`
	Phone     string         `json:"phone"`
}

type CreateClientBody struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	Phone   string `json:"phone"`
}

func (wrapper *DBWrapper) MigrateClients() {
	wrapper.db.AutoMigrate(&Client{})
}

func (wrapper *DBWrapper) CreateClient(body CreateClientBody) Client {
	client := Client{Name: body.Name, Country: body.Country, Phone: body.Phone}

	wrapper.db.Create(&client)
	return client
}

func (wrapper *DBWrapper) GetClients() ([]Client, error) {
	var clients []Client

	result := wrapper.db.Find(&clients)

	if result.Error != nil {
		return nil, result.Error
	}

	return clients, nil
}

func (wrapper *DBWrapper) GetClient(id int, client *Client) error {
	result := wrapper.db.Where("id = ?", id).First(&client, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
