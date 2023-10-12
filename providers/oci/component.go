package oci

import (
	"encoding/json"
	"errors"
	"fmt"

	v2 "github.com/phoban01/ocm-v2/api/v2"

	"github.com/phoban01/ocm-v2/api/v2/build"
	"github.com/phoban01/ocm-v2/api/v2/mutate"
	"github.com/phoban01/ocm-v2/api/v2/provider"
	"github.com/phoban01/ocm-v2/api/v2/types"
)

type component struct {
	repository v2.Repository
	version    string
	descriptor types.Descriptor
	resources  []v2.Resource
	signatures []v2.Signature
}

var _ v2.Component = (*component)(nil)

func (c *component) compute() error {
	c.resources = make([]v2.Resource, len(c.descriptor.Resources))
	for i, item := range c.descriptor.Resources {
		res, err := c.processResource(item)
		if err != nil {
			return fmt.Errorf("failed to process resource: %w", err)
		}
		c.resources[i] = res
	}
	return nil
}

func (c *component) Version() string {
	return c.descriptor.ObjectMeta.Version
}

func (c *component) Provider() (*types.Provider, error) {
	return &c.descriptor.Provider, nil
}

func (c *component) Descriptor() (*types.Descriptor, error) {
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

func (c *component) processResource(item types.Resource) (v2.Resource, error) {
	accessMeta := make(map[string]interface{})
	if err := json.Unmarshal(item.Access, &accessMeta); err != nil {
		return nil, err
	}

	accessType, ok := accessMeta["type"].(string)
	if !ok {
		return nil, errors.New("access type not found or invalid")
	}

	switch accessType {
	case "localBlob/v1":
		mediaType, ok := accessMeta["mediaType"].(string)
		if !ok {
			return nil, errors.New("media type not found or invalid")
		}
		accessor := &localAccess{
			repository: c.repository,
			desc:       c.descriptor,
			mediaType:  mediaType,
			digest:     *item.Digest,
		}
		return mutate.WithAccess(build.DecodeResource(item), accessor), nil
	default:
		return provider.GetResource(c.repository.Context(), item, accessType)
	}
}
