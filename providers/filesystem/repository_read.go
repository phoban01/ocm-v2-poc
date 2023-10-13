package filesystem

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/types"
)

func (r *repository) ReadBlob(path string) (io.ReadCloser, error) {
	fp := filepath.Join(r.blobdir, r.path, strings.TrimPrefix(path, "sha256:"))
	return os.Open(fp)
}

func (r *repository) Get(name, version string) (v2.Component, error) {
	loc := r.Context().Location()
	f, err := os.ReadFile(
		filepath.Join(strings.TrimPrefix(loc, "file://"), "component-descriptor.json"),
	)
	if err != nil {
		return nil, err
	}

	desc := types.Descriptor{}
	if err := json.Unmarshal(f, &desc); err != nil {
		return nil, err
	}

	return &component{
		repository: r,
		descriptor: desc,
	}, nil
}

func (r *repository) List() ([]v2.Component, error) {
	return nil, nil
}

func (r *repository) Delete() error {
	return nil
}
