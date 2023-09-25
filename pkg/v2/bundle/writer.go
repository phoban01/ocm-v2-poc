package bundle

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
)

func Write(path string, component v2.Component) error {
	if err := os.RemoveAll(path); err != nil {
		return err
	}

	if err := os.Mkdir(path, os.ModePerm); err != nil {
		return err
	}

	blobdir := filepath.Join(path, "blobs")
	if err := os.Mkdir(blobdir, os.ModePerm); err != nil {
		return err
	}

	resources, err := component.Resources()
	if err != nil {
		return err
	}

	for _, r := range resources {
		dig, err := r.Digest()
		if err != nil {
			return err
		}
		f, err := os.Create(filepath.Join(blobdir, dig))
		if err != nil {
			return err
		}
		defer f.Close()
		reader, err := r.Blob()
		if err != nil {
			return err
		}
		defer reader.Close()
		if _, err := io.Copy(f, reader); err != nil {
			return err
		}
	}

	desc, err := component.Descriptor()
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(path, "component-descriptor.json"))
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
