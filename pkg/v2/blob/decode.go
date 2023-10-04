package blob

import (
	"encoding/json"
	"fmt"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/decode"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

func init() {
	decode.Register(Type, BlobDecoder{})
}

type BlobDecoder struct{}

func (d BlobDecoder) Decode(resource types.Resource) (v2.Resource, error) {
	a := access{}
	if err := json.Unmarshal(resource.Access, &a); err != nil {
		return nil, fmt.Errorf("failed to unmarshal access: %w", err)
	}
	return &blob{
		name: resource.Name,
		path: a.blob.path,
	}, nil
}
