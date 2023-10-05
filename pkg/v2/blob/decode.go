package blob

import (
	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/decode"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

func init() {
	decode.Register(Type, BlobDecoder{})
}

type BlobDecoder struct{}

func (b BlobDecoder) Decode(resource types.Resource) (v2.Resource, error) {
	return &blob{
		name: resource.Name,
	}, nil
}
