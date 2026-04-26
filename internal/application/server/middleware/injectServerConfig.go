package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/agrrh/mycorp/internal/application/server/config"
)

// ServerHeader middleware adds a `Server` header to the response.
func InjectServerConfig(cfg config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Set("config", cfg)
			return next(ctx)
		}
	}
}
