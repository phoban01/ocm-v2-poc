package oci

import (
	"encoding/json"
	"fmt"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

func (f image) WithLocation(url string) v2.Resource {
	return &image{name: f.name, ref: url}
}

func (f *image) MarshalJSON() ([]byte, error) {
	return json.Marshal(f)
}

func (f *image) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, f)
}

func DecodeResource(resource types.Resource) (v2.Resource, error) {
	a := artifactAccess{}
	if err := json.Unmarshal(resource.Access, &a); err != nil {
		return nil, fmt.Errorf("failed to unmarshal access: %w", err)
	}
	return &image{
		name: resource.Name,
		ref:  a.image.ref,
	}, nil
}
