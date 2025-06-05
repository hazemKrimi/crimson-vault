package types

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
	Name       string `json:"name" validate:"required,alpha"`
	FiscalCode string `json:"fiscalCode"`
	Address    string `json:"address" validate:"required"`
	Zip        string `json:"zip" validate:"required"`
	Country    string `json:"country" validate:"required,alpha"`
	Phone      string `json:"phone" validate:"required,e164"`
	Email      string `json:"email" validate:"required,email"`
}

type UpdateClientRequestBody struct {
	Name       string `json:"name" validate:"omitempty,alpha"`
	FiscalCode string `json:"fiscalCode"`
	Address    string `json:"address"`
	Zip        string `json:"zip"`
	Country    string `json:"country" validate:"omitempty,alpha"`
	Phone      string `json:"phone" validate:"omitempty,e164"`
	Email      string `json:"email" validate:"omitempty,email"`
}
