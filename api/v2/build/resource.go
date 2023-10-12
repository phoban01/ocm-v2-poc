package build

import (
	"encoding/json"

	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/types"
)

type resource struct {
	name         string
	access       v2.Access
	deferrable   bool
	resourceType types.ResourceType
	labels       map[string]string
	digest       *types.Digest
}

var _ v2.Resource = (*resource)(nil)

type Option func(*resource)

func Deferrable(value bool) Option {
	return func(r *resource) {
		r.deferrable = value
	}
}

func NewResource(meta types.ObjectMeta, access v2.Access, opts ...Option) v2.Resource {
	r := &resource{name: meta.Name, resourceType: meta.Type, access: access}
	for _, o := range opts {
		o(r)
	}
	return r
}

func DecodeResource(r types.Resource) v2.Resource {
	return &resource{
		name:         r.Name,
		resourceType: r.Type,
	}
}

func (r *resource) Name() string {
	return r.name
}

func (r *resource) Type() types.ResourceType {
	return r.resourceType
}

func (r *resource) Labels() map[string]string {
	return r.labels
}

func (r *resource) Deferrable() bool {
	return r.deferrable
}

func (r *resource) Access() (v2.Access, error) {
	return r.access, nil
}

func (r *resource) Digest() (*types.Digest, error) {
	if r.digest != nil {
		return r.digest, nil
	}
	return r.access.Digest()
}

func (r *resource) MarshalJSON() ([]byte, error) {
	data, err := r.Access()
	if err != nil {
		return nil, err
	}
	access, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	dig, err := r.Digest()
	if err != nil {
		return nil, err
	}
	return json.Marshal(types.Resource{
		ObjectMeta: types.ObjectMeta{
			Name: r.name,
			Type: r.Type(),
		},
		Access: access,
		Digest: dig,
	})
}

func (r *resource) UnmarshalJSON(data []byte) error {
	res := types.Resource{}
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}
	r.name = res.Name
	r.resourceType = res.Type
	r.digest = res.Digest
	return nil
}
