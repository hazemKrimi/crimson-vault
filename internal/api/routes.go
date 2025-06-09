package api

import "github.com/labstack/echo/v4/middleware"

func (api *API) ClientRoutes() {
	clients := api.instance.Group("/api/clients", api.AuthSessionMiddleware)

	clients.GET("/", api.GetAllClientsHandler)
	clients.POST("/", api.CreateClientHandler)
	clients.GET("/:id/", api.GetClientHandler)
	clients.PUT("/:id/", api.UpdateClientHandler)
	clients.DELETE("/:id/", api.DeleteClientHandler)
}

func (api *API) UserRoutes() {
	users := api.instance.Group("/api/users")

	users.GET("/", api.GetAllUsersHandler)
	users.POST("/", api.CreateUserHandler)
	users.GET("/:id/", api.GetUserHandler)
	users.PUT("/:id/", api.UpdateUserHandler, api.AuthSessionMiddleware)
	users.PUT("/:id/security/", api.UpdateUserSecurityDetailsHandler)
	users.PUT("/:id/logo/", api.UpdateUserLogoHandler, middleware.BodyLimit("2M"))
	users.DELETE("/:id/", api.DeleteUserHandler, api.AuthSessionMiddleware)
	users.DELETE("/:id/logo/", api.DeleteUserLogoHandler, api.AuthSessionMiddleware)
}

func (api *API) AuthRoutes() {
	auth := api.instance.Group("/api/auth")

	auth.POST("/login/", api.LoginHandler)
	auth.DELETE("/logout/", api.LogoutHandler, api.AuthSessionMiddleware)
}
