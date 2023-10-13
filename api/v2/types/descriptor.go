package types

import (
	"encoding/json"
)

// Descriptor represents the Component Descriptor in a structured way.
type Descriptor struct {
	ObjectMeta        `json:",inline"`
	Provider          Provider            `json:"provider"`
	RepositoryContext []RepositoryContext `json:",omitempty"`
	Resources         []Resource          `json:"resources,omitempty"`
	Sources           []Signature         `json:"sources,omitempty"`
	References        []Reference         `json:"references,omitempty"`
	Signatures        []Signature         `json:"signatures,omitempty"`
}

type Provider struct {
	Name string `json:"name"`
}

type RepositoryContext struct{}

func (d *Descriptor) MarshalJSON() ([]byte, error) {
	// do something here if you would like to
	// validate the descriptor
	// before it gets written
	type DescriptorAlias Descriptor
	out := &struct {
		*DescriptorAlias
	}{
		DescriptorAlias: (*DescriptorAlias)(d),
	}
	return json.Marshal(out)
}
