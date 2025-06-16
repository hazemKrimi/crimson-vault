package models

import (
	"time"

	"github.com/google/uuid"

	"github.com/hazemKrimi/crimson-vault/internal/types"
)

func (db *DB) CreateItem(userId, invoiceId uuid.UUID, body types.CreateItemRequestBody) (types.Item, error) {
	item := types.Item{
		ID:        uuid.New().String(),
		InvoiceID: invoiceId.String(),
		UserID:    userId.String(),
		Name:      body.Name,
		Type:      body.Type,
		Quantity:  body.Quantity,
		Tax:       body.Tax,
	}

	result := db.instance.Create(&item)

	if result.Error != nil {
		return types.Item{}, result.Error
	}

	return item, nil
}

func (db *DB) CreateInvoice(userId uuid.UUID, body types.CreateInvoiceRequestBody) (types.Invoice, error) {
	dueAt, err := time.Parse("2006-01-02T15:04:05Z", body.DueAt)

	if err != nil {
		return types.Invoice{}, err
	}

	invoice := types.Invoice{
		ID:       uuid.New().String(),
		UserID:   userId.String(),
		ClientID: body.ClientID,
		DueAt:    dueAt,
		Currency: body.Currency,
		VAT:      body.VAT,
		Status:   types.Draft.String(),
	}

	result := db.instance.Create(&invoice)

	if result.Error != nil {
		return types.Invoice{}, result.Error
	}

	var items []types.Item

	for _, invoiceItem := range body.Items {
		invoiceId, err := uuid.Parse(invoice.ID)

		if err != nil {
			return types.Invoice{}, err
		}

		item, err := db.CreateItem(userId, invoiceId, invoiceItem)

		if err != nil {
			return types.Invoice{}, err
		}

		items = append(items, item)
	}

	result = db.instance.Model(&invoice).Updates(types.Invoice{
		Items: invoice.Items,
	})

	if result.Error != nil {
		return types.Invoice{}, result.Error
	}

	return invoice, nil
}

func (db *DB) GetItems(userId, id uuid.UUID) ([]types.Item, error) {
	var items []types.Item

	result := db.instance.Model(&types.Item{}).Where("user_id = ?", userId).Where("invoice_id = ?", id).Find(&items)

	if result.Error != nil {
		return nil, result.Error
	}

	return items, nil
}

func (db *DB) GetInvoices(userId uuid.UUID) ([]types.Invoice, error) {
	var invoices []types.Invoice

	result := db.instance.Model(&types.Invoice{}).Where("user_id = ?", userId).Preload("Items").Find(&invoices)

	if result.Error != nil {
		return nil, result.Error
	}

	return invoices, nil
}

func (db *DB) GetItemById(userId, id uuid.UUID, item *types.Item) error {
	result := db.instance.Model(&types.Item{}).Where("user_id = ?", userId).Where("id = ?", id).First(item)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) GetInvoiceById(userId, id uuid.UUID, invoice *types.Invoice) error {
	result := db.instance.Model(&types.Invoice{}).Preload("Items").Where("user_id = ?", userId).Where("id = ?", id).First(invoice)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) UpdateItem(userId, id uuid.UUID, body types.UpdateItemRequestBody, item *types.Item) error {
	result := db.instance.Where("user_id = ?", userId).Where("id = ?", id).First(item)

	if result.Error != nil {
		return result.Error
	}

	result = db.instance.Model(item).Updates(types.Item{
		Name:     body.Name,
		Type:     body.Type,
		Quantity: body.Quantity,
		Tax:      body.Tax,
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) UpdateInvoice(userId, id uuid.UUID, body types.UpdateInvoiceRequestBody, invoice *types.Item) error {
	result := db.instance.Where("user_id = ?", userId).Where("id = ?", id).First(invoice)

	if result.Error != nil {
		return result.Error
	}

	dueAt, err := time.Parse("2006-01-02T15:04:05.000Z", body.DueAt)

	if err != nil {
		return err
	}

	result = db.instance.Model(invoice).Updates(types.Invoice{
		DueAt:    dueAt,
		Currency: body.Currency,
		VAT:      body.VAT,
		Status:   body.Status,
	})

	return nil
}

func (db *DB) DeleteItem(userId, id uuid.UUID) error {
	result := db.instance.Unscoped().Where("user_id = ?", userId).Where("id = ?", id).Delete(&types.Item{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) DeleteInvoice(userId, id uuid.UUID) error {
	result := db.instance.Unscoped().Where("user_id = ?", userId).Where("id = ?", id).Delete(&types.Invoice{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}
