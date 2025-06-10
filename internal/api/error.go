package api

import (
	"github.com/labstack/echo/v4"

	"github.com/hazemKrimi/crimson-vault/internal/types"
)

func (api *API) CustomErrorHandler(err error, context echo.Context) {
	if context.Response().Committed {
		return
	}

	custom, ok := err.(types.Error)

	if ok {
		cause, ok := custom.Cause.(error)

		if ok {
			api.Logger.Error(cause.Error())
		}

		context.JSON(custom.Code, map[string][]string{"errors": custom.Messages})
	}

	herr, ok := err.(*echo.HTTPError)

	if ok {
		message, ok := herr.Message.(string)

		if ok {
			context.JSON(herr.Code, map[string][]string{"errors": {message}})
		}
	}

	api.instance.DefaultHTTPErrorHandler(err, context)
}
