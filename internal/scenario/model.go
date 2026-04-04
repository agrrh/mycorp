package scenario

import (
	"github.com/agrrh/mycorp/internal/metadata"
)

// Scenario represents a user‑defined scenario manifest.
type Scenario struct {
	Kind     string            `yaml:"kind" json:"kind"`
	Version  string            `yaml:"version" json:"version"`
	Metadata metadata.Metadata `yaml:"metadata" json:"metadata"`
	Spec     Spec              `yaml:"spec" json:"spec"`
}

func (s *Scenario) ExportCLI() *ScenarioCLI {
	c := &ScenarioCLI{
		Metadata: s.Metadata,
		Spec: CLISpec{
			Inputs: s.Spec.Inputs,
			Output: s.Spec.Output,
		},
	}

	return c
}
