package v2

import (
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

type Resource interface {
	ObjectMeta
	Transporter
	Type() types.ResourceType
	Access() Access
	Digest() (*types.Digest, error)
}

type ObjectMeta interface {
	Name() string
	Labels() map[string]string
}

type Transporter interface {
	WithLocation(string) Resource
	Deferrable() bool
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}
