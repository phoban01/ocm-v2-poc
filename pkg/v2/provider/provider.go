package provider

import (
	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

type Provider interface {
	v2.Access
	Decode(resource types.Resource) (v2.Access, error)
}
