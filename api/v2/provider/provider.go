package provider

import (
	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/types"
)

type Provider interface {
	v2.Access
	Decoder
}

type Decoder interface {
	Decode(v2.RepositoryContext, types.Resource) (v2.Access, error)
}
