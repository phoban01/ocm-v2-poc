package decode

import (
	"fmt"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

func Resource(resource types.Resource) (v2.Resource, error) {
	decoder, ok := lookup(resource.Type)
	if !ok {
		return nil, fmt.Errorf("unknown resource type: %s", resource)
	}
	return decoder.Decode(resource)
}
