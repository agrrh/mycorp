package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"

	"github.com/agrrh/mycorp/internal/application/server/config"
)

var (
	errFailedToLoadServerConfig error = errors.New("could not load server config from context")
)

func AuthTokens(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		serverConfig, ok := ctx.Get("config").(config.Config)
		if !ok {
			fmt.Println(ctx.Get("config"))
			return errFailedToLoadServerConfig
		}

		reqToken := ctx.Request().Header.Get("X-Token")

		if valid := slices.Contains(serverConfig.Tokens, reqToken); valid {
			return next(ctx)
		}

		return ctx.String(http.StatusUnauthorized, "Pass valid auth token via X-Token header")
	}
}
