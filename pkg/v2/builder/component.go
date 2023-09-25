package builder

import v2 "github.com/phoban01/ocm-v2/pkg/v2"

type component struct {
	name     string
	version  string
	provider string
}

var _ v2.Component = (*component)(nil)

func New(name, version, provider string) v2.Component {
	return &component{
		name:     name,
		version:  version,
		provider: provider,
	}
}

func (c *component) Version() string {
	return c.version
}

func (c *component) Provider() (*v2.Provider, error) {
	return &v2.Provider{
		Name: c.provider,
	}, nil
}

func (c *component) Descriptor() (*v2.Descriptor, error) {
	return &v2.Descriptor{
		Metadata: v2.Metadata{
			Name: c.name,
		},
		Provider: v2.Provider{
			Name: c.provider,
		},
		Version: c.version,
	}, nil
}

func (c *component) Resources() ([]v2.Resource, error) {
	return nil, nil
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
