package http

import (
	"context"
	"io"
	"net/http"

	"github.com/agrrh/mycorp/internal/modules"
	"github.com/agrrh/mycorp/internal/scenario"
)

type HttpModule struct {
	Name string
}

func New() *HttpModule {
	return &HttpModule{
		Name: "http",
	}
}

func (m *HttpModule) Run(ctx context.Context, stepName string, scInputs scenario.SpecInputs, stepParams scenario.SpecStepParams, prevStepResults modules.PrevStepsResults) (modules.StepResults, error) {
	results := make(modules.StepResults)

	resp, err := http.Get(string(stepParams["url"].(string)))
	if err != nil {
		return results, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	results["code"] = resp.StatusCode
	results["body"] = string(body)

	return results, nil
}
