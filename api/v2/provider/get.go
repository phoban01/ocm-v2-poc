package provider

import (
	"fmt"

	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/builder"
	"github.com/phoban01/ocm-v2/api/v2/mutate"
	"github.com/phoban01/ocm-v2/api/v2/types"
)

func GetResource(ctx v2.RepositoryContext, resource types.Resource) (v2.Resource, error) {
	prov, err := lookup(resource)
	if err != nil {
		return nil, fmt.Errorf("unknown provider: %w", err)
	}

	accessor, err := prov.Decode(ctx, resource)
	if err != nil {
		return nil, fmt.Errorf("unknown access type: %w", err)
	}

	return mutate.WithAccess(builder.DecodeResource(resource), accessor), nil
}
