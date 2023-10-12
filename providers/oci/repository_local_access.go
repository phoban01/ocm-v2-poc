package oci

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/types"
)

type localAccess struct {
	repository v2.Repository
	desc       v2.Descriptor
	data       func() (io.ReadCloser, error)
	mediaType  string
	digest     types.Digest
	ref        name.Digest
}

var _ v2.Access = (*localAccess)(nil)

func (a *localAccess) compute() error {
	registry := strings.TrimPrefix(a.repository.Context().Location(), "oci://")
	url := fmt.Sprintf("%s/component-descriptors/%s@%s", registry, a.desc.Name, a.digest.Value)
	ref, err := name.NewDigest(url)
	if err != nil {
		return err
	}
	a.ref = ref
	return nil
}

func (a *localAccess) Type() string {
	return "localBlob/v1"
}

func (a *localAccess) MediaType() string {
	return a.mediaType
}

func (a *localAccess) Labels() map[string]string {
	return nil
}

func (a *localAccess) Data() (io.ReadCloser, error) {
	if err := a.compute(); err != nil {
		return nil, err
	}
	if a.data != nil {
		return a.data()
	}
	layer, err := remote.Layer(a.ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return nil, err
	}
	return layer.Uncompressed()
}

func (a *localAccess) Digest() (*types.Digest, error) {
	return &a.digest, nil
}

func (a *localAccess) MarshalJSON() ([]byte, error) {
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

func (a *localAccess) UnmarshalJSON(data []byte) error {
	obj := make(map[string]string)
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	return nil
}
