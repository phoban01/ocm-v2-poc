package blob

import (
	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/provider"
)

type accessor struct {
	data     []byte
	filepath string
	labels   map[string]string
}

var _ v2.Access = (*accessor)(nil)

var MediaType = "custom-blob"

func init() {
	provider.Register(&accessor{})
}

func FromFile(path string) (v2.Access, error) {
	return &accessor{filepath: path}, nil
}

func FromBytes(data []byte) (v2.Access, error) {
	return &accessor{data: data}, nil
}
