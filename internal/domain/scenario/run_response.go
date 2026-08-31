package scenario

import "encoding/json"

type ScenarioRunResponse struct {
	Status  string          `json:"status"`
	Output  string          `json:"output"`
	Success bool            `json:"success"`
	Results json.RawMessage `json:"results,omitempty"`
}
