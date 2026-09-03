package scenario

import "fmt"

type (
	SpecInputs     []SpecInputParameter
	SpecSteps      []SpecStep
	SpecStepParams map[string]any
)

type ScenarioSpec struct {
	Access SpecAccess `yaml:"access,omitempty" json:"access,omitempty"`
	Inputs SpecInputs `yaml:"inputs" json:"inputs" validate:"dive"`
	Steps  SpecSteps  `yaml:"steps" json:"steps" validate:"dive"`
	Output SpecOutput `yaml:"output,omitempty" json:"output,omitempty"`
}

type SpecStep struct {
	Name   string         `yaml:"name" json:"name" validate:"required"`
	Module string         `yaml:"module" json:"module" validate:"required"`
	Params SpecStepParams `yaml:"params" json:"params"`
}

type SpecAccess struct {
	Allow []string `yaml:"allow,omitempty" json:"allow,omitempty"`
	Deny  []string `yaml:"deny,omitempty" json:"deny,omitempty"`
}

// Static output or templated string
type SpecOutput string

type SpecInputParameter struct {
	Name        string   `yaml:"name" json:"name" validate:"required"`
	Description string   `yaml:"description,omitempty" json:"description,omitempty"`
	Type        string   `yaml:"type,omitempty" json:"type,omitempty"`
	Aliases     []string `yaml:"aliases,omitempty" json:"aliases,omitempty"`
	Values      []string `yaml:"values,omitempty" json:"values,omitempty"`
	Default     any      `yaml:"default,omitempty" json:"default,omitempty" validate:"inputParamCheck"`
	Regexp      string   `yaml:"regexp,omitempty" json:"regexp,omitempty"`
	Optional    bool     `yaml:"optional,omitempty" json:"optional,omitempty"`
}


func (ip *SpecInputParameter) GetCLIDescription() string {
	if ip.Description != "" {
		return ip.Description
	}

	return fmt.Sprintf("defaults to %v", ip.Default)
}
