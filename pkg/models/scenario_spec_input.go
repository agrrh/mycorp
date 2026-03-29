package models

import "fmt"

type ScenarioSpecInputParameter struct {
	Name        string   `yaml:"name" json:"name"`
	Description string   `yaml:"description,omitempty" json:"description,omitempty"`
	Type        string   `yaml:"type,omitempty" json:"type,omitempty"`
	Aliases     []string `yaml:"aliases,omitempty" json:"aliases,omitempty"`
	Values      []string `yaml:"values,omitempty" json:"values,omitempty"`
	Default     any      `yaml:"default,omitempty" json:"default,omitempty"`
	Regexp      string   `yaml:"regexp,omitempty" json:"regexp,omitempty"`
	Optional    bool     `yaml:"optional,omitempty" json:"optional,omitempty"`
}

func (ip *ScenarioSpecInputParameter) GetCLIDescription() string {
	if ip.Description != "" {
		return ip.Description
	}

	return fmt.Sprintf("defaults to %v", ip.Default)
}
