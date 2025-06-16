package types

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         string         `json:"id" gorm:"type:varchar(255);primaryKey"`
	SessionID  string         `json:"-"`
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
	Username   string         `json:"username" gorm:"unique"`
	Password   string         `json:"-"`
	Clients    []Client       `json:"clients" gorm:"constraint:onDelete:CASCADE"`
}

type CreateUserRequestBody struct {
	Name       string `json:"name" validate:"required,alpha"`
	FiscalCode string `json:"fiscalCode" validate:"required"`
	Address    string `json:"address" validate:"required"`
	Zip        string `json:"zip" validate:"required"`
	Country    string `json:"country" validate:"required,alpha"`
	Phone      string `json:"phone" validate:"required,e164"`
	Email      string `json:"email" validate:"required,email"`
	Username   string `json:"username" validate:"required"`
}

type UpdateUserRequestBody struct {
	Name       string `json:"name" validate:"omitempty,alpha"`
	FiscalCode string `json:"fiscalCode"`
	Address    string `json:"address"`
	Zip        string `json:"zip"`
	Country    string `json:"country" validate:"omitempty,alpha"`
	Phone      string `json:"phone" validate:"omitempty,e164"`
	Email      string `json:"email" validate:"omitempty,email"`
	Username   string `json:"username"`
}

type UpdateUserSecurityCredentialsBody struct {
	Password        string `json:"password" validate:"password"`
	ConfirmPassword string `json:"confirmPassword" validate:"password,eqcsfield=Password"`
}
