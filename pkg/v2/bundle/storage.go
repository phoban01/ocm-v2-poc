package bundle

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/mutate"
)

type storage struct {
	path string
}

var _ v2.Storage = (*storage)(nil)

func New(path string) (v2.Storage, error) {
	p, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return &storage{path: p}, nil
}

func (s *storage) Context() *v2.StorageContext {
	return &v2.StorageContext{
		Type: "ocm.storage/bundle",
		URL:  fmt.Sprintf("file://%s", s.path),
	}
}

func (s *storage) Write(component v2.Component) error {
	if err := os.RemoveAll(s.path); err != nil {
		return err
	}

	if err := os.Mkdir(s.path, os.ModePerm); err != nil {
		return err
	}

	blobdir := filepath.Join(s.path, "blobs")
	if err := os.Mkdir(blobdir, os.ModePerm); err != nil {
		return err
	}

	resources, err := component.Resources()
	if err != nil {
		return err
	}

	for _, r := range resources {
		if r.Deferrable() {
			continue
		}

		dig, err := r.Digest()
		if err != nil {
			return err
		}

		p := filepath.Join(blobdir, dig)
		f, err := os.Create(p)
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

		r = mutate.SetAccess(r, p)
		component = mutate.ReplaceResource(component, r)
	}

	desc, err := component.Descriptor()
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(s.path, "component-descriptor.json"))
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
