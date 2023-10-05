package mutate

import (
	v2 "github.com/phoban01/ocm-v2/pkg/v2"
)

func WithAccess(base v2.Resource, access v2.Access) v2.Resource {
	return &resource{
		base:          base,
		updatedAccess: access,
	}
}

func WithResources(base v2.Component, resources ...v2.Resource) v2.Component {
	return &component{
		base:         base,
		addResources: resources,
	}
}
