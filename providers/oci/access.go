package oci

import (
	"io"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/types"
)

func (a *accessor) Type() v2.AccessType {
	return v2.AccessType("ociArtifact")
}

func (a *accessor) MediaType() string {
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

func (a *accessor) WithLocation(p string) {
	a.ref = p
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
