package decode

import (
	"fmt"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/mutate"
	"github.com/phoban01/ocm-v2/pkg/v2/provider"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

func Resource(r types.Resource) (v2.Resource, error) {
	decoder, ok := lookup(r.Type)
	if !ok {
		return nil, fmt.Errorf("unknown resource type: %s", r)
	}

	resource, err := decoder.Decode(r)
	if err != nil {
		return nil, err
	}

	accessor, err := provider.Lookup(r)
	if err != nil {
		return nil, fmt.Errorf("unknown access type: %w", err)
	}

	return mutate.WithAccess(resource, accessor), nil
}
