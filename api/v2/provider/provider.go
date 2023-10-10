package provider

import (
	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/types"
)

type Provider interface {
	v2.Access
	Decode(ctx v2.RepositoryContext, resource types.Resource) (v2.Access, error)
}
