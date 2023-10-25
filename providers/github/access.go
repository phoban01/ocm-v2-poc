package github

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/provider"
	"github.com/phoban01/ocm-v2/api/v2/types"
)

var (
	AccessType = "localBlob/v1"
	MediaType  = "application/octet-stream"
)

type accessor struct {
	data       []byte
	basepath   string
	filepath   string
	mediaType  string
	labels     map[string]string
	digest     types.Digest
	repository v2.Repository
}

var _ v2.Access = (*accessor)(nil)

func init() {
	provider.Register(&accessor{})
}

func (a *accessor) Type() string {
	return AccessType
}

func (a *accessor) MediaType() string {
	if a.mediaType != "" {
		return a.mediaType
	}
	return MediaType
}

func (a *accessor) Reference() string {
	return a.filepath
}

func (a *accessor) Length() (int64, error) {
	return 0, nil
}

func (a *accessor) Labels() map[string]string {
	return a.labels
}

func (a *accessor) Data() (io.ReadCloser, error) {
	if a.data != nil {
		return io.NopCloser(bytes.NewReader(a.data)), nil
	}
	return nil, nil
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

func (*accessor) Decode(ctx v2.RepositoryContext, resource types.Resource) (v2.Access, error) {
	a := &accessor{}
	if err := json.Unmarshal(resource.Access, a); err != nil {
		return nil, fmt.Errorf("failed to unmarshal access: %w", err)
	}
	a.basepath = filepath.Join(strings.TrimPrefix(ctx.Location(), "file://"), "blobs")
	return a, nil
}

func (a *accessor) MarshalJSON() ([]byte, error) {
	dig, err := a.Digest()
	if err != nil {
		return nil, err
	}
	result := map[string]string{
		"localReference": fmt.Sprintf("sha256:%s", dig.Value),
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
	a.mediaType = obj["mediaType"]
	return nil
}
