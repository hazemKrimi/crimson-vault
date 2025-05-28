package models

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	Name string
	Country string
	Phone string
}

func (wrapper *DBWrapper) MigrateClients() {
	wrapper.db.AutoMigrate(&Client{})
}

func (wrapper *DBWrapper) CreateClient(name string, country string, phone string) {
	wrapper.db.Create(&Client{Name: name, Country: country, Phone: phone})
}

func (wrapper *DBWrapper) GetClient(id int) (Client) {
	var client Client
	wrapper.db.First(&client, id)
	return client
}
