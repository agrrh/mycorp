package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/agrrh/mycorp/internal/domain/scenario"
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
			resultsJSON, _ := json.Marshal(results)
			resp := scenario.ScenarioRunResponse{
				Status:  "error",
				Output:  fmt.Sprintf("%v+", results),
				Success: false,
				Results: resultsJSON,
			}
			return ctx.JSON(http.StatusInternalServerError, resp)
		}

		out, err := template.ProcessOutput(sc.Spec.Output, results, sc.Spec.Inputs)
		if err != nil {
			out = fmt.Sprintf("%v+", results)
		}

		resultsJSON, _ := json.Marshal(results)
		resp := scenario.ScenarioRunResponse{
			Status:  "success",
			Output:  out,
			Success: true,
			Results: resultsJSON,
		}
		return ctx.JSON(http.StatusOK, resp)
	}

	return ctx.JSON(http.StatusNotFound, nil)
}
