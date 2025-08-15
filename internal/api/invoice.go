package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/hazemKrimi/crimson-vault/internal/lib"
	"github.com/hazemKrimi/crimson-vault/internal/types"
)

func (api *API) CreateItemHandler(context echo.Context) error {
	userId, err := uuid.Parse(context.Get("id").(string))

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Session ID not found after authorization."), Messages: []string{"Unexpected error getting User!"}}
	}

	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"Invoice ID is required to add an Invoice Item!"}}
	}

	var body types.CreateItemRequestBody

	if err := context.Bind(&body); err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"Invalid JSON!"}}
	}

	if err := context.Validate(body); err != nil {
		return err
	}

	item, err := api.db.CreateItem(userId, id, body)

	return context.JSON(http.StatusOK, item)
}

func (api *API) CreateInvoiceHandler(context echo.Context) error {
	userId, err := uuid.Parse(context.Get("id").(string))

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Session ID not found after authorization."), Messages: []string{"Unexpected error getting User!"}}
	}

	var body types.CreateInvoiceRequestBody

	if err := context.Bind(&body); err != nil {
		return types.Error{Code: http.StatusBadRequest, Cause: err, Messages: []string{"Invalid JSON!"}}
	}

	if err := context.Validate(body); err != nil {
		return err
	}

	clientId, err := uuid.Parse(body.ClientID)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Invalid Client ID."), Messages: []string{"Unexpected error getting User!"}}
	}

	var client types.Client

	if err := api.db.GetClientById(userId, clientId, &client); err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Unexpected error getting Client."), Messages: []string{"Unexpected error getting User!"}}
	}

	if client.UserID != context.Get("id") {
		return types.Error{Code: http.StatusBadRequest, Cause: err, Messages: []string{"Client does not belong to this User!"}}
	}

	var user types.User

	if err := api.db.GetUserById(userId, &user); err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Unexpected error getting User."), Messages: []string{"Unexpected error getting User!"}}
	}

	invoice, err := api.db.CreateInvoice(userId, body, user.IssuedInvoicesThisYear)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error creating Invoice!"}}
	}

	if err := api.db.UpdateUserIssuesInvoicesThisYear(userId, user.IssuedInvoicesThisYear + 1, &user); err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Unexpected error getting User."), Messages: []string{"Unexpected error getting User!"}}
	}

	path, err := lib.GenerateInvoice(invoice, user, client)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error generating Invoice!"}}
	}

	invoiceFileName := fmt.Sprintf("%s.pdf", invoice.ID)
	invoiceId, err := uuid.Parse(invoice.ID)

	if err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"Could not parse Invoice ID!"}}
	}

	err = api.db.UpdateInvoicePDF(userId, invoiceId, path, &invoice)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error generating Invoice!"}}
	}

	return context.Attachment(path, invoiceFileName)
}

func (api *API) GetAllItemsHandler(context echo.Context) error {
	userId, err := uuid.Parse(context.Get("id").(string))

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Session ID not found after authorization."), Messages: []string{"Unexpected error getting User!"}}
	}

	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"Invoice ID is required to get Items!"}}
	}

	items, err := api.db.GetItems(userId, id)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Messages: []string{"Unexpected error getting Items!"}}
	}

	return context.JSON(http.StatusOK, items)
}

func (api *API) GetAllInvoicesHandler(context echo.Context) error {
	userId, err := uuid.Parse(context.Get("id").(string))

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Session ID not found after authorization."), Messages: []string{"Unexpected error getting User!"}}
	}

	invoices, err := api.db.GetInvoices(userId)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error getting Invoices!"}}
	}

	return context.JSON(http.StatusOK, invoices)
}

func (api *API) GetItemHandler(context echo.Context) error {
	userId, err := uuid.Parse(context.Get("id").(string))

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Session ID not found after authorization."), Messages: []string{"Unexpected error getting User!"}}
	}

	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"ID is required to get an Invoice Item!"}}
	}

	var item types.Item

	if err := api.db.GetItemById(userId, id, &item); err != nil {
		return types.Error{Code: http.StatusNotFound, Messages: []string{"Invoice Item not found!"}}
	}

	return context.JSON(http.StatusOK, item)
}

func (api *API) GetInvoiceHandler(context echo.Context) error {
	userId, err := uuid.Parse(context.Get("id").(string))

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Session ID not found after authorization."), Messages: []string{"Unexpected error getting User!"}}
	}

	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"ID is required to get an Invoice!"}}
	}

	var invoice types.Invoice

	if err := api.db.GetInvoiceById(userId, id, &invoice); err != nil {
		return types.Error{Code: http.StatusNotFound, Messages: []string{"Invoice not found!"}}
	}

	return context.JSON(http.StatusOK, invoice)
}

