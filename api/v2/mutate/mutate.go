package mutate

import (
	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/types"
)

func WithResources(base v2.Component, resources ...v2.Resource) v2.Component {
	return &component{
		base:         base,
		addResources: resources,
	}
}

func WithAccess(base v2.Resource, access v2.Access) v2.Resource {
	return &resource{
		base:   base,
		access: access,
	}
}

func WithDigest(base v2.Resource, digest *types.Digest) v2.Resource {
	return &resource{
		base:   base,
		digest: digest,
	}
}

func WithAccessType(base v2.Access, newType string) v2.Access {
	return &access{
		base:       base,
		accessType: newType,
	}
}
