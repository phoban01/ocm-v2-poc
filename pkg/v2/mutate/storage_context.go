package mutate

import v2 "github.com/phoban01/ocm-v2/pkg/v2"

func AddStorageContext(base v2.Component, context ...v2.Storage) v2.Component {
	return &component{
		base:       base,
		addStorage: context,
	}
}
