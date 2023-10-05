package helm

import (
	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/provider"
)

type accessor struct {
	data     []byte
	filepath string
	labels   map[string]string
}

var _ v2.Access = (*accessor)(nil)

var MediaType = "helm"

func init() {
	provider.Register(&accessor{})
}

func FromRepository(path string) (v2.Access, error) {
	return &accessor{filepath: path}, nil
}

func FromBytes(data []byte) (v2.Access, error) {
	return &accessor{data: data}, nil
}
