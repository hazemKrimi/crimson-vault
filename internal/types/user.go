package types

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint32         `json:"id"`
	CreatedAt  time.Time      `json:"createAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	Logo       string         `json:"logo"`
	Name       string         `json:"name"`
	FiscalCode string         `json:"fiscalCode"`
	Address    string         `json:"address"`
	Zip        string         `json:"zip"`
	Country    string         `json:"country"`
	Phone      string         `json:"phone"`
	Email      string         `json:"email"`
	Username   string         `json:"username"`
	Password   string         `json:"password"`
}

type CreateUserRequestBody struct {
	Name       string `json:"name" validate:"required"`
	FiscalCode string `json:"fiscalCode" validate:"required"`
	Address    string `json:"address" validate:"required"`
	Zip        string `json:"zip" validate:"required"`
	Country    string `json:"country" validate:"required"`
	Phone      string `json:"phone" validate:"required,e164"`
	Email      string `json:"email" validate:"required,email"`
}

type UpdateUserRequestBody struct {
	Name       string `json:"name"`
	FiscalCode string `json:"fiscalCode"`
	Address    string `json:"address"`
	Zip        string `json:"zip"`
	Country    string `json:"country"`
	Phone      string `json:"phone" validate:"omitempty,e164"`
	Email      string `json:"email" validate:"omitempty,email"`
}

type UpdateUserSecurityDetailsBody struct {
	Username        string `json:"username"`
	Password        string `json:"password" validate:"password"`
	ConfirmPassword string `json:"confirmPassword" validate:"password,eqcsfield=Password"`
}
