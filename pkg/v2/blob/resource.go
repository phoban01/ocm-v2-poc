package blob

import (
	"encoding/json"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

type blob struct {
	name   string
	access *access
}

const Type types.ResourceType = "blob"

var _ v2.Resource = (*blob)(nil)

type Option func(*blob)

func WithMediaType(s string) Option {
	return func(b *blob) {
		b.access.mediaType = s
	}
}

func FromBytes(name string, data []byte, opts ...Option) v2.Resource {
	b := &blob{name: name, access: &access{
		data:      data,
		mediaType: "blob",
	}}
	for _, o := range opts {
		o(b)
	}
	return b
}

func FromFile(name string, path string, opts ...Option) v2.Resource {
	b := &blob{name: name, access: &access{
		path:      path,
		mediaType: "blob",
	}}
	for _, o := range opts {
		o(b)
	}
	return b
}

func (f *blob) Name() string {
	return f.name
}

func (f *blob) Type() types.ResourceType {
	return Type
}

func (f *blob) Labels() map[string]string {
	return map[string]string{
		"ocm.software/blobname": f.access.path,
	}
}

func (f *blob) Deferrable() bool {
	return false
}

func (f *blob) Access() v2.Access {
	return f.access
}

func (f *blob) Digest() (*types.Digest, error) {
	return f.access.Digest()
}

func (f *blob) MarshalJSON() ([]byte, error) {
	access, err := json.Marshal(f.Access())
	if err != nil {
		return nil, err
	}
	dig, err := f.Digest()
	if err != nil {
		return nil, err
	}
	r := types.Resource{
		Name:   f.name,
		Type:   f.Type(),
		Access: access,
		Digest: dig,
	}
	return json.Marshal(r)
}

func (f *blob) UnmarshalJSON(data []byte) error {
	r := types.Resource{}
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	f.name = r.Name
	return nil
}

func (f *blob) WithLocation(p string) v2.Resource {
	f.access.WithLocation(p)
	return f
}
