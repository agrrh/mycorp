package scenario

import "fmt"

// TODO: Replace with kubernetes metadata model
type Metadata struct {
	Name      string `yaml:"name" json:"name"`
	Namespace string `yaml:"namespace" json:"namespace"`
}

func (m *Metadata) GetFullName() string {
	return fmt.Sprintf("%s/%s", m.Namespace, m.Name)
}

func (m *Metadata) GetName() string {
	return m.Name
}
