package api

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/hazemKrimi/crimson-vault/internal/types"
)

func (api *API) CreateClientHandler(context echo.Context) error {
	userId, ok := context.Get("id").(string)

	if !ok {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Session ID not found after authorization."), Messages: []string{"Unexpected error getting User!"}}
	}

	id, err := uuid.Parse(userId)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error getting User!"}}
	}

	var body types.CreateClientRequestBody

	if err := context.Bind(&body); err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"Invalid JSON!"}}
	}

	if err := context.Validate(body); err != nil {
		return err
	}

	client := api.db.CreateClient(id, body)

	return context.JSON(http.StatusOK, client)
}

func (api *API) GetAllClientsHandler(context echo.Context) error {
	userId, ok := context.Get("id").(string)

	if !ok {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Session ID not found after authorization."), Messages: []string{"Unexpected error getting User!"}}
	}

	id, err := uuid.Parse(userId)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error getting User!"}}
	}

	clients, err := api.db.GetClients(id)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error getting Clients!"}}
	}

	return context.JSON(http.StatusOK, clients)
}

func (api *API) GetClientHandler(context echo.Context) error {
	userIdString, ok := context.Get("id").(string)

	if !ok {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Session ID not found after authorization."), Messages: []string{"Unexpected error getting User!"}}
	}

	userId, err := uuid.Parse(userIdString)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error getting User!"}}
	}

	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"ID is required to get a Client!"}}
	}

	var client types.Client

	if err := api.db.GetClientById(userId, id, &client); err != nil {
		return types.Error{Code: http.StatusNotFound, Messages: []string{"Client not found!"}}
	}

	return context.JSON(http.StatusOK, client)
}

func (api *API) UpdateClientHandler(context echo.Context) error {
	userIdString, ok := context.Get("id").(string)

	if !ok {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Session ID not found after authorization."), Messages: []string{"Unexpected error getting User!"}}
	}

	userId, err := uuid.Parse(userIdString)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error getting User!"}}
	}

	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"ID is required to update a Client!"}}
	}

	var body types.UpdateClientRequestBody

	if err := context.Bind(&body); err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"Invalid JSON!"}}
	}

	empty := body == types.UpdateClientRequestBody{}

	if empty {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"You must update at least one field!"}}
	}

	if err := context.Validate(body); err != nil {
		return err
	}

	var client types.Client

	if err := api.db.UpdateClient(userId, id, body, &client); err != nil {
		return types.Error{Code: http.StatusNotFound, Messages: []string{"Client not found!"}}
	}

	return context.JSON(http.StatusOK, client)
}

func (api *API) DeleteClientHandler(context echo.Context) error {
	userIdString, ok := context.Get("id").(string)

	if !ok {
		return types.Error{Code: http.StatusInternalServerError, Cause: errors.New("Session ID not found after authorization."), Messages: []string{"Unexpected error getting User!"}}
	}

	userId, err := uuid.Parse(userIdString)

	if err != nil {
		return types.Error{Code: http.StatusInternalServerError, Cause: err, Messages: []string{"Unexpected error getting User!"}}
	}

	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		return types.Error{Code: http.StatusBadRequest, Messages: []string{"ID is required to delete a Client!"}}
	}

	if err := api.db.DeleteClient(userId, id); err != nil {
		return types.Error{Code: http.StatusNotFound, Messages: []string{"Client not found!"}}
	}

	return context.JSON(http.StatusOK, map[string]string{"message": "Client deleted successfully!"})
}
