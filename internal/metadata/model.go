package metadata

import "fmt"

// TODO: Replace with kubernetes metadata model
type Metadata struct {
	Name      string `yaml:"name" json:"name"`
	Namespace string `yaml:"namespace" json:"namespace"`
}

func (sm *Metadata) GetFullName() string {
	return fmt.Sprintf("%s/%s", sm.Namespace, sm.Name)
}

func (sm *Metadata) GetName() string {
	return sm.Name
}
