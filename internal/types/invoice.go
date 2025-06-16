package types

import (
	"time"

	"gorm.io/gorm"
)

type ItemType int

const (
	Service ItemType = iota
	Product
)

func (itemType ItemType) String() string {
	switch itemType {
	case Service:
		return "service"
	case Product:
		return "product"
	default:
		return "unknown"
	}
}

type Item struct {
	ID        string `json:"id" gorm:"type:varchar(255);primaryKey"`
	InvoiceID string `json:"invoiceId" gorm:"type:varchar(255)"`
	UserID    string `json:"userId" gorm:"type:varchar(255)"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Price     uint32 `json:"price"`
	Quantity  uint32 `json:"quantity"`
	Tax       uint32 `json:"tax"`
}

type Status int

const (
	Draft Status = iota
	Posted
	Paid
	Late
)

func (status Status) String() string {
	switch status {
	case Draft:
		return "draft"
	case Posted:
		return "posted"
	case Paid:
		return "paid"
	case Late:
		return "late"
	default:
		return "unknown"
	}
}

type Invoice struct {
	ID        string         `json:"id" gorm:"type:varchar(255);primaryKey"`
	UserID    string         `json:"userId" gorm:"type:varchar(255)"`
	ClientID  string         `json:"clientId" gorm:"type:varchar(255)"`
	CreatedAt time.Time      `json:"createAt"`
	DueAt     time.Time      `json:"dueAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	Reference string         `json:"reference"`
	Status    string         `json:"status"`
	PDF       string         `json:"pdf"`
	Currency  string         `json:"currency"`
	VAT       uint32         `json:"vat"`
	Items     []Item         `json:"items"`
}

type CreateItemRequestBody struct {
	Name     string `json:"name" validate:"alpha,required"`
	Type     string `json:"type" validate:"alpha,required"`
	Price    uint32 `json:"price" validate:"number,required"`
	Quantity uint32 `json:"quantity" validate:"number,required"`
	Tax      uint32 `json:"tax" validate:"number,omitempty"`
}

type UpdateItemRequestBody struct {
	Name     string `json:"name" validate:"alpha,omitempty"`
	Type     string `json:"type" validate:"alpha,omitempty"`
	Price    uint32 `json:"price" validate:"number,omitempty"`
	Quantity uint32 `json:"quantity" validate:"number,omitempty"`
	Tax      uint32 `json:"tax" validate:"number,omitempty"`
}

type CreateInvoiceRequestBody struct {
	ClientID string                  `json:"clientId" validate:"uuid4,required"`
	DueAt    string                  `json:"dueAt" validate:"datetime=2006-01-02T15:04:05Z,required"`
	Currency string                  `json:"currency" validate:"iso4217,required"`
	VAT      uint32                  `json:"vat" validate:"number,required"`
	Items    []CreateItemRequestBody `json:"items"`
}

type UpdateInvoiceRequestBody struct {
	DueAt    string `json:"dueAt" validte:"datetime=2006-01-02T15:04:05Z,omitempty"`
	Currency string `json:"currency" validate:"iso4217,omitempty"`
	VAT      uint32 `json:"vat" validate:"number,omitempty"`
	Status   string `json:"status" validate:"alpha,omitempty"`
}
