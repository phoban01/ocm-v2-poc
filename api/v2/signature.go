package v2

import "github.com/phoban01/ocm-v2/api/v2/types"

type Signature interface {
	Name() string

	Digest() (*types.Digest, error)

	Info() *types.SignatureInfo
}
