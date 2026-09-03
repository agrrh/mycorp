package scenario

import (
	"testing"
)

func TestScenario_Validate(t *testing.T) {
	tests := []struct {
		name    string
		sc      Scenario
		wantErr bool
	}{
		{
			name: "valid scenario",
			sc: Scenario{
				Kind:    "Scenario",
				Version: "1.0",
				Metadata: Metadata{
					Name:      "test-scenario",
					Namespace: "default",
				},
				Spec: ScenarioSpec{
					Inputs: []SpecInputParameter{
						{Name: "input1", Type: "string", Default: "value"},
						{Name: "input2", Type: "int", Default: 123},
					},
					Steps: []SpecStep{
						{Name: "step1", Module: "mod1"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid kind",
			sc: Scenario{
				Kind:    "Invalid",
				Version: "1.0",
				Metadata: Metadata{
					Name:      "test-scenario",
					Namespace: "default",
				},
			},
			wantErr: true,
		},
		{
			name: "missing version",
			sc: Scenario{
				Kind:    "Scenario",
				Version: "",
				Metadata: Metadata{
					Name:      "test-scenario",
					Namespace: "default",
				},
			},
			wantErr: true,
		},
		{
			name: "missing metadata name",
			sc: Scenario{
				Kind:    "Scenario",
				Version: "1.0",
				Metadata: Metadata{
					Name:      "",
					Namespace: "default",
				},
			},
			wantErr: true,
		},
		{
			name: "missing metadata namespace",
			sc: Scenario{
				Kind:    "Scenario",
				Version: "1.0",
				Metadata: Metadata{
					Name:      "test-scenario",
					Namespace: "",
				},
			},
			wantErr: true,
		},
		{
			name: "missing input name",
			sc: Scenario{
				Kind:    "Scenario",
				Version: "1.0",
				Metadata: Metadata{
					Name:      "test-scenario",
					Namespace: "default",
				},
				Spec: ScenarioSpec{
					Inputs: []SpecInputParameter{
						{Name: ""},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing step name",
			sc: Scenario{
				Kind:    "Scenario",
				Version: "1.0",
				Metadata: Metadata{
					Name:      "test-scenario",
					Namespace: "default",
				},
				Spec: ScenarioSpec{
					Steps: []SpecStep{
						{Name: "", Module: "mod1"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing step module",
			sc: Scenario{
				Kind:    "Scenario",
				Version: "1.0",
				Metadata: Metadata{
					Name:      "test-scenario",
					Namespace: "default",
				},
				Spec: ScenarioSpec{
					Steps: []SpecStep{
						{Name: "step1", Module: ""},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "input type mismatch",
			sc: Scenario{
				Kind:    "Scenario",
				Version: "1.0",
				Metadata: Metadata{
					Name:      "test-scenario",
					Namespace: "default",
				},
				Spec: ScenarioSpec{
					Inputs: []SpecInputParameter{
						{Name: "input1", Type: "string", Default: 123},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sc.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Scenario.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
