package mutate

import v2 "github.com/phoban01/ocm-v2/pkg/v2"

func AddResources(base v2.Component, resources ...v2.Resource) v2.Component {
	return &component{
		base:         base,
		addResources: resources,
	}
}

func AddSignatures(base v2.Component, signatures ...v2.Signature) v2.Component {
	return &component{
		base:          base,
		addSignatures: signatures,
	}
}
