package mutate

import (
	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/types"
)

type resource struct {
	base      v2.Resource
	access    v2.Access
	newAccess v2.Access
}

var _ v2.Resource = (*resource)(nil)

func (r *resource) compute() error {
	if r.newAccess != nil {
		r.access = r.newAccess
	}
	return nil
}

func (r *resource) Name() string {
	return r.base.Name()
}

func (r *resource) Access() (v2.Access, error) {
	if err := r.compute(); err != nil {
		return nil, err
	}
	return r.access, nil
}

func (r *resource) Digest() (*types.Digest, error) {
	// r.compute()
	return r.access.Digest()
}

func (r *resource) Type() types.ResourceType {
	return r.base.Type()
}

func (r *resource) Labels() map[string]string {
	return r.base.Labels()
}

func (r *resource) Deferrable() bool {
	return r.base.Deferrable()
}

func (r *resource) MarshalJSON() ([]byte, error) {
	return r.base.MarshalJSON()
}

func (r *resource) UnmarshalJSON(data []byte) error {
	return r.base.UnmarshalJSON(data)
}
