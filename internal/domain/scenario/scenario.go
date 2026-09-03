package scenario

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	validate.RegisterValidation("inputParamCheck", func(fl validator.FieldLevel) bool {
		// fl.Parent() returns the struct containing the field
		param, ok := fl.Parent().Interface().(SpecInputParameter)
		if !ok {
			return false
		}

		if param.Type != "" && param.Default != nil {
			expectedType := param.Type
			actualValue := param.Default

			var isValid bool
			switch expectedType {
			case "string":
				_, isValid = actualValue.(string)
			case "int":
				_, isValid = actualValue.(int)
			case "bool":
				_, isValid = actualValue.(bool)
			case "float":
				_, isValid = actualValue.(float64)
			default:
				// If type is unknown, we might want to allow it or fail.
				// For now, let's just allow unknown types.
				isValid = true
			}

			if !isValid {
				// We use the field's name/path to report error
				return false
			}
		}
		return true
	})
}

// Scenario represents a user‑defined scenario manifest.
type Scenario struct {
	Kind     string       `yaml:"kind" json:"kind" validate:"required,eq=Scenario"`
	Version  string       `yaml:"version" json:"version" validate:"required"`
	Metadata Metadata     `yaml:"metadata" json:"metadata" validate:"required"`
	Spec     ScenarioSpec `yaml:"spec" json:"spec" validate:"required"`
}

// Validate checks if the scenario is valid.
func (s *Scenario) Validate() error {
	if err := validate.Struct(s); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}
