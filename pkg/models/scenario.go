package models

// Scenario represents a user‑defined scenario manifest.
type Scenario struct {
	Kind     string       `yaml:"kind" json:"kind"`
	Version  string       `yaml:"version" json:"version"`
	Metadata Metadata     `yaml:"metadata" json:"metadata"`
	Spec     ScenarioSpec `yaml:"spec" json:"spec"`
}

func (s *Scenario) ExportCLI() *CLIScenario {
	c := &CLIScenario{
		Metadata: s.Metadata,
		Spec: CLIScenarioSpec{
			Inputs: s.Spec.Inputs,
			Output: s.Spec.Output,
		},
	}

	return c
}
