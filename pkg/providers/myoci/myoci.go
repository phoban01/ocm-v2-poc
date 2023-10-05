package myoci

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/provider"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

type access struct {
	image v1.Image
	ref   string
}

var _ v2.Access = (*access)(nil)

var MediaType = "application/vnd.custom.image"

func Use() {
	provider.Register(&access{})
}

func (a *access) Type() v2.AccessType {
	return v2.AccessType("ociArtifact")
}

func (a *access) MediaType() string {
	return MediaType
}

func (a *access) Labels() map[string]string {
	return nil
}

func (a *access) Data() (io.ReadCloser, error) {
	if err := a.compute(); err != nil {
		return nil, err
	}
	return mutate.Extract(a.image), nil
}

func (a *access) Digest() (*types.Digest, error) {
	if err := a.compute(); err != nil {
		return nil, err
	}

	hash, err := a.image.Digest()
	if err != nil {
		return nil, err
	}

	return &types.Digest{
		HashAlgorithm:          "sha256",
		NormalisationAlgorithm: "json/v1",
		Value:                  strings.TrimPrefix(hash.String(), "sha256:"),
	}, nil
}

func (a *access) MarshalJSON() ([]byte, error) {
	result := map[string]string{
		"imageReference": a.ref,
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
	a.ref = obj["imageReference"]
	return nil
}

func (a *access) Decode(resource types.Resource) (v2.Access, error) {
	if err := json.Unmarshal(resource.Access, a); err != nil {
		return nil, fmt.Errorf("failed to unmarshal access: %w", err)
	}
	return a, nil
}

func (a *access) WithLocation(p string) {
	a.ref = p
}

func (a *access) compute() error {
	ref, err := name.ParseReference(a.ref)
	if err != nil {
		return err
	}
	img, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return err
	}
	a.image = img
	return nil
}
