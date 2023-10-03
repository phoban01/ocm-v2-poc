package archive

import (
	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/decode"
)

type component struct {
	descriptor v2.Descriptor
	resources  []v2.Resource
	signatures []v2.Signature
}

var _ v2.Component = (*component)(nil)

func (c *component) compute() error {
	for _, res := range c.descriptor.Resources {
		dr, err := decode.Resource(res)
		if err != nil {
			return err
		}
		c.resources = append(c.resources, dr)
	}
	return nil
}

func (c *component) Version() string {
	return c.descriptor.Version
}

func (c *component) Provider() (*v2.Provider, error) {
	return &c.descriptor.Provider, nil
}

func (c *component) Descriptor() (*v2.Descriptor, error) {
	return &c.descriptor, nil
}

func (c *component) RepositoryContext() ([]v2.RepositoryContext, error) {
	return nil, nil
}

func (c *component) Resources() ([]v2.Resource, error) {
	if err := c.compute(); err != nil {
		return nil, err
	}
	return c.resources, nil
}

func (c *component) Sources() ([]v2.Source, error) {
	return nil, nil
}

func (c *component) References() ([]v2.Reference, error) {
	return nil, nil
}

func (c *component) Signatures() ([]v2.Signature, error) {
	return c.signatures, nil
}
