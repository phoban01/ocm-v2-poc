package myblob

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"

	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/provider"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

type accessor struct {
	filepath string
	labels   map[string]string
}

var _ v2.Access = (*accessor)(nil)

var MediaType = "custom-blob"

func Use() {
	provider.Register(&accessor{})
}

func (a *accessor) Type() v2.AccessType {
	return v2.AccessType("localBlob/v1")
}

func (a *accessor) MediaType() string {
	return MediaType
}

func (a *accessor) Labels() map[string]string {
	return a.labels
}

func (a *accessor) Data() (io.ReadCloser, error) {
	return os.Open(a.filepath)
}

func (a *accessor) Digest() (*types.Digest, error) {
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

func (a *accessor) Decode(resource types.Resource) (v2.Access, error) {
	if err := json.Unmarshal(resource.Access, a); err != nil {
		return nil, fmt.Errorf("failed to unmarshal access: %w", err)
	}
	return a, nil
}

func (a *accessor) WithLocation(p string) {
	a.filepath = p
}

func (a *accessor) MarshalJSON() ([]byte, error) {
	result := map[string]string{
		"localReference": a.filepath,
		"type":           string(a.Type()),
		"mediaType":      a.MediaType(),
	}
	return json.Marshal(result)
}

func (a *accessor) UnmarshalJSON(data []byte) error {
	obj := make(map[string]string)
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	a.filepath = obj["localReference"]
	return nil
}
