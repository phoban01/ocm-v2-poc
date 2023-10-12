package provider

import (
	"fmt"

	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/build"
	"github.com/phoban01/ocm-v2/api/v2/mutate"
	"github.com/phoban01/ocm-v2/api/v2/types"
)

func GetAccessForResource(
	ctx v2.RepositoryContext,
	resource types.Resource,
	accessType string,
) (v2.Access, error) {
	prov, err := lookup(accessType, resource)
	if err != nil {
		return nil, fmt.Errorf("unknown provider: %w", err)
	}
	return prov.Decode(ctx, resource)
}

func GetResource(
	ctx v2.RepositoryContext,
	resource types.Resource,
	accessType string,
) (v2.Resource, error) {
	access, err := GetAccessForResource(ctx, resource, accessType)
	if err != nil {
		return nil, err
	}
	return mutate.WithAccess(build.DecodeResource(resource), access), nil
}
