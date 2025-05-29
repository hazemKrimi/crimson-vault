package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/hazemKrimi/crimson-vault/internal/models"
)

func (api *APIWrapper) CreateClientHandler(writer http.ResponseWriter, request *http.Request) {
	var body models.CreateClientBody

	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		http.Error(writer, "Invalid JSON", http.StatusBadRequest)
		log.Println(fmt.Sprintf("Error creating Client: %v.", err))
		return
	}

	client := api.dbWrapper.CreateClient(body)

	log.Println(fmt.Sprintf("Client created with ID %d.", client.ID))
	json.NewEncoder(writer).Encode(client)
}

func (api *APIWrapper) GetClientsHandler(writer http.ResponseWriter, request *http.Request) {
	idString := request.URL.Query().Get("id")

	if idString != "" {
		id, err := strconv.Atoi(idString)

		if err != nil {
			http.Error(writer, "Unexpected error getting Client.", http.StatusInternalServerError)
			return
		}

		var client models.Client

		if err := api.dbWrapper.GetClient(id, &client); err != nil {
			http.Error(writer, "Client not found.", http.StatusNotFound)
			return
		}

		log.Println(fmt.Sprintf("Got client with ID %d.", client.ID))
		json.NewEncoder(writer).Encode(client)
		return
	}

	clients, err := api.dbWrapper.GetClients()

	if err != nil {
		http.Error(writer, "Unexpected error getting Clients.", http.StatusInternalServerError)
		return
	}

	log.Println("Got all Clients.")
	json.NewEncoder(writer).Encode(clients)
}
