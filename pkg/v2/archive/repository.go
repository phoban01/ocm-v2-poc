package archive

import (
	"fmt"
	"path/filepath"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
)

type repository struct {
	path string
}

var _ v2.Repository = (*repository)(nil)

func Repository(path string) (v2.Repository, error) {
	p, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return &repository{path: p}, nil
}

func (r *repository) Context() *v2.RepositoryContext {
	return &v2.RepositoryContext{
		Type: "ocm.repository/bundle",
		URL:  fmt.Sprintf("file://%s", r.path),
	}
}
