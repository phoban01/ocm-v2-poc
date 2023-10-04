package file

import (
	"crypto/sha256"
	"encoding/json"
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

func (f file) WithLocation(p string) v2.Resource {
	return &file{name: f.name, path: p}
}

func (f *file) MarshalJSON() ([]byte, error) {
	access, err := json.Marshal(f.Access())
	if err != nil {
		return nil, err
	}
	dig, err := f.Digest()
	if err != nil {
		return nil, err
	}
	r := types.Resource{
		Name:   f.name,
		Type:   f.Type(),
		Access: access,
		Digest: dig,
	}
	return json.Marshal(r)
}

func (f *file) UnmarshalJSON(data []byte) error {
	r := types.Resource{}
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}

	a := access{}
	if err := json.Unmarshal(r.Access, &a); err != nil {
		return err
	}

	f.name = r.Name
	f.path = a.file.path

	return nil
}
