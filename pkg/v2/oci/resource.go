package oci

import (
	"encoding/json"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

type image struct {
	name   string
	access *access
}

var _ v2.Resource = (*image)(nil)

const Type types.ResourceType = "ociImage"

type Option func(*image)

func WithMediaType(s string) Option {
	return func(b *image) {
		b.access.mediaType = s
	}
}

func Resource(name, ref string, opts ...Option) v2.Resource {
	i := &image{name: name, access: &access{
		ref: ref,
	}}
	for _, o := range opts {
		o(i)
	}
	return i
}

func (f *image) Name() string {
	return f.name
}

func (f *image) Access() v2.Access {
	return f.access
}

func (f *image) Digest() (*types.Digest, error) {
	return f.access.Digest()
}

func (f *image) Deferrable() bool {
	return true
}

func (f *image) Type() types.ResourceType {
	return Type
}

func (f *image) Labels() map[string]string {
	return map[string]string{
		// "ocm.software/reference": f.access.ref,
	}
}

func (f *image) WithLocation(url string) v2.Resource {
	f.access.WithLocation(url)
	return f
}

func (f *image) MarshalJSON() ([]byte, error) {
	return json.Marshal(f)
}

func (f *image) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, f)
}
