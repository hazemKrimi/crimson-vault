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
	users.GET("/me/", api.GetUserHandler, api.AuthSessionMiddleware)
	users.PUT("/me/", api.UpdateUserHandler, api.AuthSessionMiddleware)
	users.PUT("/me/security/", api.UpdateUserSecurityCredentialsHandler, api.AuthSessionMiddleware)
	users.PUT("/me/logo/", api.UpdateUserLogoHandler, middleware.BodyLimit("2M"), api.AuthSessionMiddleware)
	users.DELETE("/me/", api.DeleteUserHandler, api.AuthSessionMiddleware)
	users.DELETE("/me/logo/", api.DeleteUserLogoHandler, api.AuthSessionMiddleware)
}

func (api *API) InvoiceRoutes() {
	invoices := api.instance.Group("/api/invoices", api.AuthSessionMiddleware)

	invoices.GET("/", api.GetAllInvoicesHandler)
	invoices.POST("/", api.CreateInvoiceHandler)
	invoices.POST("/:id/items/", api.CreateItemHandler)
	invoices.GET("/:id/", api.GetInvoiceHandler)
	invoices.GET("/:id/items/", api.GetAllItemsHandler)
	invoices.GET("/items/:id/", api.GetItemHandler)
	invoices.PUT("/:id/", api.UpdateInvoiceHandler)
	invoices.PUT("/items/:id/", api.UpdateItemHandler)
	invoices.DELETE("/:id/", api.DeleteInvoiceHandler)
	invoices.DELETE("/items/:id/", api.DeleteItemHandler)
}

func (api *API) AuthRoutes() {
	auth := api.instance.Group("/api/auth")

	auth.POST("/login/", api.LoginHandler)
	auth.DELETE("/logout/", api.LogoutHandler, api.AuthSessionMiddleware)
}
