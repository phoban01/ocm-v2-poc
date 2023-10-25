package v2

import (
	"io"

	"github.com/phoban01/ocm-v2/api/v2/types"
)

type Access interface {
	Labels() map[string]string

	Type() string

	MediaType() string

	Reference() string

	Digest() (*types.Digest, error)

	Data() (io.ReadCloser, error)

	Length() (int64, error)

	MarshalJSON() ([]byte, error)

	UnmarshalJSON([]byte) error
}
