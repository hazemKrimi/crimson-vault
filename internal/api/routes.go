package api

import "github.com/labstack/echo/v4/middleware"

func (api *API) ClientRoutes() {
	clients := api.instance.Group("/clients")
	users := api.instance.Group("/users")
	
	clients.GET("/", api.GetAllClientsHandler)
	clients.POST("/", api.CreateClientHandler)
	clients.GET("/:id", api.GetClientHandler)
	clients.PUT("/:id", api.UpdateClientHandler)
	clients.DELETE("/:id", api.DeleteClientHandler)

	users.GET("/", api.GetAllUsersHandler)
	users.POST("/", api.CreateUserHandler)
	users.GET("/:id", api.GetUserHandler)
	users.PUT("/:id", api.UpdateUserHandler)
	users.PUT("/:id/security", api.UpdateUserSecurityDetailsHandler)
	users.PUT("/:id/logo", api.UpdateUserLogoHandler, middleware.BodyLimit("2M"))
	users.DELETE("/:id", api.DeleteUserHandler)
	users.DELETE("/:id/logo", api.DeleteUserLogoHandler)
}
