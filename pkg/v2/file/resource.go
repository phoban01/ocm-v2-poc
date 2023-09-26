package file

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

type file struct {
	name string
	path string
}

var _ v2.Resource = (*file)(nil)

func Resource(name, path string) v2.Resource {
	return &file{name: name, path: path}
}

func (f *file) Name() string {
	return f.name
}

func (f *file) Access() string {
	return f.path
}

func (f *file) Deferrable() bool {
	return false
}

func (f *file) Blob() (io.ReadCloser, error) {
	return os.Open(f.Access())
}

func (f *file) Digest() (string, error) {
	data, err := f.Blob()
	if err != nil {
		return "", err
	}
	defer data.Close()

	hash := sha256.New()

	_, err = io.Copy(hash, data)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func (f *file) ResourceType() types.ResourceType {
	return "file"
}

func (f *file) MediaType() types.MediaType {
	return "application/x-yaml"
}

func (f *file) Labels() map[string]string {
	return map[string]string{
		"ocm.software/filename": f.path,
	}
}
