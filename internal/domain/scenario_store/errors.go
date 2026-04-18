package scenario_store

import (
	"errors"
)

var (
	errDuplicateScenarios     error = errors.New("Duplicate scenarios")
	errFetch                  error = errors.New("Failed to fetch scenarios")
	errReadResponse           error = errors.New("Failed to read scenarios response body")
	errParse                  error = errors.New("Failed to parse scenarios JSON")
	errObtainSpecificScenario error = errors.New("Failed to obtain specific scenario")
)
