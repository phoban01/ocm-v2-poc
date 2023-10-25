package oci

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/types"
)

type localBlob struct {
	repository *repository
	desc       types.Descriptor
	data       func() (io.ReadCloser, error)
	length     int64
	mediaType  string
	digest     types.Digest
	ref        name.Digest
}

var _ v2.Access = (*localBlob)(nil)

func (a *localBlob) compute() error {
	registry := strings.TrimPrefix(a.repository.Context().Location(), "oci://")
	url := fmt.Sprintf(
		"%s/component-descriptors/%s@%s",
		registry,
		a.desc.ObjectMeta.Name,
		a.digest.Value,
	)
	ref, err := name.NewDigest(url)
	if err != nil {
		return err
	}
	a.ref = ref
	return nil
}

func (a *localBlob) Type() string {
	return "localBlob/v1"
}

func (a *localBlob) MediaType() string {
	return a.mediaType
}

func (a *localBlob) Reference() string {
	return a.ref.Name()
}

func (a *localBlob) Labels() map[string]string {
	return nil
}

func (a *localBlob) Length() (int64, error) {
	if a.length != 0 {
		return a.length, nil
	}
	if a.repository == nil {
		return 0, nil
	}
	store, err := a.repository.storage()
	if err != nil {
		return 0, err
	}
	descriptor, _, err := store.FetchReference(context.TODO(), a.digest.Value)
	return descriptor.Size, err
}

func (a *localBlob) Data() (io.ReadCloser, error) {
	if err := a.compute(); err != nil {
		return nil, err
	}
	return a.repository.ReadBlob(context.TODO(), a.digest.Value)
}

func (a *localBlob) Digest() (*types.Digest, error) {
	return &a.digest, nil
}

func (a *localBlob) MarshalJSON() ([]byte, error) {
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

func (a *localBlob) UnmarshalJSON(data []byte) error {
	obj := make(map[string]string)
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	return nil
}
