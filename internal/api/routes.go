package api

import "net/http"

func ClientRoutes(api *APIWrapper) (*http.ServeMux) {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", api.GetClientsHandler)
	mux.HandleFunc("POST /", api.CreateClientHandler)

	return mux
}
