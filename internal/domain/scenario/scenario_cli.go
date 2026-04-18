package scenario

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

type ScenarioCLI struct {
	Metadata Metadata `yaml:"metadata" json:"metadata"`
	Spec     CLISpec  `yaml:"spec" json:"spec"`
}

type CLISpec struct {
	Inputs []SpecInputParameter `yaml:"inputs" json:"inputs"`
	Output SpecOutput           `yaml:"output,omitempty" json:"output"`
}

// Static output or variables for templated `SpecOutput` string
type CLIOutputData string

func (sc *ScenarioCLI) FromScenario(s *Scenario) error {
	sc = &ScenarioCLI{
		Metadata: s.Metadata,
		Spec: CLISpec{
			Inputs: s.Spec.Inputs,
			Output: s.Spec.Output,
		},
	}

	// TODO: Add validation

	return nil
}

func (sc *ScenarioCLI) Run(url string, output *CLIOutputData) error {
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

	fmt.Println(resp.StatusCode)
	fmt.Printf("%s", string(body[:]))

	tmp := CLIOutputData(string(body[:]))
	output = &tmp

	return nil
}
