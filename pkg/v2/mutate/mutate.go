package mutate

import v2 "github.com/phoban01/ocm-v2/pkg/v2"

func AddResources(base v2.Component, resources ...v2.Resource) v2.Component {
	return &component{
		base:         base,
		addResources: resources,
	}
}

func ReplaceResource(base v2.Component, resource v2.Resource) v2.Component {
	return &component{
		base:            base,
		replaceResource: resource,
	}
}
