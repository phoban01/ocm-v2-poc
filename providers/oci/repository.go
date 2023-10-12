package oci

import (
	"fmt"

	v2 "github.com/phoban01/ocm-v2/api/v2"
)

type repository struct {
	prefix   string
	registry string
}

var _ v2.Repository = (*repository)(nil)

func Repository(registry string) (v2.Repository, error) {
	return &repository{
		prefix:   "component-descriptors",
		registry: registry,
	}, nil
}

func (s *repository) Context() v2.RepositoryContext {
	return &repositoryContext{
		url: fmt.Sprintf("oci://%s", s.registry),
	}
}
