package mutate

import (
	"io"

	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/types"
)

type access struct {
	base       v2.Access
	accessType string
	mediaType  string
}

var _ v2.Access = (*access)(nil)

func (a *access) Type() string {
	if a.accessType != "" {
		return a.accessType
	}
	return a.base.Type()
}

func (a *access) MediaType() string {
	if a.mediaType != "" {
		return a.mediaType
	}
	return a.base.MediaType()
}

func (a *access) Labels() map[string]string {
	return a.base.Labels()
}

func (a *access) Data() (io.ReadCloser, error) {
	return a.base.Data()
}

func (a *access) Digest() (*types.Digest, error) {
	return a.base.Digest()
}

func (a *access) MarshalJSON() ([]byte, error) {
	return a.base.MarshalJSON()
}

func (a *access) UnmarshalJSON(data []byte) error {
	return a.base.UnmarshalJSON(data)
}
