package store

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/agrrh/mycorp/internal/scenario"
)

type ScenarioStore struct {
	Dir       string
	Scenarios map[string]scenario.Scenario
}

func New(dir string) *ScenarioStore {
	scStore := &ScenarioStore{
		Dir:       dir,
		Scenarios: make(map[string]scenario.Scenario),
	}

	return scStore
}

func (ss *ScenarioStore) Load() error {
	entries, err := os.ReadDir(ss.Dir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") && !strings.HasSuffix(e.Name(), ".yml") {
			continue
		}

		data, err := os.ReadFile(filepath.Join(ss.Dir, e.Name()))
		if err != nil {
			log.Printf("failed to read %s: %v", e.Name(), err)
			continue
		}

		var sc scenario.Scenario
		if err := yaml.Unmarshal(data, &sc); err != nil {
			log.Printf("failed to parse %s: %v", e.Name(), err)
			continue
		}

		if sc.Kind != "Scenario" {
			log.Printf("%s: not a Scenario, skipping", e.Name())
			continue
		}

		key := sc.Metadata.GetFullName()

		if _, exists := ss.Scenarios[key]; exists {
			return errDuplicateScenarios
		}

		// TODO: Validate scenario itself (input variable types etc.)

		ss.Scenarios[key] = sc
		log.Printf("loaded scenario %s", key)
	}

	return nil
}

func (ss *ScenarioStore) List() []*scenario.Scenario {
	var scList []*scenario.Scenario

	for _, sc := range ss.Scenarios {
		scList = append(scList, &sc)
	}

	return scList
}

func (ss *ScenarioStore) ListByNamespace(namespace string) []*scenario.Scenario {
	var scList []*scenario.Scenario

	for _, sc := range ss.Scenarios {
		if sc.Metadata.Namespace == namespace {
			scList = append(scList, &sc)
		}
	}

	return scList
}
