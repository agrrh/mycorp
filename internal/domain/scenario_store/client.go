package scenario_store

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/agrrh/mycorp/internal/domain/scenario"
)

type ScenarioStoreCLI struct {
	URL       string
	Scenarios map[string]scenario.ScenarioCLI
}

func NewCLI(url string) *ScenarioStoreCLI {
	ssc := &ScenarioStoreCLI{
		URL:       url,
		Scenarios: make(map[string]scenario.ScenarioCLI),
	}

	return ssc
}

func (ssc *ScenarioStoreCLI) generateScenarioURLforCLI(scenarioFullname string) string {
	return fmt.Sprintf("%s/%s/_cli", ssc.URL, scenarioFullname)
}

func (ssc *ScenarioStoreCLI) Fetch() error {
	// TODO: Move reading auth var to common package
	req, err := prepareRequest("GET", fmt.Sprintf("%s/", ssc.URL))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Join(errFetch, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Join(errReadResponse, err)
	}

	_ = resp.Body.Close()

	var scenarios []string
	if err := json.Unmarshal(body, &scenarios); err != nil {
		return errors.Join(errParse, err, errors.New(string(body)))
	}

	for _, v := range scenarios {
		cliURL := ssc.generateScenarioURLforCLI(v)

		req, err := prepareRequest("GET", cliURL)
		if err != nil {
			return err
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return errors.Join(errFetch, err, errObtainSpecificScenario)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Join(errReadResponse, err, errObtainSpecificScenario)
		}

		_ = resp.Body.Close()

		var sc scenario.ScenarioCLI
		if err := json.Unmarshal(body, &sc); err != nil {
			return errors.Join(errParse, err, errObtainSpecificScenario)
		}

		// Duplicate names are impossible, so not even checking

		ssc.Scenarios[sc.Metadata.GetFullName()] = sc
	}

	return nil
}

func (ssc *ScenarioStoreCLI) List() []*scenario.ScenarioCLI {
	var scList []*scenario.ScenarioCLI

	for _, sc := range ssc.Scenarios {
		scList = append(scList, &sc)
	}

	return scList
}

func prepareRequest(method, url string) (*http.Request, error) {
	authToken := os.Getenv("MYCORP_TOKEN")

	req, err := http.NewRequest(method, url, bytes.NewBuffer(nil))
	if err != nil {
		return req, errors.Join(errFetch, err)
	}

	req.Header.Set("Content-Type", "application/json")

	if authToken != "" {
		req.Header.Set("X-Token", authToken)
	}

	return req, nil
}
