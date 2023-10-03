package file

import (
	"encoding/json"
	"fmt"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

func (f file) WithLocation(p string) v2.Resource {
	return &file{name: f.name, path: p}
}

func DecodeResource(resource types.Resource) (v2.Resource, error) {
	a := access{}
	if err := json.Unmarshal(resource.Access, &a); err != nil {
		return nil, fmt.Errorf("failed to unmarshal access: %w", err)
	}
	return &file{
		name: resource.Name,
		path: a.file.path,
	}, nil
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
