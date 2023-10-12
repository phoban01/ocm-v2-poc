package build

import (
	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/types"
)

type component struct {
	name      string
	version   string
	provider  string
	resources []v2.Resource
}

var _ v2.Component = (*component)(nil)

func New(name, version, provider string) v2.Component {
	return &component{
		name:     name,
		version:  version,
		provider: provider,
	}
}

func (c *component) compute() error {
	return nil
}

func (c *component) Version() string {
	return c.version
}

func (c *component) Provider() (*types.Provider, error) {
	return &types.Provider{
		Name: c.provider,
	}, nil
}

func (c *component) Descriptor() (*types.Descriptor, error) {
	return &types.Descriptor{
		ObjectMeta: types.ObjectMeta{
			Name:    c.name,
			Version: c.version,
		},
		Provider: types.Provider{
			Name: c.provider,
		},
	}, nil
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
	return nil, nil
}
