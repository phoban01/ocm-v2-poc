package oci

import (
	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/decode"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

func init() {
	decode.Register(Type, ImageDecoder{})
}

type ImageDecoder struct{}

func (d ImageDecoder) Decode(resource types.Resource) (v2.Resource, error) {
	return &image{
		name: resource.Name,
	}, nil
}
