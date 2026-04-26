package handlers

import (
	"github.com/agrrh/mycorp/internal/domain/scenario_store"
)

type Handler struct {
	ScStore *scenario_store.ScenarioStore
}
