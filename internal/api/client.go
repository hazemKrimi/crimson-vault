package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/hazemKrimi/crimson-vault/internal/types"
)

func (api *API) CreateClientHandler(context echo.Context) error {
	userId, ok := context.Get("id").(string)

	if !ok {
		return context.String(http.StatusInternalServerError, "Unexpected error getting User!")
	}

	id, err := uuid.Parse(userId)

	if err != nil {
		return context.String(http.StatusInternalServerError, "Unexpected error getting User!")
	}

	var body types.CreateClientRequestBody

	if err := context.Bind(&body); err != nil {
		log.Println(fmt.Sprintf("Error creating Client: %v.", err))
		return context.String(http.StatusBadRequest, "Invalid JSON!")
	}

	if err := context.Validate(body); err != nil {
		return err
	}

	client := api.db.CreateClient(id, body)

	log.Println(fmt.Sprintf("Client created with ID %s.", client.ID))
	return context.JSON(http.StatusOK, client)
}

func (api *API) GetAllClientsHandler(context echo.Context) error {
	userId, ok := context.Get("id").(string)

	if !ok {
		return context.String(http.StatusInternalServerError, "Unexpected error getting User!")
	}

	id, err := uuid.Parse(userId)

	if err != nil {
		return context.String(http.StatusInternalServerError, "Unexpected error getting User!")
	}

	clients, err := api.db.GetClients(id)

	if err != nil {
		return context.String(http.StatusInternalServerError, "Unexpected error getting Clients!")
	}

	log.Println("Got all Clients.")
	return context.JSON(http.StatusOK, clients)
}

func (api *API) GetClientHandler(context echo.Context) error {
	userIdString, ok := context.Get("id").(string)

	if !ok {
		return context.String(http.StatusInternalServerError, "Unexpected error getting User!")
	}

	userId, err := uuid.Parse(userIdString)

	if err != nil {
		return context.String(http.StatusInternalServerError, "Unexpected error getting User!")
	}

	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		return context.String(http.StatusBadRequest, "ID is required to get a Client!")
	}

	var client types.Client

	if err := api.db.GetClientById(userId, id, &client); err != nil {
		return context.String(http.StatusNotFound, "Client not found!")
	}

	log.Println(fmt.Sprintf("Got User with ID %s.", client.ID))
	return context.JSON(http.StatusOK, client)
}

func (api *API) UpdateClientHandler(context echo.Context) error {
	userIdString, ok := context.Get("id").(string)

	if !ok {
		return context.String(http.StatusInternalServerError, "Unexpected error getting User!")
	}

	userId, err := uuid.Parse(userIdString)

	if err != nil {
		return context.String(http.StatusInternalServerError, "Unexpected error getting User!")
	}

	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		return context.String(http.StatusBadRequest, "ID is required to update a Client!")
	}

	var body types.UpdateClientRequestBody

	if err := context.Bind(&body); err != nil {
		log.Println(fmt.Sprintf("Error updating Client: %v.", err))
		return context.String(http.StatusBadRequest, "Invalid JSON!")
	}

	if err := context.Validate(body); err != nil {
		return err
	}

	var client types.Client

	if err := api.db.UpdateClient(userId, id, body, &client); err != nil {
		return context.String(http.StatusNotFound, "Client not found!")
	}

	log.Println(fmt.Sprintf("Updated Client with ID %s.", client.ID))
	return context.JSON(http.StatusOK, client)
}

func (api *API) DeleteClientHandler(context echo.Context) error {
	userIdString, ok := context.Get("id").(string)

	if !ok {
		return context.String(http.StatusInternalServerError, "Unexpected error getting User!")
	}

	userId, err := uuid.Parse(userIdString)

	if err != nil {
		return context.String(http.StatusInternalServerError, "Unexpected error getting User!")
	}

	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		return context.String(http.StatusBadRequest, "ID is required to delete a Client!")
	}

	if err := api.db.DeleteClient(userId, id); err != nil {
		return context.String(http.StatusNotFound, "Client not found!")
	}

	log.Println(fmt.Sprintf("Deleted Client with ID %s.", id))
	return context.String(http.StatusOK, "Client deleted successfully!")
}
