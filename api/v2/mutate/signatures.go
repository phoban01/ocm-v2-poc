package mutate

import v2 "github.com/phoban01/ocm-v2/api/v2"

func AddSignatures(base v2.Component, signatures ...v2.Signature) v2.Component {
	return &component{
		base:          base,
		addSignatures: signatures,
	}
}
