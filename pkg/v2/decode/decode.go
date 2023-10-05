package decode

import (
	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

var registry = make(map[types.ResourceType]Decoder)

type Decoder interface {
	Decode(types.Resource) (v2.Resource, error)
}

func Register(resourceType types.ResourceType, decoder Decoder) {
	registry[resourceType] = decoder
}

func lookup(resourceType types.ResourceType) (Decoder, bool) {
	decoder, found := registry[resourceType]
	return decoder, found
}
