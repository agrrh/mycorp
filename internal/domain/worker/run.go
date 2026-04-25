package worker

import (
	"context"
	"errors"
	"fmt"

	"github.com/agrrh/mycorp/internal/domain/modules"
	commandModule "github.com/agrrh/mycorp/internal/domain/modules/command"
	httpModule "github.com/agrrh/mycorp/internal/domain/modules/http"
	"github.com/agrrh/mycorp/internal/domain/scenario"
)

var (
	errUnknownModule  error = errors.New("unknown module")
	errNotImplemented error = errors.New("not Implemented")
)

type Worker struct{}

func New() *Worker {
	return &Worker{}
}

// GetModule returns concrete module.
func (w *Worker) GetModule(moduleType string) (modules.Module, error) {
	// TODO: Dynamically select module by moduleType
	switch moduleType {
	case "http":
		return httpModule.New(), nil
	case "git":
		// TODO: git module implementation
		return nil, errNotImplemented
	case "yaml":
		// TODO: yaml module implementation
		return nil, errNotImplemented
	case "command":
		return commandModule.New(), nil
	default:
		// Для demo по умолчанию возвращаем http module
		return nil, errUnknownModule
	}
}

func (w *Worker) RunScenario(sc *scenario.Scenario) (modules.PrevStepsResults, error) {
	prevStepResults := make(modules.PrevStepsResults)

	for _, step := range sc.Spec.Steps {
		fmt.Printf("Performing step: (%s) %s\n", step.Module, step.Name)

		m, err := w.GetModule(step.Module)
		if err != nil {
			return nil, err
		}

		stepResults, err := m.Run(context.TODO(), step.Name, sc.Spec.Inputs, step.Params, prevStepResults)

		prevStepResults[step.Name] = stepResults

		if err != nil {
			return prevStepResults, err
		}
	}

	return prevStepResults, nil
}
