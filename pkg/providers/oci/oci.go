package oci

import (
	v1 "github.com/google/go-containerregistry/pkg/v1"
	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/provider"
)

type accessor struct {
	image v1.Image
	ref   string
}

var _ v2.Access = (*accessor)(nil)

var MediaType = "application/vnd.docker.image"

func init() {
	provider.Register(&accessor{})
}

func FromImage(ref string) (v2.Access, error) {
	return &accessor{ref: ref}, nil
}
