package blob

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

type blob struct {
	name string
	path string
	data []byte
}

const Type types.ResourceType = "blob"

var _ v2.Resource = (*blob)(nil)

func FromBytes(name string, data []byte) v2.Resource {
	return &blob{name: name, data: data}
}

func FromFile(name string, path string) v2.Resource {
	return &blob{name: name, path: path}
}

func (f *blob) Name() string {
	return f.name
}

func (f *blob) Type() types.ResourceType {
	return Type
}

func (f *blob) Labels() map[string]string {
	return map[string]string{
		"ocm.software/blobname": f.path,
	}
}

func (f *blob) Deferrable() bool {
	return false
}

func (f *blob) Access() v2.Access {
	return &access{blob: f}
}

func (f *blob) Digest() (types.Digest, error) {
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

func (f blob) WithLocation(p string) v2.Resource {
	return &blob{name: f.name, path: p}
}

func (f *blob) MarshalJSON() ([]byte, error) {
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

func (f *blob) UnmarshalJSON(data []byte) error {
	r := types.Resource{}
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}

	a := access{}
	if err := json.Unmarshal(r.Access, &a); err != nil {
		return err
	}

	f.name = r.Name
	f.path = a.blob.path

	return nil
}
