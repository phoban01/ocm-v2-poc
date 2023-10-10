package filesystem

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/mutate"
)

func (r *repository) Write(component v2.Component) error {
	if err := os.RemoveAll(r.path); err != nil {
		return err
	}

	if err := os.Mkdir(r.path, os.ModePerm); err != nil {
		return err
	}

	blobdir := filepath.Join(r.path, "blobs")
	if err := os.Mkdir(blobdir, os.ModePerm); err != nil {
		return err
	}

	resources, err := component.Resources()
	if err != nil {
		return err
	}

	visitedResources := make([]v2.Resource, 0)
	for _, item := range resources {
		if item.Deferrable() {
			visitedResources = append(visitedResources, item)
			continue
		}

		dig, err := item.Digest()
		if err != nil {
			return err
		}

		p := filepath.Join(blobdir, dig.Value)

		f, err := os.Create(p)
		if err != nil {
			return err
		}
		defer f.Close()

		acc, err := item.Access()
		if err != nil {
			return err
		}

		reader, err := acc.Data()
		if err != nil {
			return err
		}
		defer reader.Close()

		if _, err := io.Copy(f, reader); err != nil {
			return err
		}

		newAccess, err := FromFile(p, WithMediaType(acc.MediaType()))
		if err != nil {
			return err
		}

		item = mutate.WithAccess(item, newAccess)

		visitedResources = append(visitedResources, item)
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
