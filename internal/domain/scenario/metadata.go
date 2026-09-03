package scenario

import "fmt"

// TODO: Replace with kubernetes metadata model
type Metadata struct {
	Name      string `yaml:"name" json:"name" validate:"required"`
	Namespace string `yaml:"namespace" json:"namespace" validate:"required"`
}

func (m *Metadata) GetFullName() string {
	return fmt.Sprintf("%s/%s", m.Namespace, m.Name)
}

func (m *Metadata) GetName() string {
	return m.Name
}
