package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) List(ctx echo.Context) error {
	var scList []string

	for _, sc := range h.ScStore.List() {
		scList = append(scList, sc.Metadata.GetFullName())
	}

	return ctx.JSON(http.StatusOK, scList)
}

func (h *Handler) ListByNamespace(ctx echo.Context) error {
	var scList []string

	for _, sc := range h.ScStore.ListByNamespace(ctx.Param("namespace")) {
		scList = append(scList, sc.Metadata.Name)
	}

	return ctx.JSON(http.StatusOK, scList)
}
