package v2

import (
	"io"

	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

type AccessType string

type Access interface {
	Labels() map[string]string

	Type() AccessType

	MediaType() string

	Digest() (*types.Digest, error)

	Data() (io.ReadCloser, error)

	MarshalJSON() ([]byte, error)

	UnmarshalJSON([]byte) error

	WithLocation(string)
}
