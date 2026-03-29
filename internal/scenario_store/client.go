package scenario_store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/agrrh/mycorp/pkg/models"
)

var (
	errFetch                  error = errors.New("Failed to fetch scenarios")
	errReadResponse           error = errors.New("Failed to read scenarios response body")
	errParse                  error = errors.New("Failed to parse scenarios JSON")
	errObtainSpecificScenario error = errors.New("Failed to obtain specific scenario")
)

type ScenarioStoreCLI struct {
	URL       string
	Scenarios map[string]models.CLIScenario
}

func NewCLI(url string) *ScenarioStoreCLI {
	ssc := &ScenarioStoreCLI{
		URL:       url,
		Scenarios: make(map[string]models.CLIScenario),
	}

	return ssc
}

func (ssc *ScenarioStoreCLI) generateScenarioURLforCLI(scenarioFullname string) string {
	return fmt.Sprintf("%s/%s/_cli", ssc.URL, scenarioFullname)
}

func (ssc *ScenarioStoreCLI) Fetch() error {
	resp, err := http.Get(ssc.URL)
	if err != nil {
		return errors.Join(errFetch, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Join(errReadResponse, err)
	}

	var scenarios []string
	if err := json.Unmarshal(body, &scenarios); err != nil {
		return errors.Join(errParse, err)
	}

	for _, v := range scenarios {
		cliURL := ssc.generateScenarioURLforCLI(v)

		resp, err := http.Get(cliURL)
		if err != nil {
			return errors.Join(errObtainSpecificScenario, errFetch, err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Join(errObtainSpecificScenario, errReadResponse, err)
		}

		var sc models.CLIScenario
		if err := json.Unmarshal(body, &sc); err != nil {
			return errors.Join(errObtainSpecificScenario, errParse, err)
		}

		// Duplicate names are impossible, so not even checking

		ssc.Scenarios[sc.Metadata.GetFullName()] = sc
	}

	return nil
}

func (ssc *ScenarioStoreCLI) List() []*models.CLIScenario {
	var scList []*models.CLIScenario

	for _, sc := range ssc.Scenarios {
		scList = append(scList, &sc)
	}

	return scList
}
