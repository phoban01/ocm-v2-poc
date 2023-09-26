package oci

import (
	"io"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

type image struct {
	name string
	ref  string
	i    v1.Image
}

var _ v2.Resource = (*image)(nil)

func Resource(name, ref string) v2.Resource {
	return &image{name: name, ref: ref}
}

func (f *image) compute() error {
	ref, err := name.ParseReference(f.Access())
	if err != nil {
		return err
	}

	img, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return err
	}

	f.i = img

	return nil
}

func (f *image) Name() string {
	return f.name
}

func (f *image) Access() string {
	return f.ref
}

func (f *image) Deferrable() bool {
	return true
}

func (f *image) Blob() (io.ReadCloser, error) {
	if err := f.compute(); err != nil {
		return nil, err
	}
	return mutate.Extract(f.i), nil
}

func (f *image) Digest() (string, error) {
	if err := f.compute(); err != nil {
		return "", err
	}
	hash, err := f.i.Digest()
	return strings.TrimPrefix(hash.String(), "sha256:"), err
}

func (f *image) ResourceType() types.ResourceType {
	return "image"
}

func (f *image) MediaType() types.MediaType {
	return "application/vnd.oci.image"
}

func (f *image) Labels() map[string]string {
	return map[string]string{
		"ocm.software/reference": f.ref,
	}
}
