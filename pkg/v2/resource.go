package v2

import (
	"io"

	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

type Resource interface {
	Name() string

	Blob() (io.ReadCloser, error)

	Access() string

	Digest() (string, error)

	ResourceType() types.ResourceType

	MediaType() types.MediaType

	Labels() map[string]string

	Deferrable() bool
}
