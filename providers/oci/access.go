package oci

import (
	"context"
	"fmt"
	"io"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/provider"
	"github.com/phoban01/ocm-v2/api/v2/types"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
)

type accessor struct {
	initialized bool
	repository  v2.Repository
	mediaType   string
	labels      map[string]string
	ref         string
	desc        ocispec.Descriptor
	reader      io.ReadCloser
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
	if a.initialized {
		return nil
	}

	ref, err := registry.ParseReference(a.ref)
	if err != nil {
		return err
	}

	repo, err := remote.NewRepository(fmt.Sprintf("%s/%s", ref.Registry, ref.Repository))
	if err != nil {
		return err
	}

	desc, reader, err := oras.Fetch(context.TODO(), repo, ref.Reference, oras.DefaultFetchOptions)
	if err != nil {
		return err
	}

	a.reader = reader
	a.desc = desc
	a.initialized = true

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

func (a *accessor) Length() (int64, error) {
	if err := a.compute(); err != nil {
		return 0, err
	}
	return a.desc.Size, nil
}

func (a *accessor) Labels() map[string]string {
	return a.labels
}

func (a *accessor) Reference() string {
	return a.ref
}

func (a *accessor) Data() (io.ReadCloser, error) {
	if err := a.compute(); err != nil {
		return nil, err
	}
	return a.reader, nil
}

func (a *accessor) Digest() (*types.Digest, error) {
	if err := a.compute(); err != nil {
		return nil, err
	}
	return &types.Digest{
		HashAlgorithm:          a.desc.Digest.Algorithm().String(),
		Value:                  a.desc.Digest.String(),
		NormalisationAlgorithm: "json/v1",
	}, nil
}
