package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/agrrh/mycorp/internal/domain/modules"
	"github.com/agrrh/mycorp/internal/domain/scenario"
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
	var (
		bytes []byte
		err   error
	)
	results := make(modules.StepResults)

	bytes, err = json.MarshalIndent(scInputs, "", "  ")
	if err != nil {
		return results, err
	}
	results["inputs"] = string(bytes)

	bytes, err = json.MarshalIndent(stepParams, "", "  ")
	if err != nil {
		return results, err
	}
	results["params"] = string(bytes)

	resp, err := http.Get(string(stepParams["url"].(string)))
	if err != nil {
		return results, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return results, err
	}

	_ = resp.Body.Close()

	results["code"] = resp.StatusCode
	results["body"] = string(body)

	return results, nil
}
