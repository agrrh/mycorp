package modules

import (
	"context"

	"github.com/agrrh/mycorp/internal/scenario"
)

// TODO: Move those 2 to some ScenarioRunner package
type StepResults map[string]any
type PrevStepsResults map[string]StepResults

type Module interface {
	New() *Module
	Run(ctx context.Context, stepName string, scInputs scenario.SpecInputs, stepParams scenario.SpecStepParams, prevStepResults PrevStepsResults) (StepResults, error)
}
