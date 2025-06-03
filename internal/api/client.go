package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/hazemKrimi/crimson-vault/internal/models"
	"github.com/labstack/echo/v4"
)

func (api *API) CreateClientHandler(context echo.Context) error {
	var body models.CreateClientRequestBody

	if err := context.Bind(&body); err != nil {
		log.Println(fmt.Sprintf("Error creating Client: %v.", err))
		return context.String(http.StatusBadRequest, "Invalid JSON!")
	}

	if err := context.Validate(body); err != nil {
		return err
	}

	client := api.db.CreateClient(body)

	log.Println(fmt.Sprintf("Client created with ID %d.", client.ID))
	return context.JSON(http.StatusOK, client)
}

func (api *API) GetAllClientsHandler(context echo.Context) error {
	clients, err := api.db.GetClients()

	if err != nil {
		return context.String(http.StatusInternalServerError, "Unexpected error getting Clients!")
	}

	log.Println("Got all Clients.")
	return context.JSON(http.StatusOK, clients)
}

func (api *API) GetClientHandler(context echo.Context) error {
	idString := context.Param("id")

	if idString == "" {
		return context.String(http.StatusBadRequest, "ID is required to get a Client!")
	}

	id, err := strconv.Atoi(idString)

	if err != nil {
		return context.String(http.StatusInternalServerError, "Unexpected error getting Client!")
	}

	var client models.Client

	if err := api.db.GetClient(id, &client); err != nil {
		return context.String(http.StatusNotFound, "Client not found!")
	}

	log.Println(fmt.Sprintf("Got client with ID %d.", client.ID))
	return context.JSON(http.StatusOK, client)
}

func (api *API) UpdateClientHandler(context echo.Context) error {
	idString := context.Param("id")

	if idString == "" {
		return context.String(http.StatusBadRequest, "ID is required to update a Client!")
	}

	id, err := strconv.Atoi(idString)

	if err != nil {
		return context.String(http.StatusInternalServerError, "Unexpected error updating Client!")
	}

	var body models.UpdateClientRequestBody

	if err := context.Bind(&body); err != nil {
		log.Println(fmt.Sprintf("Error updating Client: %v.", err))
		return context.String(http.StatusBadRequest, "Invalid JSON!")
	}

	if err := context.Validate(body); err != nil {
		return err
	}

	var client models.Client

	if err := api.db.UpdateClient(id, body, &client); err != nil {
		return context.String(http.StatusNotFound, "Client not found!")
	}

	log.Println(fmt.Sprintf("Updated client with ID %d.", client.ID))
	return context.JSON(http.StatusOK, client)
}

func (api *API) DeleteClientHandler(context echo.Context) error {
	idString := context.Param("id")

	if idString == "" {
		return context.String(http.StatusBadRequest, "ID is required to delete a Client!")
	}

	id, err := strconv.Atoi(idString)

	if err != nil {
		return context.String(http.StatusInternalServerError, "Unexpected error deleting Client!")
	}

	var client models.Client

	if err := api.db.DeleteClient(id); err != nil {
		return context.String(http.StatusNotFound, "Client not found!")
	}

	log.Println(fmt.Sprintf("Deleted client with ID %d.", client.ID))
	return context.String(http.StatusOK, "Client deleted successfully!")
}
