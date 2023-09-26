package mutate

import (
	"io"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

type resource struct {
	base          v2.Resource
	updatedAccess string
	access        string
}

var _ v2.Resource = (*resource)(nil)

func SetAccess(base v2.Resource, access string) v2.Resource {
	return &resource{
		base:          base,
		updatedAccess: access,
	}
}

func (r *resource) compute() error {
	r.access = r.updatedAccess
	return nil
}

func (r *resource) Name() string {
	return r.base.Name()
}

func (r *resource) Access() string {
	r.compute()
	return r.access
}

func (r *resource) Blob() (io.ReadCloser, error) {
	return r.base.Blob()
}

func (r *resource) Digest() (string, error) {
	return r.base.Digest()
}

func (r *resource) ResourceType() types.ResourceType {
	return r.base.ResourceType()
}

func (r *resource) MediaType() types.MediaType {
	return r.base.MediaType()
}

func (r *resource) Labels() map[string]string {
	return r.base.Labels()
}

func (r *resource) Deferrable() bool {
	return r.base.Deferrable()
}
