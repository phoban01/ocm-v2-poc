package file

import (
	"crypto/sha256"
	"fmt"
	"io"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

type file struct {
	name string
	path string
}

const Type types.ResourceType = "file"

var _ v2.Resource = (*file)(nil)

func Resource(name, path string) v2.Resource {
	return &file{name: name, path: path}
}

func (f *file) Name() string {
	return f.name
}

func (f *file) Type() types.ResourceType {
	return Type
}

func (f *file) Labels() map[string]string {
	return map[string]string{
		"ocm.software/filename": f.path,
	}
}

func (f *file) Deferrable() bool {
	return false
}

func (f *file) Access() v2.Access {
	return &access{file: f}
}

func (f *file) Digest() (types.Digest, error) {
	data, err := f.Access().Data()
	if err != nil {
		return types.Digest{}, err
	}
	defer data.Close()

	hash := sha256.New()

	_, err = io.Copy(hash, data)
	if err != nil {
		return types.Digest{}, err
	}

	return types.Digest{
		HashAlgorithm:          "sha256",
		NormalisationAlgorithm: "json/v1",
		Value:                  fmt.Sprintf("%x", hash.Sum(nil)),
	}, nil
}
