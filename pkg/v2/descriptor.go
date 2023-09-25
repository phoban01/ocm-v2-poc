package v2

import (
	"encoding/json"
	"fmt"

	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cuego"
)

// Descriptor represents the Component Descriptor in a structured way.
type Descriptor struct {
	Metadata       Metadata         `json:"metadata"`
	Version        string           `json:"version"`
	Provider       Provider         `json:"provider"`
	StorageContext []StorageContext `json:"storage_context"`
	Resources      []Resource       `json:"resources"`
	Sources        []Source         `json:"sources"`
	References     []Reference      `json:"references"`
	Signatures     []Signature      `json:"signatures"`
}

type Metadata struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels,omitempty"`
}

type Provider struct {
	Name string `json:"name"`
}

type serializationFormat struct {
	Metadata       Metadata         `json:"metadata"`
	Provider       Provider         `json:"provider"`
	Version        string           `json:"version"`
	StorageContext []StorageContext `json:"storage_context"`
	Resources      []resource       `json:"resources"`
	Signatures     []signature      `json:"signatures"`
}

func init() {
	c := `{
  metadata: name: =~ ".+\\..+"
  version: =~ "^v\\d+\\.\\d+\\.\\d+$"
}`

	cuego.MustConstrain(&serializationFormat{}, c)
}

type resource struct {
	Metadata     `       json:",inline"`
	Digest       string `json:"digest"`
	MediaType    string `json:"media_type"`
	ResourceType string `json:"resource_type"`
}

type signature struct {
	Metadata `       json:",inline"`
	Digest   string `json:"digest"`
}

func (d *Descriptor) MarshalJSON() ([]byte, error) {
	so, err := serialize(d)
	if err != nil {
		return nil, err
	}
	if err := cuego.Validate(so); err != nil {
		// TODO: print all validation errors
		errs := errors.Errors(err)
		return nil, fmt.Errorf("validation error: %w", errs[len(errs)-1])
	}
	return json.Marshal(so)
}

func serialize(d *Descriptor) (*serializationFormat, error) {
	sf := &serializationFormat{
		Metadata:       d.Metadata,
		StorageContext: d.StorageContext,
		Provider:       d.Provider,
		Version:        d.Version,
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
