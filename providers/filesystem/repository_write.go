package filesystem

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/mutate"
)

func (r *repository) WriteBlob(acc v2.Access) (v2.Access, error) {
	if err := os.Mkdir(r.blobdir, os.ModePerm); err != nil {
		return nil, err
	}

	dig, err := acc.Digest()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(r.blobdir, dig.Value)

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader, err := acc.Data()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return nil, err
	}

	return ReadFile(path, WithMediaType(acc.MediaType()))
}

func (r *repository) Write(component v2.Component) error {
	if err := os.RemoveAll(r.path); err != nil {
		return err
	}

	if err := os.Mkdir(r.path, os.ModePerm); err != nil {
		return err
	}

	resources, err := component.Resources()
	if err != nil {
		return err
	}

	visitedResources := make([]v2.Resource, len(resources))
	for i, item := range resources {
		if item.Deferrable() {
			visitedResources[i] = item
			continue
		}

		acc, err := item.Access()
		if err != nil {
			return err
		}

		acc, err = r.WriteBlob(acc)
		if err != nil {
			return err
		}

		visitedResources[i] = mutate.WithAccess(item, acc)
	}

	component = mutate.WithResources(component, visitedResources...)

	desc, err := component.Descriptor()
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(r.path, "component-descriptor.json"))
	if err != nil {
		return err
	}

	data, err := json.Marshal(desc)
	if err != nil {
		return err
	}

	if _, err := f.Write(data); err != nil {
		return err
	}

	return nil
}
