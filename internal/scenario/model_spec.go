package scenario

type (
	SpecInputs     []SpecInputParameter
	SpecSteps      []SpecStep
	SpecStepParams map[string]any
)

type Spec struct {
	Access SpecAccess `yaml:"access,omitempty" json:"access,omitempty"`
	Inputs SpecInputs `yaml:"inputs" json:"inputs"`
	Steps  SpecSteps  `yaml:"steps" json:"steps"`
	Output SpecOutput `yaml:"output,omitempty" json:"output,omitempty"`
}

type SpecStep struct {
	Name   string         `yaml:"name" json:"name"`
	Module string         `yaml:"module" json:"module"`
	Params SpecStepParams `yaml:"params" json:"params"`
}

type SpecAccess struct {
	Allow []string `yaml:"allow,omitempty" json:"allow,omitempty"`
	Deny  []string `yaml:"deny,omitempty" json:"deny,omitempty"`
}

// Static output or templated string
type SpecOutput string
