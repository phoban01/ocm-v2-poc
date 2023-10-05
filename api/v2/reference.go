package v2

import (
	"crypto"
)

type Reference interface {
	Name() (string, error)

	Version() (string, error)

	Digest() (crypto.Hash, error)
}
