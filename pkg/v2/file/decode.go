package file

import (
	"encoding/json"
	"fmt"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/decode"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

func init() {
	decode.Register(Type, FileDecoder{})
}

type FileDecoder struct{}

func (d FileDecoder) Decode(resource types.Resource) (v2.Resource, error) {
	a := access{}
	if err := json.Unmarshal(resource.Access, &a); err != nil {
		return nil, fmt.Errorf("failed to unmarshal access: %w", err)
	}
	return &file{
		name: resource.Name,
		path: a.file.path,
	}, nil
}
