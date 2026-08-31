package template

import (
	"testing"

	"github.com/agrrh/mycorp/internal/domain/modules"
	"github.com/agrrh/mycorp/internal/domain/scenario"
)

func TestRenderOutput_StaticOutput(t *testing.T) {
	tmplStr := "Static output"
	stepResults := make(modules.PrevStepsResults)
	inputs := scenario.SpecInputs{}

	result, err := RenderOutput(tmplStr, stepResults, inputs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != tmplStr {
		t.Errorf("expected %q, got %q", tmplStr, result)
	}
}

func TestRenderOutput_StepResults(t *testing.T) {
	tmplStr := `Result: {{ steps["step1"].stdout }}`
	stepResults := modules.PrevStepsResults{
		"step1": modules.StepResults{
			"stdout": "output from step 1",
		},
	}
	inputs := scenario.SpecInputs{}

	result, err := RenderOutput(tmplStr, stepResults, inputs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Result: output from step 1"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestRenderOutput_Conditional(t *testing.T) {
	tmplStr := `{% if steps["step1"].exit_code == 0 %}Success{% else %}Failed{% end %}`
	stepResults := modules.PrevStepsResults{
		"step1": modules.StepResults{
			"exit_code": 0,
		},
	}
	inputs := scenario.SpecInputs{}

	result, err := RenderOutput(tmplStr, stepResults, inputs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Success"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestRenderOutput_ConditionalElse(t *testing.T) {
	tmplStr := `{% if steps["step1"].exit_code == 0 %}Success{% else %}Failed{% end %}`
	stepResults := modules.PrevStepsResults{
		"step1": modules.StepResults{
			"exit_code": 1,
		},
	}
	inputs := scenario.SpecInputs{}

	result, err := RenderOutput(tmplStr, stepResults, inputs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Failed"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestProcessOutput_Empty(t *testing.T) {
	stepResults := make(modules.PrevStepsResults)
	inputs := scenario.SpecInputs{}

	result, err := ProcessOutput(scenario.SpecOutput(""), stepResults, inputs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestProcessOutput_WithSteps(t *testing.T) {
	tmplStr := scenario.SpecOutput(`Output from steps`)
	stepResults := modules.PrevStepsResults{
		"step1": modules.StepResults{
			"stdout": "hello",
		},
	}
	inputs := scenario.SpecInputs{}

	result, err := ProcessOutput(tmplStr, stepResults, inputs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != string(tmplStr) {
		t.Errorf("expected %q, got %q", tmplStr, result)
	}
}

func TestRenderOutput_MultiLine(t *testing.T) {
	tmplStr := `Created:
  - {{ steps["step1"].stdout }}
  - {{ steps["step2"].exit_code }}`
	stepResults := modules.PrevStepsResults{
		"step1": modules.StepResults{
			"stdout": "hello world",
		},
		"step2": modules.StepResults{
			"exit_code": 0,
		},
	}
	inputs := scenario.SpecInputs{}

	result, err := RenderOutput(tmplStr, stepResults, inputs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `Created:
  - hello world
  - 0`
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
