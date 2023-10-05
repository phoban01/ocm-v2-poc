package blob

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/provider"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

type access struct {
	data      []byte
	path      string
	labels    map[string]string
	mediaType string
}

var _ v2.Access = (*access)(nil)

func init() {
	provider.Register(&access{})
}

func (a *access) Type() v2.AccessType {
	return v2.AccessType("localBlob/v1")
}

func (a *access) MediaType() string {
	if a.mediaType != "" {
		return a.mediaType
	}
	return "blob"
}

func (a *access) WithLocation(p string) {
	a.path = p
}

func (a *access) Labels() map[string]string {
	return a.labels
}

func (a *access) Data() (io.ReadCloser, error) {
	if a.data != nil {
		return io.NopCloser(bytes.NewReader(a.data)), nil
	}
	return os.Open(a.path)
}

func (a *access) Digest() (*types.Digest, error) {
	data, err := a.Data()
	if err != nil {
		return nil, err
	}
	defer data.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, data)
	if err != nil {
		return nil, err
	}

	return &types.Digest{
		HashAlgorithm:          "sha256",
		NormalisationAlgorithm: "json/v1",
		Value:                  fmt.Sprintf("%x", hash.Sum(nil)),
	}, nil
}

func (a *access) MarshalJSON() ([]byte, error) {
	result := map[string]string{
		"localReference": a.path,
		"type":           string(a.Type()),
		"mediaType":      a.MediaType(),
	}
	return json.Marshal(result)
}

func (a *access) UnmarshalJSON(data []byte) error {
	obj := make(map[string]string)
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	a.path = obj["localReference"]
	a.mediaType = obj["mediaType"]
	return nil
}

func (a *access) Decode(resource types.Resource) (v2.Access, error) {
	if err := json.Unmarshal(resource.Access, a); err != nil {
		return nil, fmt.Errorf("failed to unmarshal access: %w", err)
	}
	return a, nil
}
