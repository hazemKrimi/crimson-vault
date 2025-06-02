package api

func (api *API) ClientRoutes() {
	group := api.instance.Group("/clients")
	
	group.GET("/", api.GetAllClientsHandler)
	group.POST("/", api.CreateClientHandler)
	group.GET("/:id", api.GetClientHandler)
	group.PUT("/:id", api.UpdateClientHandler)
	group.DELETE("/:id", api.DeleteClientHandler)
}
