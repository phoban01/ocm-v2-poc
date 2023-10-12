package v2

import (
	"github.com/phoban01/ocm-v2/api/v2/types"
)

type Resource interface {
	ObjectMeta
	Transferable
	Type() types.ResourceType
	Access() (Access, error)
	Digest() (*types.Digest, error)
}

type ObjectMeta interface {
	Name() string
	Labels() map[string]string
	Version() string
}

type Transferable interface {
	Deferrable() bool
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}
