package mutate

import v2 "github.com/phoban01/ocm-v2/pkg/v2"

func WithRepositoryContext(base v2.Component, context ...v2.Repository) v2.Component {
	return &component{
		base:          base,
		addRepository: context,
	}
}

func AppendRepositoryContext(base v2.Component, context ...v2.Repository) v2.Component {
	return &component{
		base:          base,
		addRepository: context,
	}
}
