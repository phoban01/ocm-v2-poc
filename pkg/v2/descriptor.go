package v2

import "encoding/json"

// Descriptor represents the Component Descriptor in a structured way.
type Descriptor struct {
	Metadata   Metadata    `json:"metadata"`
	Provider   Provider    `json:"provider"`
	Version    string      `json:"version"`
	Resources  []Resource  `json:"resources"`
	Sources    []Source    `json:"sources"`
	References []Reference `json:"references"`
	Signatures []Signature `json:"signatures"`
}

type Metadata struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels,omitempty"`
}

type Provider struct {
	Name string `json:"name"`
}

func (d *Descriptor) MarshalJSON() ([]byte, error) {
	so, err := serialize(d)
	if err != nil {
		return nil, err
	}
	return json.Marshal(so)
}

type serializationFormat struct {
	Metadata   Metadata    `json:"metadata"`
	Provider   Provider    `json:"provider"`
	Version    string      `json:"version"`
	Resources  []resource  `json:"resources"`
	Signatures []signature `json:"signatures"`
}

type resource struct {
	Metadata     `json:",inline"`
	Digest       string `json:"digest"`
	MediaType    string `json:"media_type"`
	ResourceType string `json:"resource_type"`
}

type signature struct {
	Metadata `json:",inline"`
	Digest   string `json:"digest"`
}

func serialize(d *Descriptor) (*serializationFormat, error) {
	sf := &serializationFormat{
		Metadata: d.Metadata,
		Provider: d.Provider,
		Version:  d.Version,
	}

	for _, r := range d.Resources {
		dig, err := r.Digest()
		if err != nil {
			return nil, err
		}
		rx := resource{
			Metadata: Metadata{
				Name:   r.Name(),
				Labels: r.Labels(),
			},
			ResourceType: string(r.ResourceType()),
			MediaType:    string(r.MediaType()),
			Digest:       dig,
		}
		sf.Resources = append(sf.Resources, rx)
	}

	for _, r := range d.Signatures {
		dig, err := r.Digest()
		if err != nil {
			return nil, err
		}
		sx := signature{
			Metadata: Metadata{
				Name: r.Name(),
			},
			Digest: dig,
		}
		sf.Signatures = append(sf.Signatures, sx)
	}
	return sf, nil
}
