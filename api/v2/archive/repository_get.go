package archive

import (
	"encoding/json"
	"os"
	"path/filepath"

	v2 "github.com/phoban01/ocm-v2/api/v2"
)

func (r *repository) Get(name, version string) (v2.Component, error) {
	f, err := os.ReadFile(filepath.Join(r.path, "component-descriptor.json"))
	if err != nil {
		return nil, err
	}

	desc := v2.Descriptor{}
	if err := json.Unmarshal(f, &desc); err != nil {
		return nil, err
	}

	return &component{
		descriptor: desc,
	}, nil
}

func (r *repository) List() ([]v2.Component, error) {
	return nil, nil
}

func (r *repository) Delete() error {
	return nil
}
