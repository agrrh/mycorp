package models

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	errCall         error = errors.New("Error calling scenario")
	errReadResponse error = errors.New("Error reading scenario call response")
	// errParseResponse error = errors.New("Error parsing scenario call response")
)

type CLIScenario struct {
	Metadata Metadata        `yaml:"metadata" json:"metadata"`
	Spec     CLIScenarioSpec `yaml:"spec" json:"spec"`
}

type CLIScenarioSpec struct {
	Inputs []ScenarioSpecInputParameter `yaml:"inputs" json:"inputs"`
	Output ScenarioSpecOutput           `yaml:"output,omitempty" json:"output"`
}

// Static output or variables for templated `ScenarioSpecOutput` string
type ScenarioCLIOutputData string

func (cs *CLIScenario) Run(url string, output *ScenarioCLIOutputData) error {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(nil))
	if err != nil {
		return errors.Join(errCall, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Join(errReadResponse, err)
	}

	// if err := json.Unmarshal(body, &output); err != nil {
	// 	return errors.Join(errParseResponse, err)
	// }

	fmt.Printf("%s", string(body[:]))

	tmp := ScenarioCLIOutputData(string(body[:]))
	output = &tmp

	return nil
}
