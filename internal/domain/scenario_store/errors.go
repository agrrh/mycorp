package scenario_store

import (
	"errors"
)

var (
	errDuplicateScenarios     error = errors.New("duplicate scenarios")
	errFetch                  error = errors.New("failed to fetch scenarios")
	errReadResponse           error = errors.New("failed to read scenarios response body")
	errParse                  error = errors.New("failed to parse scenarios JSON")
	errObtainSpecificScenario error = errors.New("failed to obtain specific scenario")
)
