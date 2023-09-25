package mutate

import (
	"sync"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
)

type component struct {
	base          v2.Component
	addResources  []v2.Resource
	addSignatures []v2.Signature
	addStorage    []v2.Storage

	version        string
	computed       bool
	resources      []v2.Resource
	signatures     []v2.Signature
	storageContext []v2.StorageContext
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

	res, err := c.base.Resources()
	if err != nil {
		return err
	}

	for _, add := range c.addResources {
		res = append(res, add)
	}

	c.resources = res

	sigs, err := c.base.Signatures()
	if err != nil {
		return err
	}

	for _, add := range c.addSignatures {
		sigs = append(sigs, add)
	}

	c.signatures = sigs

	sctx, err := c.base.StorageContext()
	if err != nil {
		return err
	}

	for _, add := range c.addStorage {
		sctx = append(sctx, *add.Context())
	}

	c.storageContext = sctx

	od, err := c.base.Descriptor()
	if err != nil {
		return err
	}

	c.descriptor = od

	c.descriptor.Resources = c.resources

	c.descriptor.Signatures = c.signatures

	c.descriptor.StorageContext = c.storageContext

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

func (c *component) StorageContext() ([]v2.StorageContext, error) {
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
