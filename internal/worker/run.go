package worker

import (
	"context"
	"fmt"

	"github.com/agrrh/mycorp/internal/modules"
	httpModule "github.com/agrrh/mycorp/internal/modules/http"
	"github.com/agrrh/mycorp/internal/scenario"
)

type Worker struct{}

func New() *Worker {
	return &Worker{}
}

func (w *Worker) RunScenario(sc *scenario.Scenario) (modules.PrevStepsResults, error) {
	prevStepResults := make(modules.PrevStepsResults)

	for _, step := range sc.Spec.Steps {
		fmt.Printf("Performing step: (%s) %s\n", step.Module, step.Name)

		// TODO: Demo purposes, use modules selector
		m := httpModule.New()

		stepResults, err := m.Run(context.TODO(), step.Name, sc.Spec.Inputs, step.Params, prevStepResults)

		prevStepResults[step.Name] = stepResults

		if err != nil {
			return prevStepResults, err
		}
	}

	return prevStepResults, nil
}
