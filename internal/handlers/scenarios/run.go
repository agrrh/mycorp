package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Run(ctx echo.Context) error {
	scName := fmt.Sprintf(
		"%s/%s",
		ctx.Param("namespace"),
		ctx.Param("name"),
	)

	if sc, exists := h.ScStore.Scenarios[scName]; exists {
		return ctx.String(http.StatusOK, string(sc.Spec.Output))
	}

	return ctx.JSON(http.StatusNotFound, nil)
}
