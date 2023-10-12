package oci

import (
	"io"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/provider"
	"github.com/phoban01/ocm-v2/api/v2/types"
)

type accessor struct {
	repository v2.Repository
	image      v1.Image
	ref        string
	mediaType  string
}

var _ v2.Access = (*accessor)(nil)

var (
	AccessType = "ociArtifact"
	MediaType  = "application/vnd.docker.image"
)

func init() {
	provider.Register(&accessor{})
}

func (a *accessor) compute() error {
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

func (a *accessor) Type() string {
	return AccessType
}

func (a *accessor) MediaType() string {
	if a.mediaType != "" {
		return a.mediaType
	}
	return MediaType
}

func (a *accessor) Labels() map[string]string {
	return nil
}

func (a *accessor) Data() (io.ReadCloser, error) {
	if err := a.compute(); err != nil {
		return nil, err
	}
	return mutate.Extract(a.image), nil
}

func (a *accessor) Digest() (*types.Digest, error) {
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