func (api *API) DownloadInvoiceHandler(context echo.Context) error {
	userId, err := uuid.Parse(context.Get("id").(string))

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Session ID not found after authorization."), Messages: []string{"Unexpected error getting User!"}}
	}

	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"ID is required to download an Invoice!"}}
	}

	var invoice types.Invoice

	if err := api.db.GetInvoiceById(userId, id, &invoice); err != nil {
		return types.Error{Code: http.StatusNotFound, Messages: []string{"Invoice not found!"}}
	}

	if invoice.PDF != "" {
		invoiceFileName := fmt.Sprintf("%s.pdf", invoice.ID)
		return context.Attachment(invoice.PDF, invoiceFileName)
	}

	var user types.User

	if err := api.db.GetUserById(userId, &user); err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Unexpected error getting User."), Messages: []string{"Unexpected error getting User!"}}
	}

	clientId, err := uuid.Parse(invoice.ClientID)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Invalid Client ID."), Messages: []string{"Unexpected error getting User!"}}
	}

	var client types.Client

	if err := api.db.GetClientById(userId, clientId, &client); err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Unexpected error getting Client."), Messages: []string{"Unexpected error getting User!"}}
	}

	path, err := lib.GenerateInvoice(invoice, user, client)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error generating Invoice!"}}
	}

	invoiceFileName := fmt.Sprintf("%s.pdf", invoice.ID)
	invoiceId, err := uuid.Parse(invoice.ID)

	if err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"Could not parse Invoice ID!"}}
	}

	err = api.db.UpdateInvoicePDF(userId, invoiceId, path, &invoice)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error generating Invoice!"}}
	}

	return context.Attachment(path, invoiceFileName)
}

func (api *API) UpdateItemHandler(context echo.Context) error {
	userId, err := uuid.Parse(context.Get("id").(string))

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Session ID not found after authorization."), Messages: []string{"Unexpected error getting User!"}}
	}

	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"ID is required to update an Invoice Item!"}}
	}

	var body types.UpdateItemRequestBody

	if err := context.Bind(&body); err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"Invalid JSON!"}}
	}

	empty := body == types.UpdateItemRequestBody{}

	if empty {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"You must update at lease one field!"}}
	}

	var item types.Item

	if err := api.db.UpdateItem(userId, id, body, &item); err != nil {
		return types.Error{Code: http.StatusNotFound, Messages: []string{"Invoice Item not found!"}}
	}

	return context.JSON(http.StatusOK, item)
}

func (api *API) UpdateInvoiceHandler(context echo.Context) error {
	userId, err := uuid.Parse(context.Get("id").(string))

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Session ID not found after authorization."), Messages: []string{"Unexpected error getting User!"}}
	}

	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"ID is required to update an Invoice!"}}
	}

	var body types.UpdateInvoiceRequestBody

	if err := context.Bind(&body); err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"Invalid JSON!"}}
	}

	empty := body == types.UpdateInvoiceRequestBody{}

	if empty {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"You must update at lease one field!"}}
	}

	var item types.Item

	if err := api.db.UpdateInvoice(userId, id, body, &item); err != nil {
		return types.Error{Code: http.StatusNotFound, Messages: []string{"Invoice not found!"}}
	}

	return context.JSON(http.StatusOK, item)
}

func (api *API) DeleteItemHandler(context echo.Context) error {
	userId, err := uuid.Parse(context.Get("id").(string))

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Session ID not found after authorization."), Messages: []string{"Unexpected error getting User!"}}
	}

	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"ID is required to delete an Invoice Item!"}}
	}

	if err := api.db.DeleteItem(userId, id); err != nil {
		return types.Error{Code: http.StatusNotFound, Messages: []string{"Invoice Item not found!"}}
	}

	return context.JSON(http.StatusOK, map[string]string{"message": "Invoice Item deleted successfully!"})
}

func (api *API) DeleteInvoiceHandler(context echo.Context) error {
	userId, err := uuid.Parse(context.Get("id").(string))

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Session ID not found after authorization."), Messages: []string{"Unexpected error getting User!"}}
	}

	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"ID is required to delete an Invoice Item!"}}
	}

	if err := api.db.DeleteInvoice(userId, id); err != nil {
		return types.Error{Code: http.StatusNotFound, Messages: []string{"Invoice Item not found!"}}
	}

	return context.JSON(http.StatusOK, map[string]string{"message": "Invoice Item deleted successfully!"})
}
