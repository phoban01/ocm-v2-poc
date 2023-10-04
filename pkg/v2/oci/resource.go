package oci

import (
	"encoding/json"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	v2 "github.com/phoban01/ocm-v2/pkg/v2"
	"github.com/phoban01/ocm-v2/pkg/v2/types"
)

type image struct {
	name   string
	ref    string
	digest *types.Digest
	img    v1.Image
}

var _ v2.Resource = (*image)(nil)

const Type types.ResourceType = "ociImage"

func Resource(name, ref string) v2.Resource {
	return &image{name: name, ref: ref}
}

func (f *image) compute() error {
	ref, err := name.ParseReference(f.ref)
	if err != nil {
		return err
	}

	img, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return err
	}

	f.img = img

	hash, err := img.Digest()
	if err != nil {
		return err
	}

	f.digest = &types.Digest{
		HashAlgorithm:          "sha256",
		NormalisationAlgorithm: "json/v1",
		Value:                  strings.TrimPrefix(hash.String(), "sha256:"),
	}

	return nil
}

func (f *image) Name() string {
	return f.name
}

func (f *image) Access() v2.Access {
	if err := f.compute(); err != nil {
		return nil
	}
	return &artifactAccess{image: f}
}

func (f *image) Deferrable() bool {
	return true
}

func (f *image) Type() types.ResourceType {
	return Type
}

func (f *image) Labels() map[string]string {
	return map[string]string{
		"ocm.software/reference": f.ref,
	}
}

func (f *image) Digest() (types.Digest, error) {
	if f.digest != nil {
		return *f.digest, nil
	}
	if err := f.compute(); err != nil {
		return types.Digest{}, err
	}
	return *f.digest, nil
}

func (f image) WithLocation(url string) v2.Resource {
	return &image{name: f.name, ref: url}
}

func (f *image) MarshalJSON() ([]byte, error) {
	return json.Marshal(f)
}

func (f *image) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, f)
}
