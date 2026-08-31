package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/agrrh/mycorp/internal/domain/template"
	"github.com/agrrh/mycorp/internal/domain/worker"
)

func (h *Handler) Run(ctx echo.Context) error {
	scName := fmt.Sprintf(
		"%s/%s",
		ctx.Param("namespace"),
		ctx.Param("name"),
	)

	w := worker.New()

	if sc, exists := h.ScStore.Scenarios[scName]; exists {
		results, err := w.RunScenario(&sc)

		if err != nil {
			out := fmt.Sprintf("%v+", results)
			return ctx.String(http.StatusInternalServerError, out)
		}

		out, err := template.ProcessOutput(sc.Spec.Output, results, sc.Spec.Inputs)
		if err != nil {
			out = fmt.Sprintf("%v+", results)
		}

		return ctx.String(http.StatusOK, out)
	}

	return ctx.JSON(http.StatusNotFound, nil)
}
