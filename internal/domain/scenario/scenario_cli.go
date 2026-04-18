package scenario

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	errCall         error = errors.New("error calling scenario")
	errReadResponse error = errors.New("error reading scenario call response")
	// errParseResponse error = errors.New("error parsing scenario call response")
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
	sc.Metadata = s.Metadata
	sc.Spec = CLISpec{
		Inputs: s.Spec.Inputs,
		Output: s.Spec.Output,
	}

	// TODO: Add validation

	return nil
}

func (sc *ScenarioCLI) Run(url string) (CLIOutputData, error) {
	output := CLIOutputData(string(""))

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(nil))
	if err != nil {
		return output, errors.Join(errCall, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return output, errors.Join(errReadResponse, err)
	}

	_ = resp.Body.Close()

	// if err := json.Unmarshal(body, &output); err != nil {
	// 	return errors.Join(errParseResponse, err)
	// }

	fmt.Println(resp.StatusCode)
	fmt.Printf("%s", string(body[:]))

	output = CLIOutputData(string(body[:]))

	return output, nil
}
