package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetCLI(ctx echo.Context) error {
	scName := fmt.Sprintf("%s/%s", ctx.Param("namespace"), ctx.Param("name"))

	if cliSpec, exists := h.ScStore.Scenarios[scName]; exists {
		return ctx.JSON(http.StatusOK, cliSpec)
	}

	return ctx.JSON(http.StatusNotFound, nil)
}
