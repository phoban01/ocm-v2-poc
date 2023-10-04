package oci

import (
	"encoding/json"
	"fmt"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/decode"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

func init() {
	decode.Register(Type, ImageDecoder{})
}

type ImageDecoder struct{}

func (d ImageDecoder) Decode(resource types.Resource) (v2.Resource, error) {
	a := artifactAccess{}
	if err := json.Unmarshal(resource.Access, &a); err != nil {
		return nil, fmt.Errorf("failed to unmarshal access: %w", err)
	}
	return &image{
		name: resource.Name,
		ref:  a.image.ref,
	}, nil
}
