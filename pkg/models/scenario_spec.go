package models

type ScenarioSpec struct {
	Access ScenarioSpecAccess           `yaml:"access,omitempty" json:"access,omitempty"`
	Inputs []ScenarioSpecInputParameter `yaml:"inputs" json:"inputs"`
	Steps  []ScenarioSpecStep           `yaml:"steps" json:"steps"`
	Output ScenarioSpecOutput           `yaml:"output,omitempty" json:"output,omitempty"`
}

type ScenarioSpecStep struct {
	Name   string `yaml:"name" json:"name"`
	Module string `yaml:"module" json:"module"`
}

type ScenarioSpecAccess struct {
	Allow []string `yaml:"allow,omitempty" json:"allow,omitempty"`
	Deny  []string `yaml:"deny,omitempty" json:"deny,omitempty"`
}

// Static output or templated string
type ScenarioSpecOutput string
