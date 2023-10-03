package archive

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/mutate"
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

		reader, err := item.Access().Data()
		if err != nil {
			return err
		}
		defer reader.Close()
		if _, err := io.Copy(f, reader); err != nil {
			return err
		}

		visitedResources = append(visitedResources, item.WithLocation(p))
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
