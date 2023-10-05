package provider

import (
	"fmt"

	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/builder"
	"github.com/phoban01/ocm-v2/api/v2/mutate"
	"github.com/phoban01/ocm-v2/api/v2/types"
)

func GetResource(r types.Resource) (v2.Resource, error) {
	resource := builder.DecodeResource(r)

	accessor, err := lookup(r)
	if err != nil {
		return nil, fmt.Errorf("unknown access type: %w", err)
	}

	return mutate.WithAccess(resource, accessor), nil
}
