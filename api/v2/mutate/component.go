package mutate

import (
	"encoding/json"
	"sync"

	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/types"
)

type component struct {
	base          v2.Component
	addResources  []v2.Resource
	addSignatures []v2.Signature
	addRepository []v2.Repository

	version        string
	computed       bool
	resources      []v2.Resource
	signatures     []v2.Signature
	storageContext []v2.RepositoryContext
	descriptor     *v2.Descriptor

	sync.Mutex
}

var _ v2.Component = (*component)(nil)

func (c *component) compute() error {
	c.Lock()
	defer c.Unlock()

	if c.computed {
		return nil
	}

	var updatedResources bool
	if c.addResources != nil {
		c.resources = c.addResources
		updatedResources = true
	}
	c.addResources = nil

	sigs, err := c.base.Signatures()
	if err != nil {
		return err
	}

	for _, add := range c.addSignatures {
		sigs = append(sigs, add)
	}

	c.signatures = sigs

	sctx, err := c.base.RepositoryContext()
	if err != nil {
		return err
	}

	for _, add := range c.addRepository {
		sctx = append(sctx, add.Context())
	}

	c.storageContext = sctx

	od, err := c.base.Descriptor()
	if err != nil {
		return err
	}

	c.descriptor = od

	if updatedResources {
		for _, r := range c.resources {
			acc, err := r.Access()
			if err != nil {
				return err
			}

			accData, err := json.Marshal(acc)
			if err != nil {
				return err
			}

			dig, err := r.Digest()
			if err != nil {
				return err
			}

			re := types.Resource{
				ObjectMeta: types.ObjectMeta{
					Name: r.Name(),
					Type: r.Type(),
				},
				Access: accData,
				Digest: dig,
			}
			c.descriptor.Resources = append(c.descriptor.Resources, re)
		}
	}

	c.descriptor.RepositoryContext = c.storageContext

	return nil
}

func (c *component) Version() string {
	if c.version != "" {
		return c.version
	}
	return c.base.Version()
}

func (c *component) Provider() (*v2.Provider, error) {
	return c.base.Provider()
}

func (c *component) RepositoryContext() ([]v2.RepositoryContext, error) {
	if err := c.compute(); err != nil {
		return nil, err
	}
	return c.storageContext, nil
}

func (c *component) Descriptor() (*v2.Descriptor, error) {
	if err := c.compute(); err != nil {
		return nil, err
	}
	return c.descriptor, nil
}

func (c *component) Resources() ([]v2.Resource, error) {
	if err := c.compute(); err != nil {
		return nil, err
	}
	return c.resources, nil
}

func (c *component) Sources() ([]v2.Source, error) {
	return c.base.Sources()
}

func (c *component) References() ([]v2.Reference, error) {
	return c.base.References()
}

func (c *component) Signatures() ([]v2.Signature, error) {
	if err := c.compute(); err != nil {
		return nil, err
	}
	return c.signatures, nil
}
